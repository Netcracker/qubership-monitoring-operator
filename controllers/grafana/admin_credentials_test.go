package grafana

import (
	"context"
	"testing"
	"time"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/go-logr/logr"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/kubernetes"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const testNamespace = "monitoring"

func adminSecret(namespace, name, password string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: name + "-admin-credentials", Namespace: namespace},
		Data: map[string][]byte{
			"GF_SECURITY_ADMIN_USER":     []byte("admin"),
			"GF_SECURITY_ADMIN_PASSWORD": []byte(password),
		},
	}
}

func adminCredentialsCR(grafana *monv1.Grafana) *monv1.PlatformMonitoring {
	if grafana == nil {
		grafana = &monv1.Grafana{}
	}
	return &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Name: "monitoring", Namespace: testNamespace},
		Spec:       monv1.PlatformMonitoringSpec{Grafana: grafana},
	}
}

func newAdminCredentialsReconciler(t *testing.T, kubeClient kubernetes.Interface, objects ...client.Object) *GrafanaReconciler {
	t.Helper()
	testScheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(testScheme))
	require.NoError(t, monv1.AddToScheme(testScheme))
	require.NoError(t, grafv1.AddToScheme(testScheme))
	if kubeClient == nil {
		kubeClient = kubernetesfake.NewSimpleClientset()
	}
	return &GrafanaReconciler{
		KubeClient: kubeClient,
		ComponentReconciler: &utils.ComponentReconciler{
			Client: fake.NewClientBuilder().WithScheme(testScheme).WithObjects(objects...).Build(),
			Scheme: testScheme,
			Dc:     &fakediscovery.FakeDiscovery{Fake: &k8stesting.Fake{}},
			Log:    logr.Discard(),
		},
	}
}

func shortenGrafanaPodWait(t *testing.T) {
	t.Helper()
	timeout, interval := grafanaPodReadyTimeout, grafanaPodPollInterval
	grafanaPodReadyTimeout, grafanaPodPollInterval = 20*time.Millisecond, 5*time.Millisecond
	t.Cleanup(func() { grafanaPodReadyTimeout, grafanaPodPollInterval = timeout, interval })
}

func existingGrafana(t *testing.T, cr *monv1.PlatformMonitoring, syncedChecksum string) *grafv1.Grafana {
	t.Helper()
	m, err := grafanaWithDefaultSources(cr)
	require.NoError(t, err)
	if syncedChecksum != "" {
		m.SetAnnotations(map[string]string{adminPasswordSyncedAnnotation: syncedChecksum})
	}
	return m
}

func readyPod(name string, ready bool) *corev1.Pod {
	status := corev1.ConditionFalse
	if ready {
		status = corev1.ConditionTrue
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: name, Namespace: testNamespace,
			Labels: map[string]string{"app.kubernetes.io/name": utils.GrafanaComponentName},
		},
		Status: corev1.PodStatus{
			Phase:      corev1.PodRunning,
			Conditions: []corev1.PodCondition{{Type: corev1.PodReady, Status: status}},
		},
	}
}

func TestAdminPasswordChecksum(t *testing.T) {
	first := adminPasswordChecksum(adminSecret(testNamespace, "grafana", "first"))
	assert.Len(t, first, 64)
	assert.Equal(t, first, adminPasswordChecksum(adminSecret(testNamespace, "grafana", "first")))
	assert.NotEqual(t, first, adminPasswordChecksum(adminSecret(testNamespace, "grafana", "second")))

	userOnly := adminSecret(testNamespace, "grafana", "first")
	userOnly.Data["GF_SECURITY_ADMIN_USER"] = []byte("other")
	assert.Equal(t, first, adminPasswordChecksum(userOnly), "the user name must not affect the checksum")

	assert.Empty(t, adminPasswordChecksum(adminSecret(testNamespace, "grafana", "")))
	assert.Empty(t, adminPasswordChecksum(&corev1.Secret{}))
}

func TestReadAdminPasswordChecksumUsesCustomNameAndNamespace(t *testing.T) {
	cr := adminCredentialsCR(&monv1.Grafana{Name: "custom", Namespace: "grafana-ns"})
	secret := adminSecret("grafana-ns", "custom", "secret")
	r := newAdminCredentialsReconciler(t, nil, secret)

	checksum, err := r.readAdminPasswordChecksum(cr)
	require.NoError(t, err)
	assert.Equal(t, adminPasswordChecksum(secret), checksum)
}

