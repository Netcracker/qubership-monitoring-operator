package utils

import (
	"context"
	"errors"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

const createOrUpdateTestNamespace = "monitoring"

var errCreateOrUpdateTest = errors.New("injected failure")

func TestCreateOrUpdate(t *testing.T) {
	tests := []struct {
		name        string
		live        []client.Object
		desired     map[string]string
		wantLabels  map[string]string
		wantUpdated bool
	}{
		{
			name:       "creates a missing object",
			desired:    map[string]string{"app": "new"},
			wantLabels: map[string]string{"app": "new"},
		},
		{
			name:       "skips an unchanged object",
			live:       []client.Object{liveServiceAccount("sa", map[string]string{"app": "same"})},
			desired:    map[string]string{"app": "same"},
			wantLabels: map[string]string{"app": "same"},
		},
		{
			name:        "updates a changed object",
			live:        []client.Object{liveServiceAccount("sa", map[string]string{"app": "old"})},
			desired:     map[string]string{"app": "new"},
			wantLabels:  map[string]string{"app": "new"},
			wantUpdated: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newCreateOrUpdateTestReconciler(t, interceptor.Funcs{}, tt.live...)
			desired := testServiceAccount("sa", tt.desired)
			existing := &corev1.ServiceAccount{ObjectMeta: *desired.ObjectMeta.DeepCopy()}

			require.NoError(t, CreateOrUpdate(r, testPlatformMonitoring(), desired, existing, SyncLabels))

			got := &corev1.ServiceAccount{}
			require.NoError(t, r.Client.Get(context.TODO(), client.ObjectKeyFromObject(desired), got))
			assert.Equal(t, tt.wantLabels, got.GetLabels())
			if len(tt.live) > 0 {
				assert.Equal(t, tt.wantUpdated, got.ResourceVersion != tt.live[0].GetResourceVersion())
			}
		})
	}
}

func TestCreateOrUpdateReturnsReadError(t *testing.T) {
	r := newCreateOrUpdateTestReconciler(t, interceptor.Funcs{Get: failingGet})
	desired := testServiceAccount("sa", nil)
	existing := &corev1.ServiceAccount{ObjectMeta: desired.ObjectMeta}

	err := CreateOrUpdate(r, testPlatformMonitoring(), desired, existing, SyncLabels)
	assert.ErrorIs(t, err, errCreateOrUpdateTest)
}

func TestCreateOrUpdateWorkload(t *testing.T) {
	tests := []struct {
		name         string
		liveSelector map[string]string
		liveLabels   map[string]string
		wantReplaced bool
		wantUpdated  bool
	}{
		{name: "skips an unchanged workload", liveSelector: map[string]string{"app": "x"}, liveLabels: map[string]string{"app": "x"}},
		{name: "updates changed fields", liveSelector: map[string]string{"app": "x"}, liveLabels: map[string]string{"app": "old"}, wantUpdated: true},
		{name: "recreates a workload with a different selector", liveSelector: map[string]string{"app": "old"}, liveLabels: map[string]string{"app": "x"}, wantReplaced: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			live := testDeployment(tt.liveSelector, tt.liveLabels)
			live.UID = types.UID("live-uid")
			live.ResourceVersion = "1"
			r := newCreateOrUpdateTestReconciler(t, interceptor.Funcs{}, live)
			desired := testDeployment(map[string]string{"app": "x"}, map[string]string{"app": "x"})
			existing := &appsv1.Deployment{ObjectMeta: *desired.ObjectMeta.DeepCopy()}

			require.NoError(t, CreateOrUpdateWorkload(r, testPlatformMonitoring(), desired, existing, syncTestDeployment))

			got := &appsv1.Deployment{}
			require.NoError(t, r.Client.Get(context.TODO(), client.ObjectKeyFromObject(desired), got))
			assert.Equal(t, desired.Spec.Selector, got.Spec.Selector)
			assert.Equal(t, desired.GetLabels(), got.GetLabels())
			assert.Equal(t, tt.wantReplaced, got.UID != live.UID, "replaced")
			assert.Equal(t, tt.wantUpdated, !tt.wantReplaced && got.ResourceVersion != "1", "updated")
		})
	}
}

func TestCreateOrUpdateWorkloadCreatesMissingWorkload(t *testing.T) {
	r := newCreateOrUpdateTestReconciler(t, interceptor.Funcs{})
	desired := testDeployment(map[string]string{"app": "x"}, nil)
	existing := &appsv1.Deployment{ObjectMeta: *desired.ObjectMeta.DeepCopy()}

	require.NoError(t, CreateOrUpdateWorkload(r, testPlatformMonitoring(), desired, existing, syncTestDeployment))
	require.NoError(t, r.Client.Get(context.TODO(), client.ObjectKeyFromObject(desired), &appsv1.Deployment{}))
}

