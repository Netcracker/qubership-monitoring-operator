package prometheus

import (
	"context"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/labelsassert"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	fakediscovery "k8s.io/client-go/discovery/fake"
	ktesting "k8s.io/client-go/testing"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	cr              *monv1.PlatformMonitoring
	labelKey        = "label.key"
	labelValue      = "label-value"
	annotationKey   = "annotation.key"
	annotationValue = "annotation-value"
)

func TestPrometheusManifests(t *testing.T) {
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Prometheus: &monv1.Prometheus{
				Annotations: map[string]string{annotationKey: annotationValue},
				Labels:      map[string]string{labelKey: labelValue},
			},
		},
	}
	t.Run("Test Prometheus manifest", func(t *testing.T) {
		m, err := prometheus(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Prometheus manifest should not be empty")
		assert.NotNil(t, m.Spec.PodMetadata.Labels)
		assert.Equal(t, labelValue, m.Spec.PodMetadata.Labels[labelKey])
		assert.NotNil(t, m.Spec.PodMetadata.Annotations)
		assert.Equal(t, annotationValue, m.Spec.PodMetadata.Annotations[annotationKey])
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Prometheus: &monv1.Prometheus{},
		},
	}
	t.Run("Test Prometheus manifest with nil labels and annotation", func(t *testing.T) {
		m, err := prometheus(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Prometheus manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Nil(t, m.GetAnnotations())
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
	}
	t.Run("Test ServiceAccount manifest", func(t *testing.T) {
		crWithSALabels := &monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
			Spec: monv1.PlatformMonitoringSpec{
				Prometheus: &monv1.Prometheus{
					ServiceAccount: &monv1.EmbeddedObjectMetadata{
						Labels: map[string]string{labelKey: labelValue},
					},
				},
			},
		}
		m, err := prometheusServiceAccount(crWithSALabels)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ServiceAccount manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey], "ServiceAccount.Labels should be merged")
	})
	t.Run("Test ClusterRole manifest", func(t *testing.T) {
		m, err := prometheusClusterRole(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ClusterRole manifest should not be empty")
	})
	t.Run("Test ClusterRoleBinding manifest", func(t *testing.T) {
		m, err := prometheusClusterRoleBinding(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ClusterRoleBinding manifest should not be empty")
	})
	t.Run("Test Ingress v1 manifest", func(t *testing.T) {
		m, err := prometheusIngressV1(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Ingress v1 manifest should not be empty")
	})
	t.Run("Test PodMonitor manifest", func(t *testing.T) {
		crWithLabels := &monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring", Labels: map[string]string{labelKey: labelValue}},
			Spec:       monv1.PlatformMonitoringSpec{Prometheus: &monv1.Prometheus{}},
		}
		m, err := prometheusPodMonitor(crWithLabels)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "PodMonitor manifest should not be empty")
		labelsassert.AssertCRLabels(t, m.GetLabels(), utils.PrometheusComponentName, "prometheus-operator", map[string]string{labelKey: labelValue})
	})
}

func TestPrometheusReconciliationPreservesExistingPVC(t *testing.T) {
	install := true
	disabled := false
	storage := &promv1.StorageSpec{}
	tests := []struct {
		name       string
		prometheus *monv1.Prometheus
		runs       int
	}{
		{name: "missing component configuration", runs: 1},
		{name: "disabled component", prometheus: &monv1.Prometheus{Install: &disabled}, runs: 1},
		{
			name:       "removed storage configuration",
			prometheus: &monv1.Prometheus{Install: &install, Paused: true},
			runs:       1,
		},
		{
			name: "paused component",
			prometheus: &monv1.Prometheus{
				Install: &install,
				Paused:  true,
				Storage: storage,
			},
			runs: 1,
		},
		{
			name: "ordinary repeated reconciliation",
			prometheus: &monv1.Prometheus{
				Install: &install,
				Storage: storage,
			},
			runs: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reconciler, kubeClient := newPrometheusTestReconciler(t, &corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "prometheus-k8s-db-prometheus-k8s-0",
					Namespace: "monitoring",
				},
			})
			platformMonitoring := &monv1.PlatformMonitoring{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "platformmonitoring",
					Namespace: "monitoring",
					UID:       types.UID("platformmonitoring"),
				},
				Spec: monv1.PlatformMonitoringSpec{Prometheus: tt.prometheus},
			}

			for range tt.runs {
				require.NoError(t, reconciler.Run(platformMonitoring))
			}

			pvc := &corev1.PersistentVolumeClaim{}
			require.NoError(t, kubeClient.Get(
				context.Background(),
				types.NamespacedName{Name: "prometheus-k8s-db-prometheus-k8s-0", Namespace: "monitoring"},
				pvc,
			))
		})
	}
}

func newPrometheusTestReconciler(t *testing.T, objects ...client.Object) (*PrometheusReconciler, client.Client) {
	t.Helper()

	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, promv1.AddToScheme(scheme))

	kubeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).Build()
	discoveryClient := &fakediscovery.FakeDiscovery{Fake: &ktesting.Fake{}}

	return NewPrometheusReconciler(kubeClient, scheme, discoveryClient), kubeClient
}