func TestReadAdminPasswordChecksumMissingSecret(t *testing.T) {
	r := newAdminCredentialsReconciler(t, nil)

	checksum, err := r.readAdminPasswordChecksum(adminCredentialsCR(nil))
	require.NoError(t, err)
	assert.Empty(t, checksum)
}

func TestReadAdminPasswordChecksumReturnsReadError(t *testing.T) {
	r := newAdminCredentialsReconciler(t, nil)
	r.Client = fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build()

	_, err := r.readAdminPasswordChecksum(adminCredentialsCR(nil))
	assert.ErrorContains(t, err, "cannot read Grafana admin credentials Secret monitoring/grafana-admin-credentials")
}

func TestAddAdminPasswordChecksumOnCreate(t *testing.T) {
	userAnnotations := map[string]string{"example.com/user": "value"}
	cr := adminCredentialsCR(&monv1.Grafana{Annotations: userAnnotations})
	secret := adminSecret(testNamespace, "grafana", "secret")
	r := newAdminCredentialsReconciler(t, nil, secret)
	m, err := grafanaWithDefaultSources(cr)
	require.NoError(t, err)

	require.NoError(t, r.addAdminPasswordChecksum(cr, m, nil, true))

	checksum := adminPasswordChecksum(secret)
	assert.Equal(t, checksum, m.Spec.Deployment.Spec.Template.Annotations[adminSecretChecksumAnnotation])
	assert.Equal(t, checksum, m.GetAnnotations()[adminPasswordSyncedAnnotation])
	assert.Equal(t, "value", m.GetAnnotations()["example.com/user"])
	assert.Equal(t, map[string]string{"example.com/user": "value"}, userAnnotations,
		"the PlatformMonitoring annotations must not be modified")
}

func TestAddAdminPasswordChecksumOnUpdateDoesNotMarkSynced(t *testing.T) {
	cr := adminCredentialsCR(nil)
	secret := adminSecret(testNamespace, "grafana", "secret")
	r := newAdminCredentialsReconciler(t, nil, secret)
	m, err := grafanaWithDefaultSources(cr)
	require.NoError(t, err)

	require.NoError(t, r.addAdminPasswordChecksum(cr, m, map[string]string{adminSecretChecksumAnnotation: "old"}, false))

	assert.Equal(t, adminPasswordChecksum(secret), m.Spec.Deployment.Spec.Template.Annotations[adminSecretChecksumAnnotation])
	assert.NotContains(t, m.GetAnnotations(), adminPasswordSyncedAnnotation)
}

func TestAddAdminPasswordChecksumKeepsExistingChecksumWithoutPassword(t *testing.T) {
	cr := adminCredentialsCR(nil)
	r := newAdminCredentialsReconciler(t, nil)
	m, err := grafanaWithDefaultSources(cr)
	require.NoError(t, err)

	require.NoError(t, r.addAdminPasswordChecksum(cr, m, map[string]string{adminSecretChecksumAnnotation: "old"}, false))
	assert.Equal(t, "old", m.Spec.Deployment.Spec.Template.Annotations[adminSecretChecksumAnnotation])

	fresh, err := grafanaWithDefaultSources(cr)
	require.NoError(t, err)
	require.NoError(t, r.addAdminPasswordChecksum(cr, fresh, nil, true))
	assert.NotContains(t, fresh.Spec.Deployment.Spec.Template.Annotations, adminSecretChecksumAnnotation)
	assert.NotContains(t, fresh.GetAnnotations(), adminPasswordSyncedAnnotation)
}

func TestHandleGrafanaCreatesGrafanaWithSyncedPasswordChecksum(t *testing.T) {
	cr := adminCredentialsCR(nil)
	secret := adminSecret(testNamespace, "grafana", "secret")
	r := newAdminCredentialsReconciler(t, nil, secret)

	require.NoError(t, r.handleGrafana(cr))

	created := &grafv1.Grafana{}
	require.NoError(t, r.Client.Get(context.Background(), types.NamespacedName{Name: "grafana", Namespace: testNamespace}, created))
	checksum := adminPasswordChecksum(secret)
	assert.Equal(t, checksum, created.Spec.Deployment.Spec.Template.Annotations[adminSecretChecksumAnnotation])
	assert.Equal(t, checksum, created.GetAnnotations()[adminPasswordSyncedAnnotation])
}

