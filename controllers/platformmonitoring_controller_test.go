package controllers

import (
	"context"
	"testing"

	qubershiporgv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestRequestsForGrafanaExtraVars(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, qubershiporgv1.AddToScheme(scheme))
	reconciler := &PlatformMonitoringReconciler{Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(
		&qubershiporgv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Name: "custom", Namespace: "monitoring"},
			Spec: qubershiporgv1.PlatformMonitoringSpec{Grafana: &qubershiporgv1.Grafana{
				Namespace: "grafana",
			}},
		},
		&qubershiporgv1.PlatformMonitoring{ObjectMeta: metav1.ObjectMeta{Name: "other", Namespace: "other"}},
	).Build()}

	requests := reconciler.requestsForPlatformMonitorings(context.Background(), &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "grafana-extra-vars", Namespace: "grafana"},
	})

	assert.Equal(t, []reconcile.Request{{NamespacedName: types.NamespacedName{
		Name: "custom", Namespace: "monitoring",
	}}}, requests)
}

func TestGrafanaExtraVarsResourcePredicates(t *testing.T) {
	assert.True(t, isGrafanaExtraVarsConfigMap(&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "grafana-extra-vars"}}))
	assert.False(t, isGrafanaExtraVarsConfigMap(&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "grafana-extra-vars-secret"}}))
	assert.True(t, isGrafanaExtraVarsSecret(&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "grafana-extra-vars-secret"}}))
	assert.False(t, isGrafanaExtraVarsSecret(&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "grafana-extra-vars"}}))
}
