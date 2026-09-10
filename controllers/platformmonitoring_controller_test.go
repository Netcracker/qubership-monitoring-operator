package controllers

import (
	"context"
	"testing"
	"time"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	"github.com/go-logr/logr"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	v1beta1ext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery/fake"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	k8stesting "k8s.io/client-go/testing"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	clientfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestRequestsForGrafanaExtraVars(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	reconciler := &PlatformMonitoringReconciler{Client: clientfake.NewClientBuilder().WithScheme(scheme).WithObjects(
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
		Client: clientfake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
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

type countingStatusWriter struct {
	client.SubResourceWriter
	updates int
}

func (w *countingStatusWriter) Update(
	ctx context.Context,
	obj client.Object,
	opts ...client.SubResourceUpdateOption,
) error {
	w.updates++
	return w.SubResourceWriter.Update(ctx, obj, opts...)
}

type statusCountingClient struct {
	client.Client
	statusWriter *countingStatusWriter
}

func newStatusCountingClient(delegate client.Client) *statusCountingClient {
	return &statusCountingClient{
		Client: delegate,
		statusWriter: &countingStatusWriter{
			SubResourceWriter: delegate.Status(),
		},
	}
}

func (c *statusCountingClient) Status() client.SubResourceWriter {
	return c.statusWriter
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

func TestReconcileStatusWrites(t *testing.T) {
	const transitionTime = "2026-08-11 12:00:00 +0000 UTC"

	tests := []struct {
		name               string
		generation         int64
		observedGeneration int64
		wantUpdates        int
	}{
		{
			name:               "idle observed generation",
			generation:         1,
			observedGeneration: 1,
			wantUpdates:        0,
		},
		{
			name:               "new generation",
			generation:         2,
			observedGeneration: 1,
			wantUpdates:        2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := &monv1.PlatformMonitoring{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "monitoring",
					Namespace:  "monitoring",
					Generation: tt.generation,
				},
				Status: monv1.PlatformMonitoringStatus{
					ObservedGeneration: tt.observedGeneration,
					Conditions: []monv1.PlatformMonitoringCondition{{
						Type:               "Successful",
						Status:             "True",
						Reason:             "ReconcileCycleStatus",
						Message:            "Monitoring service reconcile cycle succeeded",
						LastTransitionTime: transitionTime,
					}},
				},
			}
			reconciler, countingClient := newStatusTestReconciler(t, resource)

			result, err := reconciler.Reconcile(
				context.Background(),
				ctrl.Request{NamespacedName: types.NamespacedName{Name: resource.Name, Namespace: resource.Namespace}},
			)

			require.NoError(t, err)
			assert.Equal(t, 60*time.Second, result.RequeueAfter)
			assert.Equal(t, tt.wantUpdates, countingClient.statusWriter.updates)

			stored := &monv1.PlatformMonitoring{}
			require.NoError(t, countingClient.Get(
				context.Background(),
				client.ObjectKeyFromObject(resource),
				stored,
			))
			assert.Equal(t, tt.generation, stored.Status.ObservedGeneration)
			assert.Equal(t, "Successful", stored.Status.Conditions[0].Type)
			if tt.wantUpdates == 0 {
				assert.Equal(t, transitionTime, stored.Status.Conditions[0].LastTransitionTime)
			}
		})
	}
}

func TestReconcileStatusRecoversFromPriorFailure(t *testing.T) {
	resource := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "monitoring",
			Namespace:  "monitoring",
			Generation: 1,
		},
		Status: monv1.PlatformMonitoringStatus{
			ObservedGeneration: 1,
			Conditions: []monv1.PlatformMonitoringCondition{
				{
					Type:               "Failed",
					Status:             "False",
					Reason:             "ReconcilePrometheusStatus",
					Message:            "Prometheus reconcile cycle failed",
					LastTransitionTime: "2026-08-11 12:00:00 +0000 UTC",
				},
				{
					Type:               "Failed",
					Status:             "False",
					Reason:             "ReconcileCycleStatus",
					Message:            "Monitoring service reconcile cycle failed",
					LastTransitionTime: "2026-08-11 12:00:00 +0000 UTC",
				},
			},
		},
	}
	reconciler, countingClient := newStatusTestReconciler(t, resource)

	result, err := reconciler.Reconcile(
		context.Background(),
		ctrl.Request{NamespacedName: client.ObjectKeyFromObject(resource)},
	)

	require.NoError(t, err)
	assert.Equal(t, 60*time.Second, result.RequeueAfter)
	assert.Equal(t, 1, countingClient.statusWriter.updates)

	stored := &monv1.PlatformMonitoring{}
	require.NoError(t, countingClient.Get(context.Background(), client.ObjectKeyFromObject(resource), stored))
	require.Len(t, stored.Status.Conditions, 1)
	assert.Equal(t, "ReconcileCycleStatus", stored.Status.Conditions[0].Reason)
	assert.Equal(t, "Successful", stored.Status.Conditions[0].Type)
	assert.Equal(t, "True", stored.Status.Conditions[0].Status)
}

