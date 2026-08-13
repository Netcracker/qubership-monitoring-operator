package controllers

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
)

func TestManagedPrometheusExplicitlyEnabled(t *testing.T) {
	t.Parallel()

	enabled := true
	disabled := false
	tests := []struct {
		name       string
		prometheus *monv1.Prometheus
		want       bool
	}{
		{name: "missing prometheus spec"},
		{name: "missing install value", prometheus: &monv1.Prometheus{}},
		{name: "explicitly disabled", prometheus: &monv1.Prometheus{Install: &disabled}},
		{name: "explicitly enabled", prometheus: &monv1.Prometheus{Install: &enabled}, want: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			instance := &monv1.PlatformMonitoring{
				Spec: monv1.PlatformMonitoringSpec{Prometheus: test.prometheus},
			}
			if got := managedPrometheusExplicitlyEnabled(instance); got != test.want {
				t.Fatalf("managedPrometheusExplicitlyEnabled() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestSyncPrometheusDeprecationCondition(t *testing.T) {
	t.Parallel()

	reconciler := &PlatformMonitoringReconciler{}
	unrelated := monv1.PlatformMonitoringCondition{
		Type:               "Successful",
		Status:             "True",
		Reason:             "ReconcileCycleStatus",
		Message:            "Monitoring service reconcile cycle succeeded",
		LastTransitionTime: "unchanged",
	}
	instance := &monv1.PlatformMonitoring{
		Status: monv1.PlatformMonitoringStatus{Conditions: []monv1.PlatformMonitoringCondition{unrelated}},
	}

	if changed := reconciler.syncPrometheusDeprecationCondition(instance, true); !changed {
		t.Fatal("expected enabling managed Prometheus to add the deprecation condition")
	}
	if len(instance.Status.Conditions) != 2 {
		t.Fatalf("expected two conditions, got %d", len(instance.Status.Conditions))
	}
	condition := conditionByReason(t, instance, prometheusDeprecationReason)
	if condition.Type != prometheusDeprecationType || condition.Status != prometheusDeprecationStatus ||
		condition.Message != prometheusDeprecationMessage || condition.LastTransitionTime == "" {
		t.Fatalf("unexpected deprecation condition: %#v", condition)
	}
	originalTransitionTime := condition.LastTransitionTime

	if changed := reconciler.syncPrometheusDeprecationCondition(instance, true); changed {
		t.Fatal("expected an unchanged deprecation condition to remain stable")
	}
	condition = conditionByReason(t, instance, prometheusDeprecationReason)
	if condition.LastTransitionTime != originalTransitionTime {
		t.Fatalf("transition time changed from %q to %q", originalTransitionTime, condition.LastTransitionTime)
	}

	if changed := reconciler.syncPrometheusDeprecationCondition(instance, false); !changed {
		t.Fatal("expected disabling managed Prometheus to remove the deprecation condition")
	}
	if len(instance.Status.Conditions) != 1 || instance.Status.Conditions[0] != unrelated {
		t.Fatalf("expected only the unrelated condition to remain, got %#v", instance.Status.Conditions)
	}
	if changed := reconciler.syncPrometheusDeprecationCondition(instance, false); changed {
		t.Fatal("expected removing an absent deprecation condition to be a no-op")
	}
}

func TestSyncPrometheusDeprecationConditionRepairsStaleFields(t *testing.T) {
	t.Parallel()

	reconciler := &PlatformMonitoringReconciler{}
	instance := &monv1.PlatformMonitoring{
		Status: monv1.PlatformMonitoringStatus{Conditions: []monv1.PlatformMonitoringCondition{{
			Type:               prometheusDeprecationType,
			Status:             prometheusDeprecationStatus,
			Reason:             prometheusDeprecationReason,
			Message:            "stale message",
			LastTransitionTime: "stale transition time",
		}}},
	}

	if changed := reconciler.syncPrometheusDeprecationCondition(instance, true); !changed {
		t.Fatal("expected stale deprecation fields to be repaired")
	}
	condition := conditionByReason(t, instance, prometheusDeprecationReason)
	if condition.Message != prometheusDeprecationMessage {
		t.Fatalf("message = %q, want %q", condition.Message, prometheusDeprecationMessage)
	}
	if condition.LastTransitionTime == "stale transition time" {
		t.Fatal("expected stale transition time to be replaced")
	}
}

func conditionByReason(
	t *testing.T,
	instance *monv1.PlatformMonitoring,
	reason string,
) monv1.PlatformMonitoringCondition {
	t.Helper()
	for _, condition := range instance.Status.Conditions {
		if condition.Reason == reason {
			return condition
		}
	}
	t.Fatalf("condition with reason %q was not found", reason)
	return monv1.PlatformMonitoringCondition{}
}