func TestCreateOrUpdateWorkloadReturnsErrors(t *testing.T) {
	tests := []struct {
		name  string
		funcs interceptor.Funcs
	}{
		{name: "read", funcs: interceptor.Funcs{Get: failingGet}},
		{name: "delete", funcs: interceptor.Funcs{Delete: func(context.Context, client.WithWatch, client.Object, ...client.DeleteOption) error {
			return errCreateOrUpdateTest
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newCreateOrUpdateTestReconciler(t, tt.funcs, testDeployment(map[string]string{"app": "old"}, nil))
			desired := testDeployment(map[string]string{"app": "x"}, nil)
			existing := &appsv1.Deployment{ObjectMeta: *desired.ObjectMeta.DeepCopy()}

			err := CreateOrUpdateWorkload(r, testPlatformMonitoring(), desired, existing, syncTestDeployment)
			assert.ErrorIs(t, err, errCreateOrUpdateTest)
		})
	}
}

func TestCreateOrUpdateReplacingLegacy(t *testing.T) {
	tests := []struct {
		name       string
		live       []client.Object
		wantLabels map[string]string
	}{
		{name: "creates when neither object exists", wantLabels: map[string]string{"app": "new"}},
		{
			name:       "deletes the legacy object before it creates",
			live:       []client.Object{testServiceAccount("legacy", map[string]string{"app": "legacy"})},
			wantLabels: map[string]string{"app": "new"},
		},
		{
			name:       "updates the current object and keeps the legacy one",
			live:       []client.Object{testServiceAccount("sa", map[string]string{"app": "old"}), testServiceAccount("legacy", nil)},
			wantLabels: map[string]string{"app": "new"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newCreateOrUpdateTestReconciler(t, interceptor.Funcs{}, tt.live...)
			desired := testServiceAccount("sa", map[string]string{"app": "new"})
			existing := &corev1.ServiceAccount{ObjectMeta: *desired.ObjectMeta.DeepCopy()}
			legacy := &corev1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: "legacy", Namespace: createOrUpdateTestNamespace}}

			require.NoError(t, CreateOrUpdateReplacingLegacy(r, testPlatformMonitoring(), desired, existing, legacy, SyncLabels))

			got := &corev1.ServiceAccount{}
			require.NoError(t, r.Client.Get(context.TODO(), client.ObjectKeyFromObject(desired), got))
			assert.Equal(t, tt.wantLabels, got.GetLabels())
			legacyErr := r.Client.Get(context.TODO(), client.ObjectKeyFromObject(legacy), &corev1.ServiceAccount{})
			wantLegacy := len(tt.live) == 2
			assert.Equal(t, wantLegacy, legacyErr == nil, "legacy object present")
		})
	}
}

func TestCreateOrUpdateReplacingLegacyReturnsErrors(t *testing.T) {
	failingDelete := func(context.Context, client.WithWatch, client.Object, ...client.DeleteOption) error {
		return errCreateOrUpdateTest
	}
	tests := []struct {
		name  string
		funcs interceptor.Funcs
	}{
		{name: "read", funcs: interceptor.Funcs{Get: failingGet}},
		{name: "legacy delete", funcs: interceptor.Funcs{Delete: failingDelete}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newCreateOrUpdateTestReconciler(t, tt.funcs, testServiceAccount("legacy", nil))
			desired := testServiceAccount("sa", nil)
			existing := &corev1.ServiceAccount{ObjectMeta: *desired.ObjectMeta.DeepCopy()}
			legacy := &corev1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: "legacy", Namespace: createOrUpdateTestNamespace}}

			err := CreateOrUpdateReplacingLegacy(r, testPlatformMonitoring(), desired, existing, legacy, SyncLabels)
			assert.ErrorIs(t, err, errCreateOrUpdateTest)
		})
	}
}

func TestUpdateIfChanged(t *testing.T) {
	live := testServiceAccount("sa", map[string]string{"app": "old"})
	r := newCreateOrUpdateTestReconciler(t, interceptor.Funcs{}, live)
	obj := &corev1.ServiceAccount{}
	require.NoError(t, r.Client.Get(context.TODO(), client.ObjectKeyFromObject(live), obj))
	obj.SetLabels(map[string]string{"app": "new"})

	require.NoError(t, r.UpdateIfChanged(obj, false))
	got := &corev1.ServiceAccount{}
	require.NoError(t, r.Client.Get(context.TODO(), client.ObjectKeyFromObject(live), got))
	assert.Equal(t, map[string]string{"app": "old"}, got.GetLabels(), "no Update when unchanged")

	require.NoError(t, r.UpdateIfChanged(obj, true))
	require.NoError(t, r.Client.Get(context.TODO(), client.ObjectKeyFromObject(live), got))
	assert.Equal(t, map[string]string{"app": "new"}, got.GetLabels())
}

func syncTestDeployment(existing, desired *appsv1.Deployment) bool {
	return SyncLabels(existing, desired)
}

func failingGet(context.Context, client.WithWatch, client.ObjectKey, client.Object, ...client.GetOption) error {
	return errCreateOrUpdateTest
}

func newCreateOrUpdateTestReconciler(t *testing.T, funcs interceptor.Funcs, objs ...client.Object) *ComponentReconciler {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, appsv1.AddToScheme(scheme))
	return &ComponentReconciler{
		Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).WithInterceptorFuncs(funcs).Build(),
		Scheme: scheme,
		Log:    Logger("create_or_update_test"),
	}
}

func testPlatformMonitoring() *monv1.PlatformMonitoring {
	return &monv1.PlatformMonitoring{ObjectMeta: metav1.ObjectMeta{
		Name:      "platformmonitoring",
		Namespace: createOrUpdateTestNamespace,
		UID:       types.UID("platform-monitoring-uid"),
	}}
}

func testServiceAccount(name string, labels map[string]string) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{
		Name:      name,
		Namespace: createOrUpdateTestNamespace,
		Labels:    labels,
	}}
}

func liveServiceAccount(name string, labels map[string]string) *corev1.ServiceAccount {
	sa := testServiceAccount(name, labels)
	sa.ResourceVersion = "1"
	return sa
}

func testDeployment(selector, labels map[string]string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "workload", Namespace: createOrUpdateTestNamespace, Labels: labels},
		Spec:       appsv1.DeploymentSpec{Selector: &metav1.LabelSelector{MatchLabels: selector}},
	}
}
