//go:build envtest

package envtests

import (
	"context"
	"os"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("PlatformMonitoring status stability", func() {
	It("preserves status during an idle reconciliation", func() {
		ctx := context.Background()
		const namespaceName = "status-stability"

		originalInterval, intervalWasSet := os.LookupEnv("RECONCILIATION_INTERVAL")
		Expect(os.Setenv("RECONCILIATION_INTERVAL", "60")).To(Succeed())
		DeferCleanup(func() {
			if intervalWasSet {
				Expect(os.Setenv("RECONCILIATION_INTERVAL", originalInterval)).To(Succeed())
				return
			}
			Expect(os.Unsetenv("RECONCILIATION_INTERVAL")).To(Succeed())
		})

		namespace := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespaceName}}
		Expect(k8sClient.Create(ctx, namespace)).To(Succeed())
		DeferCleanup(func() {
			Expect(client.IgnoreNotFound(k8sClient.Delete(ctx, namespace))).To(Succeed())
		})

		resource := &monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "status-stability",
				Namespace: namespaceName,
			},
		}
		Expect(k8sClient.Create(ctx, resource)).To(Succeed())

		reconciler := &controllers.PlatformMonitoringReconciler{
			Client:          k8sClient,
			Scheme:          clientgoscheme.Scheme,
			Config:          cfg,
			DiscoveryClient: discoveryClient,
		}
		request := ctrl.Request{NamespacedName: client.ObjectKeyFromObject(resource)}

		_, err := reconciler.Reconcile(ctx, request)
		Expect(err).NotTo(HaveOccurred())

		first := &monv1.PlatformMonitoring{}
		Expect(k8sClient.Get(ctx, request.NamespacedName, first)).To(Succeed())
		Expect(first.Status.ObservedGeneration).To(Equal(first.Generation))
		Expect(first.Status.Conditions).To(HaveLen(1))
		Expect(first.Status.Conditions[0].Type).To(Equal("Successful"))

		firstResourceVersion := first.ResourceVersion
		firstConditions := append([]monv1.PlatformMonitoringCondition(nil), first.Status.Conditions...)

		_, err = reconciler.Reconcile(ctx, request)
		Expect(err).NotTo(HaveOccurred())

		second := &monv1.PlatformMonitoring{}
		Expect(k8sClient.Get(ctx, request.NamespacedName, second)).To(Succeed())
		Expect(second.ResourceVersion).To(Equal(firstResourceVersion))
		Expect(second.Status.ObservedGeneration).To(Equal(first.Status.ObservedGeneration))
		Expect(second.Status.Conditions).To(Equal(firstConditions))
	})
})
