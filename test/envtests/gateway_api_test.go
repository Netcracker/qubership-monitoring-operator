package bdd_tests

import (
	"context"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/alertmanager"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/grafana"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/prometheus"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/pushgateway"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	vmagent "github.com/Netcracker/qubership-monitoring-operator/controllers/victoriametrics/vmagent"
	vmalert "github.com/Netcracker/qubership-monitoring-operator/controllers/victoriametrics/vmalert"
	vmalertmanager "github.com/Netcracker/qubership-monitoring-operator/controllers/victoriametrics/vmalertmanager"
	vmauth "github.com/Netcracker/qubership-monitoring-operator/controllers/victoriametrics/vmauth"
	vmsingle "github.com/Netcracker/qubership-monitoring-operator/controllers/victoriametrics/vmsingle"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func newGatewayCR(name, namespace string) monv1.PlatformMonitoring {
	created := monv1.PlatformMonitoring{}
	err = yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, "assets/platformmonitoring.yaml"), 100).Decode(&created)
	Expect(err).NotTo(HaveOccurred())

	created.SetName(name)
	created.SetNamespace(namespace)
	created.SetUID(uuid.NewUUID())

	created.Spec.GatewayAPI = &monv1.GatewayAPI{
		AddIngressIgnoreAnnotation: true,
		ParentRefs: []monv1.GatewayParentRef{
			{
				Name:        "gateway",
				Namespace:   "gateway-infra",
				Group:       "gateway.networking.k8s.io",
				Kind:        "Gateway",
				SectionName: "http",
			},
		},
	}
	return created
}

func ensureNamespace(name string) {
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: name}}
	_ = k8sClient.Create(context.TODO(), ns)
}

func expectIngressAnnotation(namespace, name string) {
	ingress := netv1.Ingress{ObjectMeta: metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
	}}
	err = k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(&ingress), &ingress)
	Expect(err).NotTo(HaveOccurred())
	Expect(ingress.Annotations).NotTo(BeNil())
	Expect(ingress.Annotations["gateway-api-converter.netcracker.com/ignore"]).To(Equal("true"))
}

func expectHTTPRoute(namespace, name, hostname string) {
	route := &unstructured.Unstructured{}
	route.SetGroupVersionKind(schema.GroupVersionKind{Group: "gateway.networking.k8s.io", Version: "v1", Kind: "HTTPRoute"})
	route.SetName(name)
	route.SetNamespace(namespace)
	err = k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(route), route)
	Expect(err).NotTo(HaveOccurred())

	spec, ok := route.Object["spec"].(map[string]interface{})
	Expect(ok).To(BeTrue())
	hostnames := spec["hostnames"].([]interface{})
	Expect(hostnames).To(ContainElement(hostname))
	parentRefs := spec["parentRefs"].([]interface{})
	Expect(parentRefs).To(HaveLen(1))
	parent := parentRefs[0].(map[string]interface{})
	Expect(parent["name"]).To(Equal("gateway"))
	Expect(parent["namespace"]).To(Equal("gateway-infra"))
}