func TestSyncGrafanaAdminPasswordSkipsWithoutPassword(t *testing.T) {
	r := newAdminCredentialsReconciler(t, nil)
	assert.NoError(t, r.syncGrafanaAdminPassword(adminCredentialsCR(nil)))
}

func TestSyncGrafanaAdminPasswordSkipsMissingGrafana(t *testing.T) {
	r := newAdminCredentialsReconciler(t, nil, adminSecret(testNamespace, "grafana", "secret"))
	assert.NoError(t, r.syncGrafanaAdminPassword(adminCredentialsCR(nil)))
}

func TestSyncGrafanaAdminPasswordSkipsSyncedPassword(t *testing.T) {
	cr := adminCredentialsCR(nil)
	secret := adminSecret(testNamespace, "grafana", "secret")
	kubeClient := kubernetesfake.NewSimpleClientset()
	r := newAdminCredentialsReconciler(t, kubeClient, secret, existingGrafana(t, cr, adminPasswordChecksum(secret)))

	require.NoError(t, r.syncGrafanaAdminPassword(cr))
	assert.Empty(t, kubeClient.Actions(), "no Grafana pod lookup is expected for a synced password")
}

func TestSyncGrafanaAdminPasswordKeepsChecksumWhenResetFails(t *testing.T) {
	shortenGrafanaPodWait(t)
	cr := adminCredentialsCR(nil)
	secret := adminSecret(testNamespace, "grafana", "new")
	r := newAdminCredentialsReconciler(t, nil, secret, existingGrafana(t, cr, "old"))

	err := r.syncGrafanaAdminPassword(cr)
	assert.ErrorContains(t, err, "no ready Grafana pod")

	current := &grafv1.Grafana{}
	require.NoError(t, r.Client.Get(context.Background(), types.NamespacedName{Name: "grafana", Namespace: testNamespace}, current))
	assert.Equal(t, "old", current.GetAnnotations()[adminPasswordSyncedAnnotation],
		"a failed reset must be retried on the next reconcile")
}

func TestSyncGrafanaAdminPasswordReturnsGrafanaReadError(t *testing.T) {
	r := newAdminCredentialsReconciler(t, nil, adminSecret(testNamespace, "grafana", "secret"))
	testScheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(testScheme))
	r.Client = fake.NewClientBuilder().WithScheme(testScheme).
		WithObjects(adminSecret(testNamespace, "grafana", "secret")).Build()

	assert.Error(t, r.syncGrafanaAdminPassword(adminCredentialsCR(nil)))
}

func TestResetGrafanaCredentialsSkipsMissingSecret(t *testing.T) {
	kubeClient := kubernetesfake.NewSimpleClientset()
	r := newAdminCredentialsReconciler(t, kubeClient)

	require.NoError(t, r.resetGrafanaCredentials(adminCredentialsCR(nil)))
	assert.Empty(t, kubeClient.Actions())
}

func TestResetGrafanaCredentialsSkipsEmptyPassword(t *testing.T) {
	kubeClient := kubernetesfake.NewSimpleClientset()
	r := newAdminCredentialsReconciler(t, kubeClient, adminSecret(testNamespace, "grafana", ""))

	require.NoError(t, r.resetGrafanaCredentials(adminCredentialsCR(nil)))
	assert.Empty(t, kubeClient.Actions())
}

func TestWaitForReadyGrafanaPod(t *testing.T) {
	terminating := readyPod("grafana-terminating", true)
	terminating.DeletionTimestamp = &metav1.Time{Time: time.Now()}
	terminating.Finalizers = []string{"example.com/keep"}
	kubeClient := kubernetesfake.NewSimpleClientset(
		terminating,
		readyPod("grafana-starting", false),
		readyPod("grafana-ready", true),
	)
	r := newAdminCredentialsReconciler(t, kubeClient)

	name, err := r.waitForReadyGrafanaPod(testNamespace, utils.GrafanaComponentName)
	require.NoError(t, err)
	assert.Equal(t, "grafana-ready", name)
}

