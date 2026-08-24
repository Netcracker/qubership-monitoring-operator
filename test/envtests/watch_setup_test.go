//go:build envtest

package envtests

import (
	"context"
	"os"
	"time"

	"github.com/Netcracker/qubership-monitoring-operator/controllers"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	routev1 "github.com/openshift/api/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/kubernetes/scheme"
	ktesting "k8s.io/client-go/testing"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlconfig "sigs.k8s.io/controller-runtime/pkg/config"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

var _ = Describe("Watch-based SetupWithManager", func() {
	It("starts, syncs cache, sees Grafana GVKs, and does not treat Route as served", func() {
		Expect(os.Setenv("RECONCILIATION_INTERVAL", "0")).To(Succeed())
		defer os.Unsetenv("RECONCILIATION_INTERVAL")

		prev := utils.PrivilegedRights
		utils.PrivilegedRights = true
		defer func() { utils.PrivilegedRights = prev }()

		mgr, err := ctrl.NewManager(cfg, ctrl.Options{
			Scheme:  scheme.Scheme,
			Metrics: metricsserver.Options{BindAddress: "0"},
		})
		Expect(err).NotTo(HaveOccurred())

		r := &controllers.PlatformMonitoringReconciler{
			Client:          mgr.GetClient(),
			Scheme:          mgr.GetScheme(),
			Log:             ctrl.Log.WithName("watch-spike"),
			Config:          cfg,
			DiscoveryClient: discoveryClient,
		}
		Expect(r.SetupWithManager(mgr)).To(Succeed())

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		go func() {
			defer GinkgoRecover()
			_ = mgr.Start(ctx)
		}()
		Expect(mgr.GetCache().WaitForCacheSync(ctx)).To(BeTrue())

		ok, err := utils.ResourceExists(discoveryClient, grafv1.SchemeGroupVersion.String(), "Grafana")
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue(), "envtest loads Grafana CRDs; Owns must be safe to register")

		ok, err = utils.ResourceExists(discoveryClient, grafv1.SchemeGroupVersion.String(), "GrafanaDashboard")
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())

		ok, err = utils.ResourceExists(discoveryClient, routev1.GroupVersion.String(), "Route")
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeFalse(), "vanilla envtest must not serve OpenShift Route")

		ok, err = utils.ResourceExists(discoveryClient, "gateway.networking.k8s.io/v1", "HTTPRoute")
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue(), "envtest loads the local HTTPRoute CRD; Owns must be safe to register")
	})

	It("starts when optional VM GVKs are absent from discovery", func() {
		Expect(os.Setenv("RECONCILIATION_INTERVAL", "0")).To(Succeed())
		defer os.Unsetenv("RECONCILIATION_INTERVAL")

		prev := utils.PrivilegedRights
		utils.PrivilegedRights = false
		defer func() { utils.PrivilegedRights = prev }()

		dc := &fakediscovery.FakeDiscovery{Fake: &ktesting.Fake{}}
		dc.Resources = []*metav1.APIResourceList{
			{
				GroupVersion: grafv1.SchemeGroupVersion.String(),
				APIResources: []metav1.APIResource{{Kind: "Grafana"}, {Kind: "GrafanaDashboard"}},
			},
		}

		mgr, err := ctrl.NewManager(cfg, ctrl.Options{
			Scheme:     scheme.Scheme,
			Metrics:    metricsserver.Options{BindAddress: "0"},
			Controller: ctrlconfig.Controller{SkipNameValidation: ptr.To(true)},
		})
		Expect(err).NotTo(HaveOccurred())

		r := &controllers.PlatformMonitoringReconciler{
			Client:          mgr.GetClient(),
			Scheme:          mgr.GetScheme(),
			Log:             ctrl.Log.WithName("watch-spike-no-vm"),
			Config:          cfg,
			DiscoveryClient: dc,
		}
		Expect(r.SetupWithManager(mgr)).To(Succeed())
	})
})