func TestReconcileInvalidIntervalRunsComponents(t *testing.T) {
	resource := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "monitoring",
			Namespace:  "monitoring",
			Generation: 1,
		},
		Spec: monv1.PlatformMonitoringSpec{
			KubernetesMonitors: map[string]monv1.Monitor{
				"apiserverServiceMonitor": {Install: ptr.To(false)},
			},
		},
		Status: monv1.PlatformMonitoringStatus{
			ObservedGeneration: 1,
			Conditions: []monv1.PlatformMonitoringCondition{{
				Type:               "Successful",
				Status:             "True",
				Reason:             "ReconcileCycleStatus",
				Message:            "Monitoring service reconcile cycle succeeded",
				LastTransitionTime: "2026-08-11 12:00:00 +0000 UTC",
			}},
		},
	}
	staleServiceMonitor := &promv1.ServiceMonitor{ObjectMeta: metav1.ObjectMeta{
		Name:      "monitoring-kube-apiserver-service-monitor",
		Namespace: "monitoring",
	}}
	reconciler, countingClient := newStatusTestReconciler(t, resource, staleServiceMonitor)
	t.Setenv("RECONCILIATION_INTERVAL", "invalid")

	_, err := reconciler.Reconcile(
		context.Background(),
		ctrl.Request{NamespacedName: client.ObjectKeyFromObject(resource)},
	)

	require.Error(t, err)
	assert.ErrorContains(t, err, "invalid syntax")
	assert.Equal(t, 1, countingClient.statusWriter.updates)

	deleted := &promv1.ServiceMonitor{}
	err = countingClient.Get(context.Background(), client.ObjectKeyFromObject(staleServiceMonitor), deleted)
	assert.True(t, apierrors.IsNotFound(err))

	stored := &monv1.PlatformMonitoring{}
	require.NoError(t, countingClient.Get(context.Background(), client.ObjectKeyFromObject(resource), stored))
	assert.Equal(t, int64(1), stored.Status.ObservedGeneration)
	require.Len(t, stored.Status.Conditions, 1)
	assert.Equal(t, "In progress", stored.Status.Conditions[0].Type)
}

func newStatusTestReconciler(
	t *testing.T,
	objects ...client.Object,
) (*PlatformMonitoringReconciler, *statusCountingClient) {
	t.Helper()
	t.Setenv("LOG_LEVEL", "error")

	scheme := runtime.NewScheme()
	require.NoError(t, clientgoscheme.AddToScheme(scheme))
	require.NoError(t, monv1.AddToScheme(scheme))
	require.NoError(t, promv1.AddToScheme(scheme))
	require.NoError(t, grafv1.AddToScheme(scheme))
	require.NoError(t, vmetricsv1b1.AddToScheme(scheme))
	require.NoError(t, v1beta1ext.AddToScheme(scheme))
	baseClient := clientfake.NewClientBuilder().
		WithScheme(scheme).
		WithStatusSubresource(&monv1.PlatformMonitoring{}).
		WithObjects(objects...).
		Build()
	countingClient := newStatusCountingClient(baseClient)

	return &PlatformMonitoringReconciler{
		Client: countingClient,
		Log:    logr.Discard(),
		Scheme: scheme,
		Config: &rest.Config{},
		DiscoveryClient: &fake.FakeDiscovery{
			Fake:               &k8stesting.Fake{},
			FakedServerVersion: &version.Info{Major: "1", Minor: "30"},
		},
	}, countingClient
}
