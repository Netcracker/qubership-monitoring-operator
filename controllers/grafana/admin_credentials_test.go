package grafana

import (
	"context"
	"io"
	"maps"
	"net/url"
	"testing"
	"time"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/go-logr/logr"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/kubernetes"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	k8stesting "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/remotecommand"
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

func existingGrafana(t *testing.T, cr *monv1.PlatformMonitoring, syncedChecksum string) *grafv1.Grafana {
	t.Helper()
	m, err := grafanaWithDefaultSources(cr)
	require.NoError(t, err)
	if syncedChecksum != "" {
		m.SetAnnotations(map[string]string{adminPasswordSyncedAnnotation: syncedChecksum})
	}
	return m
}

// grafanaDeployment mirrors the Deployment that Grafana Operator creates: its selector and the
// pod label it enforces are app=<Grafana resource name>.
func grafanaDeployment(name string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: name + "-deployment", Namespace: testNamespace},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": name}},
		},
	}
}

func grafanaPod(name string, labels map[string]string, ready bool) *corev1.Pod {
	status := corev1.ConditionFalse
	if ready {
		status = corev1.ConditionTrue
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: testNamespace, Labels: labels},
		Status: corev1.PodStatus{
			Phase:      corev1.PodRunning,
			Conditions: []corev1.PodCondition{{Type: corev1.PodReady, Status: status}},
		},
	}
}

func readyPod(name string, ready bool) *corev1.Pod {
	return grafanaPod(name, map[string]string{"app": utils.GrafanaComponentName}, ready)
}

// fakeExecutor records the exec request and the stdin it receives.
type fakeExecutor struct {
	url    *url.URL
	stdin  string
	stdout string
	err    error
}

func (e *fakeExecutor) Stream(options remotecommand.StreamOptions) error {
	return e.StreamWithContext(context.Background(), options)
}

func (e *fakeExecutor) StreamWithContext(_ context.Context, options remotecommand.StreamOptions) error {
	stdin, err := io.ReadAll(options.Stdin)
	if err != nil {
		return err
	}
	e.stdin = string(stdin)
	if _, err = io.WriteString(options.Stdout, e.stdout); err != nil {
		return err
	}
	return e.err
}

// execClientset wraps the fake clientset with a real core REST client, so exec requests can be
// built. The fake clientset returns a nil REST client.
type execClientset struct {
	*kubernetesfake.Clientset
	core corev1client.CoreV1Interface
}

func (c execClientset) CoreV1() corev1client.CoreV1Interface { return c.core }

type execCoreV1 struct {
	corev1client.CoreV1Interface
	restClient rest.Interface
}

func (c execCoreV1) RESTClient() rest.Interface { return c.restClient }

func newExecClientset(t *testing.T, objects ...runtime.Object) kubernetes.Interface {
	t.Helper()
	restClientset, err := kubernetes.NewForConfig(&rest.Config{Host: "http://127.0.0.1:1"})
	require.NoError(t, err)
	fakeClient := kubernetesfake.NewSimpleClientset(objects...)
	return execClientset{
		Clientset: fakeClient,
		core:      execCoreV1{CoreV1Interface: fakeClient.CoreV1(), restClient: restClientset.CoreV1().RESTClient()},
	}
}

func useFakeExecutor(r *GrafanaReconciler, executor *fakeExecutor) {
	r.newExecutor = func(_ *rest.Config, method string, url *url.URL) (remotecommand.Executor, error) {
		executor.url = url
		return executor, nil
	}
}

