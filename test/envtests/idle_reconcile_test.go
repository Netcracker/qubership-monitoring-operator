//go:build envtest

package envtests

import (
	"context"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

// updateRecordingClient returns a client that records "Kind/name" for every
// Update request. The API server keeps the resourceVersion of an object when
// an Update changes nothing, so only the requests themselves reveal a
// redundant write.
func updateRecordingClient(updates *[]string) client.Client {
	base, err := client.NewWithWatch(cfg, client.Options{Scheme: clientgoscheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	return interceptor.NewClient(base, interceptor.Funcs{
		Update: func(ctx context.Context, c client.WithWatch, obj client.Object, opts ...client.UpdateOption) error {
			gvk, err := c.GroupVersionKindFor(obj)
			if err != nil {
				return err
			}
			*updates = append(*updates, gvk.Kind+"/"+obj.GetName())
			return c.Update(ctx, obj, opts...)
		},
	})
}

var _ = Describe("Idle reconciliation", func() {
	It("does not update managed resources that already match", func() {
		ctx := context.Background()
		const namespaceName = "idle-reconcile"

		ensureNamespace(namespaceName)

		resource := &monv1.PlatformMonitoring{}
		Expect(yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, "assets/platformmonitoring.yaml"), 100).Decode(resource)).To(Succeed())
		resource.SetNamespace(namespaceName)
		// The asset's Alertmanager NodePort belongs to the Reconcile specs.
		resource.Spec.AlertManager.Port = 30904
		resource.Spec.GatewayAPI = &monv1.GatewayAPI{
			ParentRefs: []monv1.GatewayParentRef{{Name: "gateway", Namespace: "gateway-infra", SectionName: "http"}},
		}
		httpRoute := func(hostname string) *monv1.GatewayHTTPRoute {
			return &monv1.GatewayHTTPRoute{Install: ptr.To(true), Hostnames: []string{hostname}}
		}
		resource.Spec.Prometheus.HTTPRoute = httpRoute("prometheus.idle.test")
		resource.Spec.AlertManager.HTTPRoute = httpRoute("alertmanager.idle.test")
		resource.Spec.Grafana.HTTPRoute = httpRoute("grafana.idle.test")
		resource.Spec.Pushgateway.HTTPRoute = httpRoute("pushgateway.idle.test")
		Expect(k8sClient.Create(ctx, resource)).To(Succeed())

		var updates []string
		reconciler := &controllers.PlatformMonitoringReconciler{
			Client:          updateRecordingClient(&updates),
			Scheme:          clientgoscheme.Scheme,
			Config:          cfg,
			DiscoveryClient: discoveryClient,
		}
		request := ctrl.Request{NamespacedName: client.ObjectKeyFromObject(resource)}

		_, err := reconciler.Reconcile(ctx, request)
		Expect(err).NotTo(HaveOccurred())
		for _, key := range []client.ObjectKey{
			{Namespace: namespaceName, Name: utils.GrafanaOperatorComponentName},
			{Namespace: namespaceName, Name: utils.PushgatewayComponentName},
		} {
			Expect(k8sClient.Get(ctx, key, &appsv1.Deployment{})).To(Succeed())
		}
		Expect(k8sClient.Get(ctx, client.ObjectKey{Namespace: namespaceName, Name: "k8s"}, &promv1.Prometheus{})).To(Succeed())
		routeName := namespaceName + "-" + utils.PrometheusComponentName + "-http-route"
		expectHTTPRoute(namespaceName, routeName, "prometheus.idle.test")
		route := &unstructured.Unstructured{}
		route.SetGroupVersionKind(schema.GroupVersionKind{Group: "gateway.networking.k8s.io", Version: "v1", Kind: "HTTPRoute"})
		Expect(k8sClient.Get(ctx, client.ObjectKey{Namespace: namespaceName, Name: routeName}, route)).To(Succeed())
		Expect(route.GetLabels()).NotTo(BeEmpty())
		rules, found, err := unstructured.NestedSlice(route.Object, "spec", "rules")
		Expect(err).NotTo(HaveOccurred())
		Expect(found).To(BeTrue())
		Expect(rules[0]).To(HaveKeyWithValue("backendRefs", ConsistOf(HaveKeyWithValue("weight", BeEquivalentTo(1)))),
			"the HTTPRoute CRD must apply schema defaults")

		updates = nil
		_, err = reconciler.Reconcile(ctx, request)
		Expect(err).NotTo(HaveOccurred())

		Expect(updates).To(BeEmpty())
	})
})
