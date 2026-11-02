package utils

import (
	"testing"

	secv1 "github.com/openshift/api/security/v1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

// assertSync checks that sync reports no change for equal objects, copies a
// change made by mutate, and then reports no change again.
func assertSync[T any](t *testing.T, sync func(existing, desired T) bool, newObject func() T, mutate func(T)) {
	t.Helper()
	existing, desired := newObject(), newObject()
	assert.False(t, sync(existing, desired), "equal objects")

	mutate(desired)
	assert.True(t, sync(existing, desired), "changed desired object")
	assert.True(t, apiequality.Semantic.DeepEqual(existing, desired), "existing matches desired after sync")
	assert.False(t, sync(existing, desired), "second sync")
}

func syncTestMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{Name: "obj", Namespace: "monitoring", Labels: map[string]string{"app": "x"}}
}

func TestSyncFunctions(t *testing.T) {
	t.Run("labels", func(t *testing.T) {
		assertSync(t, SyncLabels[*corev1.ServiceAccount],
			func() *corev1.ServiceAccount { return &corev1.ServiceAccount{ObjectMeta: syncTestMeta()} },
			func(o *corev1.ServiceAccount) { o.Labels["app"] = "y" })
	})
	t.Run("cluster role", func(t *testing.T) {
		assertSync(t, SyncClusterRole,
			func() *rbacv1.ClusterRole { return &rbacv1.ClusterRole{ObjectMeta: syncTestMeta()} },
			func(o *rbacv1.ClusterRole) { o.Rules = []rbacv1.PolicyRule{{Verbs: []string{"get"}}} })
	})
	t.Run("role", func(t *testing.T) {
		assertSync(t, SyncRole,
			func() *rbacv1.Role { return &rbacv1.Role{ObjectMeta: syncTestMeta()} },
			func(o *rbacv1.Role) { o.Rules = []rbacv1.PolicyRule{{Verbs: []string{"list"}}} })
	})
	t.Run("service", func(t *testing.T) {
		assertSync(t, SyncService,
			func() *corev1.Service { return &corev1.Service{ObjectMeta: syncTestMeta()} },
			func(o *corev1.Service) {
				o.Spec.Ports = []corev1.ServicePort{{Name: "http", Port: 80}}
				o.Spec.Selector = map[string]string{"app": "x"}
			})
	})
	t.Run("ingress", func(t *testing.T) {
		assertSync(t, SyncIngress,
			func() *networkingv1.Ingress { return &networkingv1.Ingress{ObjectMeta: syncTestMeta()} },
			func(o *networkingv1.Ingress) {
				o.Annotations = map[string]string{"a": "b"}
				o.Spec.Rules = []networkingv1.IngressRule{{Host: "example.com"}}
				o.Spec.TLS = []networkingv1.IngressTLS{{SecretName: "tls"}}
			})
	})
	t.Run("service monitor", func(t *testing.T) {
		assertSync(t, SyncServiceMonitor,
			func() *promv1.ServiceMonitor { return &promv1.ServiceMonitor{ObjectMeta: syncTestMeta()} },
			func(o *promv1.ServiceMonitor) {
				o.Spec.JobLabel = "job"
				o.Spec.Endpoints = []promv1.Endpoint{{Port: "http"}}
				o.Spec.NamespaceSelector = promv1.NamespaceSelector{Any: true}
				o.Spec.Selector = metav1.LabelSelector{MatchLabels: map[string]string{"app": "x"}}
			})
	})
	t.Run("pod monitor", func(t *testing.T) {
		assertSync(t, SyncPodMonitor,
			func() *promv1.PodMonitor { return &promv1.PodMonitor{ObjectMeta: syncTestMeta()} },
			func(o *promv1.PodMonitor) {
				o.Spec.JobLabel = "job"
				o.Spec.PodMetricsEndpoints = []promv1.PodMetricsEndpoint{{Port: ptr.To("http")}}
				o.Spec.NamespaceSelector = promv1.NamespaceSelector{Any: true}
				o.Spec.Selector = metav1.LabelSelector{MatchLabels: map[string]string{"app": "x"}}
			})
	})
	t.Run("security context constraints", func(t *testing.T) {
		assertSync(t, SyncSecurityContextConstraints,
			func() *secv1.SecurityContextConstraints {
				return &secv1.SecurityContextConstraints{ObjectMeta: syncTestMeta()}
			},
			func(o *secv1.SecurityContextConstraints) {
				o.Labels["app"] = "y"
				o.AllowHostNetwork = true
			})
	})
	t.Run("pod spec", func(t *testing.T) {
		assertSync(t, SyncPodSpec,
			func() *corev1.PodSpec { return &corev1.PodSpec{} },
			func(o *corev1.PodSpec) {
				o.SecurityContext = &corev1.PodSecurityContext{RunAsNonRoot: ptr.To(true)}
				o.Containers = []corev1.Container{{Name: "app", Image: "image"}}
				o.Volumes = []corev1.Volume{{Name: "data", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}}}
				o.ServiceAccountName = "sa"
				o.Affinity = &corev1.Affinity{}
				o.PriorityClassName = "high"
			})
	})
}

func TestSyncPodSpecLeavesUnmanagedFields(t *testing.T) {
	existing := &corev1.PodSpec{DNSPolicy: corev1.DNSClusterFirst, NodeSelector: map[string]string{"a": "b"}}
	desired := &corev1.PodSpec{ServiceAccountName: "sa"}

	assert.True(t, SyncPodSpec(existing, desired))
	assert.Equal(t, corev1.DNSClusterFirst, existing.DNSPolicy)
	assert.Equal(t, map[string]string{"a": "b"}, existing.NodeSelector)
	assert.Equal(t, "sa", existing.ServiceAccountName)
}

func TestAnyChanged(t *testing.T) {
	assert.False(t, AnyChanged())
	assert.False(t, AnyChanged(false, false))
	assert.True(t, AnyChanged(false, true, false))

	calls := 0
	track := func(result bool) bool {
		calls++
		return result
	}
	assert.True(t, AnyChanged(track(true), track(false), track(true)))
	assert.Equal(t, 3, calls, "every argument is evaluated")
}