var _ = Describe("Gateway API", func() {
	It("creates HTTPRoute and adds ingress ignore annotation", func() {
		gatewayCr := newGatewayCR("platformmonitoring-gateway", "gateway-test")
		gatewayCr.Spec.Grafana.HTTPRoute = &monv1.GatewayHTTPRoute{
			Install:   ptr.To(true),
			Hostnames: []string{"grafana.gateway.test"},
		}
		ensureNamespace(gatewayCr.GetNamespace())
		gReconciler := grafana.NewGrafanaReconciler(k8sClient, scheme.Scheme, discoveryClient, cfg)
		err = gReconciler.Run(&gatewayCr)
		Expect(err).NotTo(HaveOccurred())

		expectIngressAnnotation(gatewayCr.GetNamespace(), gatewayCr.GetNamespace()+"-"+utils.GrafanaComponentName)
		expectHTTPRoute(gatewayCr.GetNamespace(), gatewayCr.GetNamespace()+"-"+utils.GrafanaComponentName+"-http-route", "grafana.gateway.test")
	})

	It("creates Prometheus HTTPRoute and adds ingress ignore annotation", func() {
		cr := newGatewayCR("platformmonitoring-prom", "gateway-prom")
		cr.Spec.Prometheus.HTTPRoute = &monv1.GatewayHTTPRoute{
			Install:   ptr.To(true),
			Hostnames: []string{"prom.gateway.test"},
		}
		ensureNamespace(cr.GetNamespace())

		pReconciler := prometheus.NewPrometheusReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = pReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		expectIngressAnnotation(cr.GetNamespace(), cr.GetNamespace()+"-"+utils.PrometheusComponentName)
		expectHTTPRoute(cr.GetNamespace(), cr.GetNamespace()+"-"+utils.PrometheusComponentName+"-http-route", "prom.gateway.test")
	})

	It("creates Alertmanager HTTPRoute and adds ingress ignore annotation", func() {
		cr := newGatewayCR("platformmonitoring-am", "gateway-am")
		cr.Spec.AlertManager.HTTPRoute = &monv1.GatewayHTTPRoute{
			Install:   ptr.To(true),
			Hostnames: []string{"am.gateway.test"},
		}
		cr.Spec.AlertManager.Port = 30913
		ensureNamespace(cr.GetNamespace())

		aReconciler := alertmanager.NewAlertManagerReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = aReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		expectIngressAnnotation(cr.GetNamespace(), cr.GetNamespace()+"-"+utils.AlertManagerComponentName)
		expectHTTPRoute(cr.GetNamespace(), cr.GetNamespace()+"-"+utils.AlertManagerComponentName+"-http-route", "am.gateway.test")
	})

	It("creates Pushgateway HTTPRoute and adds ingress ignore annotation", func() {
		cr := newGatewayCR("platformmonitoring-pg", "gateway-pg")
		cr.Spec.Pushgateway.HTTPRoute = &monv1.GatewayHTTPRoute{
			Install:   ptr.To(true),
			Hostnames: []string{"pg.gateway.test"},
		}
		ensureNamespace(cr.GetNamespace())

		pgReconciler := pushgateway.NewPushgatewayReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = pgReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		expectIngressAnnotation(cr.GetNamespace(), cr.GetNamespace()+"-"+utils.PushgatewayComponentName)
		expectHTTPRoute(cr.GetNamespace(), cr.GetNamespace()+"-"+utils.PushgatewayComponentName+"-http-route", "pg.gateway.test")
	})

	It("creates VmAgent HTTPRoute and adds ingress ignore annotation", func() {
		cr := newGatewayCR("platformmonitoring-vmagent", "gateway-vmagent")
		cr.Spec.Victoriametrics = &monv1.Victoriametrics{
			VmAgent: monv1.VmAgent{
				Install: ptr.To(true),
				Image:   "victoriametrics/vmagent:v1.99.0",
				RemoteWrite: []vmetricsv1b1.VMAgentRemoteWriteSpec{
					{URL: "http://example.com/api/v1/write"},
				},
				Ingress: &monv1.Ingress{Install: ptr.To(true), Host: "vmagent.gateway.test"},
				HTTPRoute: &monv1.GatewayHTTPRoute{
					Install:   ptr.To(true),
					Hostnames: []string{"vmagent.gateway.test"},
				},
			},
		}
		ensureNamespace(cr.GetNamespace())

		r := vmagent.NewVmAgentReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = r.Run(context.TODO(), &cr)
		Expect(err).NotTo(HaveOccurred())

		expectIngressAnnotation(cr.GetNamespace(), cr.GetNamespace()+"-"+utils.VmAgentServiceName)
		expectHTTPRoute(cr.GetNamespace(), cr.GetNamespace()+"-"+utils.VmAgentServiceName+"-http-route", "vmagent.gateway.test")
	})

	It("creates VmAuth HTTPRoute and adds ingress ignore annotation", func() {
		cr := newGatewayCR("platformmonitoring-vmauth", "gateway-vmauth")
		cr.Spec.Victoriametrics = &monv1.Victoriametrics{
			VmAuth: monv1.VmAuth{
				Install: ptr.To(true),
				Image:   "victoriametrics/vmauth:v1.99.0",
				Ingress: &monv1.Ingress{Install: ptr.To(true), Host: "vmauth.gateway.test"},
				HTTPRoute: &monv1.GatewayHTTPRoute{
					Install:   ptr.To(true),
					Hostnames: []string{"vmauth.gateway.test"},
				},
			},
		}
		ensureNamespace(cr.GetNamespace())

		r := vmauth.NewVmAuthReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = r.Run(context.TODO(), &cr)
		Expect(err).NotTo(HaveOccurred())

		expectIngressAnnotation(cr.GetNamespace(), cr.GetNamespace()+"-"+utils.VmAuthServiceName)
		expectHTTPRoute(cr.GetNamespace(), cr.GetNamespace()+"-"+utils.VmAuthServiceName+"-http-route", "vmauth.gateway.test")
	})

	It("creates VmAlert HTTPRoute and adds ingress ignore annotation", func() {
		cr := newGatewayCR("platformmonitoring-vmalert", "gateway-vmalert")
		cr.Spec.Victoriametrics = &monv1.Victoriametrics{
			VmAlert: monv1.VmAlert{
				Install: ptr.To(true),
				Image:   "victoriametrics/vmalert:v1.99.0",
				Ingress: &monv1.Ingress{Install: ptr.To(true), Host: "vmalert.gateway.test"},
				HTTPRoute: &monv1.GatewayHTTPRoute{
					Install:   ptr.To(true),
					Hostnames: []string{"vmalert.gateway.test"},
				},
			},
		}
		ensureNamespace(cr.GetNamespace())

		r := vmalert.NewVmAlertReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = r.Run(context.TODO(), &cr)
		Expect(err).NotTo(HaveOccurred())

		expectIngressAnnotation(cr.GetNamespace(), cr.GetNamespace()+"-"+utils.VmAlertServiceName)
		expectHTTPRoute(cr.GetNamespace(), cr.GetNamespace()+"-"+utils.VmAlertServiceName+"-http-route", "vmalert.gateway.test")
	})

	It("creates VmAlertManager HTTPRoute and adds ingress ignore annotation", func() {
		cr := newGatewayCR("platformmonitoring-vmalertmanager", "gateway-vmalertmanager")
		cr.Spec.Victoriametrics = &monv1.Victoriametrics{
			VmAlertManager: monv1.VmAlertManager{
				Install: ptr.To(true),
				Image:   "victoriametrics/vmalertmanager:v1.99.0",
				Ingress: &monv1.Ingress{Install: ptr.To(true), Host: "vmalertmanager.gateway.test"},
				HTTPRoute: &monv1.GatewayHTTPRoute{
					Install:   ptr.To(true),
					Hostnames: []string{"vmalertmanager.gateway.test"},
				},
			},
		}
		ensureNamespace(cr.GetNamespace())

		r := vmalertmanager.NewVmAlertManagerReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = r.Run(context.TODO(), &cr)
		Expect(err).NotTo(HaveOccurred())

		expectIngressAnnotation(cr.GetNamespace(), cr.GetNamespace()+"-"+utils.VmAlertManagerServiceName)
		expectHTTPRoute(cr.GetNamespace(), cr.GetNamespace()+"-"+utils.VmAlertManagerServiceName+"-http-route", "vmalertmanager.gateway.test")
	})

	It("creates VmSingle HTTPRoute and adds ingress ignore annotation", func() {
		cr := newGatewayCR("platformmonitoring-vmsingle", "gateway-vmsingle")
		cr.Spec.Victoriametrics = &monv1.Victoriametrics{
			VmSingle: monv1.VmSingle{
				Install: ptr.To(true),
				Image:   "victoriametrics/vmsingle:v1.99.0",
				// VMSingle CRD requires retentionPeriod to be non-empty.
				RetentionPeriod: "1d",
				Ingress:         &monv1.Ingress{Install: ptr.To(true), Host: "vmsingle.gateway.test"},
				HTTPRoute: &monv1.GatewayHTTPRoute{
					Install:   ptr.To(true),
					Hostnames: []string{"vmsingle.gateway.test"},
				},
			},
		}
		ensureNamespace(cr.GetNamespace())

		r := vmsingle.NewVmSingleReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = r.Run(context.TODO(), &cr)
		Expect(err).NotTo(HaveOccurred())

		expectIngressAnnotation(cr.GetNamespace(), cr.GetNamespace()+"-"+utils.VmSingleServiceName)
		expectHTTPRoute(cr.GetNamespace(), cr.GetNamespace()+"-"+utils.VmSingleServiceName+"-http-route", "vmsingle.gateway.test")
	})

})