func getGrafana(t *testing.T, r *GrafanaReconciler, name string) *grafv1.Grafana {
	t.Helper()
	current := &grafv1.Grafana{}
	require.NoError(t, r.Client.Get(context.Background(), types.NamespacedName{Name: name, Namespace: testNamespace}, current))
	return current
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

func TestReadAdminPasswordChecksumUsesCustomName(t *testing.T) {
	cr := adminCredentialsCR(&monv1.Grafana{Name: "custom"})
	secret := adminSecret(testNamespace, "custom", "secret")
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

func TestAddAdminPasswordChecksum(t *testing.T) {
	cr := adminCredentialsCR(nil)
	secret := adminSecret(testNamespace, "grafana", "secret")
	r := newAdminCredentialsReconciler(t, nil, secret)
	m, err := grafanaWithDefaultSources(cr)
	require.NoError(t, err)

	require.NoError(t, r.addAdminPasswordChecksum(cr, m, map[string]string{adminSecretChecksumAnnotation: "old"}))

	assert.Equal(t, adminPasswordChecksum(secret), m.Spec.Deployment.Spec.Template.Annotations[adminSecretChecksumAnnotation])
	assert.NotContains(t, m.GetAnnotations(), adminPasswordSyncedAnnotation)
}

func TestAddAdminPasswordChecksumKeepsExistingChecksumWithoutPassword(t *testing.T) {
	cr := adminCredentialsCR(nil)
	r := newAdminCredentialsReconciler(t, nil)
	m, err := grafanaWithDefaultSources(cr)
	require.NoError(t, err)

	require.NoError(t, r.addAdminPasswordChecksum(cr, m, map[string]string{adminSecretChecksumAnnotation: "old"}))
	assert.Equal(t, "old", m.Spec.Deployment.Spec.Template.Annotations[adminSecretChecksumAnnotation])

	fresh, err := grafanaWithDefaultSources(cr)
	require.NoError(t, err)
	require.NoError(t, r.addAdminPasswordChecksum(cr, fresh, nil))
	assert.NotContains(t, fresh.Spec.Deployment.Spec.Template.Annotations, adminSecretChecksumAnnotation)
}

func TestAddAdminPasswordChecksumReturnsReadError(t *testing.T) {
	cr := adminCredentialsCR(nil)
	r := newAdminCredentialsReconciler(t, nil)
	r.Client = fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build()
	m, err := grafanaWithDefaultSources(cr)
	require.NoError(t, err)

	assert.Error(t, r.addAdminPasswordChecksum(cr, m, nil))
}

func TestHandleGrafanaCreatesGrafanaWithUnsyncedPassword(t *testing.T) {
	cr := adminCredentialsCR(nil)
	secret := adminSecret(testNamespace, "grafana", "secret")
	r := newAdminCredentialsReconciler(t, nil, secret)

	require.NoError(t, r.handleGrafana(cr))

	created := getGrafana(t, r, "grafana")
	assert.Equal(t, adminPasswordChecksum(secret), created.Spec.Deployment.Spec.Template.Annotations[adminSecretChecksumAnnotation])
	assert.NotContains(t, created.GetAnnotations(), adminPasswordSyncedAnnotation,
		"a new Grafana resource can adopt an existing database, so it must get one reset")
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

func TestSyncGrafanaAdminPasswordResetsNewGrafana(t *testing.T) {
	cr := adminCredentialsCR(nil)
	secret := adminSecret(testNamespace, "grafana", "new-password")
	r := newAdminCredentialsReconciler(t,
		newExecClientset(t, grafanaDeployment("grafana"), readyPod("grafana-ready", true)),
		secret, existingGrafana(t, cr, ""))
	executor := &fakeExecutor{stdout: "Admin password changed successfully"}
	useFakeExecutor(r, executor)

	require.NoError(t, r.syncGrafanaAdminPassword(cr))

	require.NotNil(t, executor.url)
	assert.Equal(t, "/api/v1/namespaces/monitoring/pods/grafana-ready/exec", executor.url.Path)
	query := executor.url.Query()
	assert.Equal(t, grafanaContainerName, query.Get("container"))
	assert.Contains(t, query["command"], "--password-from-stdin")
	assert.NotContains(t, executor.url.String(), "new-password", "the password must not be in the exec URL")
	assert.Equal(t, "new-password\n", executor.stdin)
	assert.Equal(t, adminPasswordChecksum(secret), getGrafana(t, r, "grafana").GetAnnotations()[adminPasswordSyncedAnnotation])

	executor.url = nil
	require.NoError(t, r.syncGrafanaAdminPassword(cr))
	assert.Nil(t, executor.url, "a synced password must not be reset again")
}

func TestSyncGrafanaAdminPasswordPendingWithoutDeployment(t *testing.T) {
	cr := adminCredentialsCR(nil)
	r := newAdminCredentialsReconciler(t, nil, adminSecret(testNamespace, "grafana", "new"), existingGrafana(t, cr, "old"))

	err := r.syncGrafanaAdminPassword(cr)
	assert.ErrorIs(t, err, ErrAdminPasswordSyncPending)
	assert.True(t, IsPending(err))
	assert.Equal(t, "old", getGrafana(t, r, "grafana").GetAnnotations()[adminPasswordSyncedAnnotation])
}

func TestSyncGrafanaAdminPasswordKeepsChecksumWhenResetFails(t *testing.T) {
	cr := adminCredentialsCR(nil)
	r := newAdminCredentialsReconciler(t,
		newExecClientset(t, grafanaDeployment("grafana"), readyPod("grafana-ready", true)),
		adminSecret(testNamespace, "grafana", "new-password"), existingGrafana(t, cr, "old"))
	useFakeExecutor(r, &fakeExecutor{stdout: "output", err: assert.AnError})

	err := r.syncGrafanaAdminPassword(cr)
	assert.ErrorIs(t, err, assert.AnError)
	assert.False(t, IsPending(err))
	assert.ErrorContains(t, err, "grafana cli reset-admin-password failed in pod grafana-ready")
	assert.ErrorContains(t, err, "stdout: output")
	assert.NotContains(t, err.Error(), "new-password")
	assert.Equal(t, "old", getGrafana(t, r, "grafana").GetAnnotations()[adminPasswordSyncedAnnotation],
		"a failed reset must be retried on the next reconcile")
}

func TestSyncGrafanaAdminPasswordReturnsGrafanaReadError(t *testing.T) {
	r := newAdminCredentialsReconciler(t, nil)
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

func TestResetGrafanaCredentialsReturnsSecretReadError(t *testing.T) {
	r := newAdminCredentialsReconciler(t, nil)
	r.Client = fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build()

	assert.Error(t, r.resetGrafanaCredentials(adminCredentialsCR(nil)))
}

func TestExecInPodReturnsExecutorCreationError(t *testing.T) {
	r := newAdminCredentialsReconciler(t, newExecClientset(t))
	r.newExecutor = func(*rest.Config, string, *url.URL) (remotecommand.Executor, error) {
		return nil, assert.AnError
	}

	_, _, err := r.execInPod(testNamespace, "grafana-ready", grafanaContainerName, []string{"true"}, "")
	assert.ErrorIs(t, err, assert.AnError)
	assert.ErrorContains(t, err, "cannot create exec client for pod grafana-ready")
}

func TestExecInPodUsesSPDYExecutorByDefault(t *testing.T) {
	r := newUnreachableExecReconciler(t)

	// The default executor connects to an unreachable API server, so the stream fails.
	_, _, err := r.execInPod(testNamespace, "grafana-ready", grafanaContainerName, []string{"true"}, "")
	assert.Error(t, err)
}

func newUnreachableExecReconciler(t *testing.T) *GrafanaReconciler {
	t.Helper()
	r := newAdminCredentialsReconciler(t, newExecClientset(t))
	r.config = &rest.Config{Host: "http://127.0.0.1:1", Timeout: time.Second}
	return r
}

func TestFindReadyGrafanaPodUsesDeploymentSelector(t *testing.T) {
	// spec.grafana.labels can override app.kubernetes.io/name on the pod template; the lookup must
	// still find the pod through the Deployment selector.
	cr := adminCredentialsCR(&monv1.Grafana{Labels: map[string]string{"app.kubernetes.io/name": "overridden"}})
	m, err := grafanaWithDefaultSources(cr)
	require.NoError(t, err)
	podLabels := maps.Clone(m.Spec.Deployment.Spec.Template.Labels)
	require.Equal(t, "overridden", podLabels["app.kubernetes.io/name"])
	podLabels["app"] = m.GetName()

	terminating := grafanaPod("grafana-terminating", podLabels, true)
	terminating.DeletionTimestamp = &metav1.Time{Time: time.Now()}
	terminating.Finalizers = []string{"example.com/keep"}
	kubeClient := kubernetesfake.NewSimpleClientset(
		grafanaDeployment(m.GetName()),
		terminating,
		grafanaPod("grafana-starting", podLabels, false),
		grafanaPod("grafana-ready", podLabels, true),
		grafanaPod("other-ready", map[string]string{"app": "other"}, true),
	)
	r := newAdminCredentialsReconciler(t, kubeClient)

	name, err := r.findReadyGrafanaPod(testNamespace, m.GetName())
	require.NoError(t, err)
	assert.Equal(t, "grafana-ready", name)
}

func TestFindReadyGrafanaPodUsesCustomName(t *testing.T) {
	kubeClient := kubernetesfake.NewSimpleClientset(
		grafanaDeployment("custom"),
		readyPod("grafana-ready", true),
		grafanaPod("custom-ready", map[string]string{"app": "custom"}, true),
	)
	r := newAdminCredentialsReconciler(t, kubeClient)

	name, err := r.findReadyGrafanaPod(testNamespace, "custom")
	require.NoError(t, err)
	assert.Equal(t, "custom-ready", name)
}

func TestFindReadyGrafanaPodPendingWithoutReadyPod(t *testing.T) {
	r := newAdminCredentialsReconciler(t, kubernetesfake.NewSimpleClientset(
		grafanaDeployment("grafana"), readyPod("grafana-starting", false)))

	_, err := r.findReadyGrafanaPod(testNamespace, utils.GrafanaComponentName)
	assert.ErrorIs(t, err, ErrAdminPasswordSyncPending)
	assert.ErrorContains(t, err, "no ready pod of Deployment monitoring/grafana-deployment")
}

func TestFindReadyGrafanaPodRejectsEmptySelector(t *testing.T) {
	deployment := grafanaDeployment("grafana")
	deployment.Spec.Selector = &metav1.LabelSelector{}
	r := newAdminCredentialsReconciler(t, kubernetesfake.NewSimpleClientset(deployment))

	_, err := r.findReadyGrafanaPod(testNamespace, utils.GrafanaComponentName)
	assert.ErrorContains(t, err, "has an empty selector")
	assert.False(t, IsPending(err))
}

func TestFindReadyGrafanaPodRejectsInvalidSelector(t *testing.T) {
	deployment := grafanaDeployment("grafana")
	deployment.Spec.Selector = &metav1.LabelSelector{MatchExpressions: []metav1.LabelSelectorRequirement{
		{Key: "app", Operator: "Bogus"},
	}}
	r := newAdminCredentialsReconciler(t, kubernetesfake.NewSimpleClientset(deployment))

	_, err := r.findReadyGrafanaPod(testNamespace, utils.GrafanaComponentName)
	assert.ErrorContains(t, err, "invalid selector in Grafana Deployment")
}

func TestFindReadyGrafanaPodReturnsAPIErrors(t *testing.T) {
	forbidden := func(resource string) k8stesting.ReactionFunc {
		return func(k8stesting.Action) (bool, runtime.Object, error) {
			return true, nil, apierrors.NewForbidden(schema.GroupResource{Resource: resource}, "", assert.AnError)
		}
	}

	deploymentClient := kubernetesfake.NewSimpleClientset()
	deploymentClient.PrependReactor("get", "deployments", forbidden("deployments"))
	_, err := newAdminCredentialsReconciler(t, deploymentClient).findReadyGrafanaPod(testNamespace, "grafana")
	assert.ErrorContains(t, err, "cannot get Grafana Deployment monitoring/grafana-deployment")
	assert.False(t, IsPending(err))

	podClient := kubernetesfake.NewSimpleClientset(grafanaDeployment("grafana"))
	podClient.PrependReactor("list", "pods", forbidden("pods"))
	_, err = newAdminCredentialsReconciler(t, podClient).findReadyGrafanaPod(testNamespace, "grafana")
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

func TestIsPending(t *testing.T) {
	assert.True(t, IsPending(ErrAdminPasswordSyncPending))
	assert.True(t, IsPending(ErrDatasourceMigrationPending))
	assert.False(t, IsPending(assert.AnError))
	assert.False(t, IsPending(nil))
}