func TestWaitForReadyGrafanaPodTimesOut(t *testing.T) {
	shortenGrafanaPodWait(t)
	r := newAdminCredentialsReconciler(t, kubernetesfake.NewSimpleClientset(readyPod("grafana-starting", false)))

	_, err := r.waitForReadyGrafanaPod(testNamespace, utils.GrafanaComponentName)
	assert.ErrorContains(t, err, "no ready Grafana pod with label app.kubernetes.io/name=grafana")
}

func TestWaitForReadyGrafanaPodReturnsListError(t *testing.T) {
	kubeClient := kubernetesfake.NewSimpleClientset()
	kubeClient.PrependReactor("list", "pods", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, apierrors.NewForbidden(schema.GroupResource{Resource: "pods"}, "", assert.AnError)
	})
	r := newAdminCredentialsReconciler(t, kubeClient)

	_, err := r.waitForReadyGrafanaPod(testNamespace, utils.GrafanaComponentName)
	assert.ErrorContains(t, err, "cannot list Grafana pods")
}

func TestIsPodReady(t *testing.T) {
	assert.True(t, isPodReady(readyPod("ready", true)))
	assert.False(t, isPodReady(readyPod("not-ready", false)))

	pending := readyPod("pending", true)
	pending.Status.Phase = corev1.PodPending
	assert.False(t, isPodReady(pending))

	noConditions := readyPod("no-conditions", true)
	noConditions.Status.Conditions = nil
	assert.False(t, isPodReady(noConditions))
}

type recordedPodCommand struct {
	namespace, pod, container string
	command                   []string
	stdin                     string
}

func recordPodCommands(r *GrafanaReconciler, err error) *[]recordedPodCommand {
	calls := &[]recordedPodCommand{}
	r.runPodCommand = func(namespace, pod, container string, command []string, stdin string) (string, string, error) {
		*calls = append(*calls, recordedPodCommand{namespace, pod, container, command, stdin})
		return "out", "err", err
	}
	return calls
}

func TestSyncGrafanaAdminPasswordResetsChangedPassword(t *testing.T) {
	cr := adminCredentialsCR(nil)
	secret := adminSecret(testNamespace, "grafana", "new-password")
	r := newAdminCredentialsReconciler(t, kubernetesfake.NewSimpleClientset(readyPod("grafana-ready", true)),
		secret, existingGrafana(t, cr, "old"))
	calls := recordPodCommands(r, nil)

	require.NoError(t, r.syncGrafanaAdminPassword(cr))

	require.Len(t, *calls, 1)
	call := (*calls)[0]
	assert.Equal(t, testNamespace, call.namespace)
	assert.Equal(t, "grafana-ready", call.pod)
	assert.Equal(t, grafanaContainerName, call.container)
	assert.Equal(t, "new-password\n", call.stdin)
	assert.Contains(t, call.command, "--password-from-stdin")
	assert.NotContains(t, call.command, "new-password", "the password must not be passed as an argument")

	current := &grafv1.Grafana{}
	require.NoError(t, r.Client.Get(context.Background(), types.NamespacedName{Name: "grafana", Namespace: testNamespace}, current))
	assert.Equal(t, adminPasswordChecksum(secret), current.GetAnnotations()[adminPasswordSyncedAnnotation])

	require.NoError(t, r.syncGrafanaAdminPassword(cr))
	assert.Len(t, *calls, 1, "a synced password must not be reset again")
}

func TestSyncGrafanaAdminPasswordReturnsCommandError(t *testing.T) {
	cr := adminCredentialsCR(nil)
	r := newAdminCredentialsReconciler(t, kubernetesfake.NewSimpleClientset(readyPod("grafana-ready", true)),
		adminSecret(testNamespace, "grafana", "new-password"), existingGrafana(t, cr, "old"))
	recordPodCommands(r, assert.AnError)

	err := r.syncGrafanaAdminPassword(cr)
	assert.ErrorIs(t, err, assert.AnError)
	assert.ErrorContains(t, err, "grafana cli reset-admin-password failed in pod grafana-ready")
	assert.NotContains(t, err.Error(), "new-password")

	current := &grafv1.Grafana{}
	require.NoError(t, r.Client.Get(context.Background(), types.NamespacedName{Name: "grafana", Namespace: testNamespace}, current))
	assert.Equal(t, "old", current.GetAnnotations()[adminPasswordSyncedAnnotation])
}
