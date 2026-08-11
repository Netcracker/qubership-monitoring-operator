package controllers

import (
	"context"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/go-logr/logr"
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
	require.NoError(t, monv1.AddToScheme(scheme))
	reconciler := &PlatformMonitoringReconciler{Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(
		&monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Name: "custom", Namespace: "monitoring"},
			Spec: monv1.PlatformMonitoringSpec{Grafana: &monv1.Grafana{
				Namespace: "grafana",
			}},
		},
		&monv1.PlatformMonitoring{ObjectMeta: metav1.ObjectMeta{Name: "other", Namespace: "other"}},
	).Build()}

	requests := reconciler.requestsForPlatformMonitorings(context.Background(), &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "grafana-extra-vars", Namespace: "grafana"},
	})

	assert.Equal(t, []reconcile.Request{{NamespacedName: types.NamespacedName{
		Name: "custom", Namespace: "monitoring",
	}}}, requests)
}

func TestRequestsForGrafanaExtraVarsReturnsNoRequestsWhenListFails(t *testing.T) {
	reconciler := &PlatformMonitoringReconciler{
		Client: fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
		Log:    logr.Discard(),
	}

	requests := reconciler.requestsForPlatformMonitorings(context.Background(), &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "grafana-extra-vars", Namespace: "grafana"},
	})

	assert.Empty(t, requests)
}

func TestGrafanaExtraVarsResourcePredicates(t *testing.T) {
	assert.True(t, isGrafanaExtraVarsConfigMap(&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "grafana-extra-vars"}}))
	assert.False(t, isGrafanaExtraVarsConfigMap(&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "grafana-extra-vars-secret"}}))
	assert.True(t, isGrafanaExtraVarsSecret(&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "grafana-extra-vars-secret"}}))
	assert.False(t, isGrafanaExtraVarsSecret(&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "grafana-extra-vars"}}))
}

func TestPrepareStatusForUpdatePreservesTransitionTime(t *testing.T) {
	const transitionTime = "2026-08-11 12:00:00 +0000 UTC"

	tests := []struct {
		name        string
		oldMessage  string
		newMessage  string
		wantChanged bool
	}{
		{
			name:        "unchanged condition",
			oldMessage:  "Monitoring service reconcile cycle succeeded",
			newMessage:  "Monitoring service reconcile cycle succeeded",
			wantChanged: false,
		},
		{
			name:        "message-only update",
			oldMessage:  "Old diagnostic details",
			newMessage:  "New diagnostic details",
			wantChanged: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := &monv1.PlatformMonitoring{
				Status: monv1.PlatformMonitoringStatus{
					Conditions: []monv1.PlatformMonitoringCondition{{
						Type:               "Successful",
						Status:             "True",
						Reason:             "ReconcileCycleStatus",
						Message:            tt.oldMessage,
						LastTransitionTime: transitionTime,
					}},
				},
			}
			reconciler := &PlatformMonitoringReconciler{}

			changed := reconciler.prepareStatusForUpdate(
				resource,
				"Successful",
				"True",
				"ReconcileCycleStatus",
				tt.newMessage,
			)

			assert.Equal(t, tt.wantChanged, changed)
			assert.Equal(t, tt.newMessage, resource.Status.Conditions[0].Message)
			assert.Equal(t, transitionTime, resource.Status.Conditions[0].LastTransitionTime)
		})
	}
}

func TestPrepareStatusForUpdateAdvancesTransitionTime(t *testing.T) {
	const transitionTime = "2026-08-11 12:00:00 +0000 UTC"

	tests := []struct {
		name      string
		oldType   string
		newType   string
		oldStatus string
		newStatus string
		oldTime   string
	}{
		{
			name:      "type transition",
			oldType:   "Successful",
			newType:   "Failed",
			oldStatus: "True",
			newStatus: "True",
			oldTime:   transitionTime,
		},
		{
			name:      "status transition",
			oldType:   "Successful",
			newType:   "Successful",
			oldStatus: "False",
			newStatus: "True",
			oldTime:   transitionTime,
		},
		{
			name:      "missing transition time",
			oldType:   "Successful",
			newType:   "Successful",
			oldStatus: "True",
			newStatus: "True",
			oldTime:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := &monv1.PlatformMonitoring{
				Status: monv1.PlatformMonitoringStatus{
					Conditions: []monv1.PlatformMonitoringCondition{{
						Type:               tt.oldType,
						Status:             tt.oldStatus,
						Reason:             "ReconcileCycleStatus",
						Message:            "Monitoring service reconcile cycle completed",
						LastTransitionTime: tt.oldTime,
					}},
				},
			}
			reconciler := &PlatformMonitoringReconciler{}

			changed := reconciler.prepareStatusForUpdate(
				resource,
				tt.newType,
				tt.newStatus,
				"ReconcileCycleStatus",
				"Monitoring service reconcile cycle completed",
			)

			assert.True(t, changed)
			assert.NotEmpty(t, resource.Status.Conditions[0].LastTransitionTime)
			assert.NotEqual(t, tt.oldTime, resource.Status.Conditions[0].LastTransitionTime)
		})
	}
}

func TestHasFailedComponentConditionIgnoresCycleCondition(t *testing.T) {
	tests := []struct {
		name       string
		conditions []monv1.PlatformMonitoringCondition
		want       bool
	}{
		{
			name: "aggregate failure only",
			conditions: []monv1.PlatformMonitoringCondition{{
				Type:   "Failed",
				Reason: "ReconcileCycleStatus",
			}},
			want: false,
		},
		{
			name: "component failure",
			conditions: []monv1.PlatformMonitoringCondition{{
				Type:   "Failed",
				Reason: "ReconcilePrometheusStatus",
			}},
			want: true,
		},
		{
			name: "successful component condition",
			conditions: []monv1.PlatformMonitoringCondition{{
				Type:   "Successful",
				Reason: "ReconcilePrometheusStatus",
			}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, hasFailedComponentCondition(tt.conditions))
		})
	}
}
