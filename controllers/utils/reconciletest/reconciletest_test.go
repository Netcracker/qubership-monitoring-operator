package reconciletest

import (
	"context"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestUpdateCounterCountsUpdates(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	counter := &UpdateCounter{}
	configMap := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "cm", Namespace: "monitoring"}}
	c := NewClient(scheme, counter, configMap)

	got := &corev1.ConfigMap{}
	require.NoError(t, c.Get(context.TODO(), client.ObjectKeyFromObject(configMap), got))
	got.Data = map[string]string{"k": "v"}
	require.NoError(t, c.Update(context.TODO(), got))
	assert.Equal(t, int64(1), counter.Updates())

	counter.Reset()
	assert.Zero(t, counter.Updates())
}

func TestAssertIdleReconcileSkipsUpdates(t *testing.T) {
	counter := &UpdateCounter{}
	c := NewClient(NewScheme(t), counter)
	calls := 0
	createOnce := func(cr *monv1.PlatformMonitoring) error {
		calls++
		configMap := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "cm", Namespace: cr.GetNamespace()}}
		err := c.Get(context.TODO(), client.ObjectKeyFromObject(configMap), configMap)
		if apierrors.IsNotFound(err) {
			return c.Create(context.TODO(), configMap)
		}
		return err
	}
	cr := &monv1.PlatformMonitoring{ObjectMeta: metav1.ObjectMeta{Name: "pm", Namespace: "monitoring"}}

	AssertIdleReconcileSkipsUpdates(t, cr, counter, map[string]Handler{"ConfigMap": createOnce})

	assert.Equal(t, 2, calls)
}
