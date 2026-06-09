package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/kubernetes/fake"
	ktesting "k8s.io/client-go/testing"
)

const (
	validPrivateKey = "-----BEGIN PRIVATE KEY-----\nkey\n-----END PRIVATE KEY-----"
	validCACert     = "-----BEGIN CERTIFICATE-----\nca\n-----END CERTIFICATE-----"
	validPeerCert   = "-----BEGIN CERTIFICATE-----\ncert\n-----END CERTIFICATE-----"
)

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    slog.Level
		wantErr bool
	}{
		{name: "debug", value: "debug", want: slog.LevelDebug},
		{name: "info", value: "info", want: slog.LevelInfo},
		{name: "warn", value: "warn", want: slog.LevelWarn},
		{name: "warning alias", value: "warning", want: slog.LevelWarn},
		{name: "error", value: "error", want: slog.LevelError},
		{name: "case insensitive", value: "DEBUG", want: slog.LevelDebug},
		{name: "invalid", value: "trace", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLogLevel(tt.value)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLoggerDefaultsToInfoWhenConfigured(t *testing.T) {
	logger := newLogger(slog.LevelInfo)

	if logger.Enabled(context.TODO(), slog.LevelDebug) {
		t.Fatal("debug log should not be enabled at info level")
	}
	if !logger.Enabled(context.TODO(), slog.LevelInfo) {
		t.Fatal("info log should be enabled at info level")
	}
}

func TestParseOptions(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		env        map[string]string
		wantNS     string
		wantSecret string
		wantLevel  slog.Level
		wantErr    bool
	}{
		{
			name:       "defaults fall back to monitoring",
			args:       []string{"prog"},
			wantNS:     "monitoring",
			wantSecret: "kube-etcd-client-certs",
			wantLevel:  slog.LevelInfo,
		},
		{
			name:       "namespace env wins over default",
			args:       []string{"prog"},
			env:        map[string]string{"NAMESPACE": "from-env"},
			wantNS:     "from-env",
			wantSecret: "kube-etcd-client-certs",
			wantLevel:  slog.LevelInfo,
		},
		{
			name:       "flags override env",
			args:       []string{"prog", "--namespace=from-flag", "--secret=custom", "--log-level=debug"},
			env:        map[string]string{"NAMESPACE": "from-env"},
			wantNS:     "from-flag",
			wantSecret: "custom",
			wantLevel:  slog.LevelDebug,
		},
		{
			name:    "invalid log level surfaces error",
			args:    []string{"prog", "--log-level=trace"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withArgsAndEnv(t, tt.args, tt.env)
			opts, err := parseOptions()
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if opts.namespace != tt.wantNS {
				t.Fatalf("got namespace %q, want %q", opts.namespace, tt.wantNS)
			}
			if opts.secretName != tt.wantSecret {
				t.Fatalf("got secret %q, want %q", opts.secretName, tt.wantSecret)
			}
			if opts.logLevel != tt.wantLevel {
				t.Fatalf("got level %v, want %v", opts.logLevel, tt.wantLevel)
			}
		})
	}
}

func withArgsAndEnv(t *testing.T, args []string, env map[string]string) {
	t.Helper()
	origArgs := os.Args
	origCmdLine := flag.CommandLine

	flag.CommandLine = flag.NewFlagSet(args[0], flag.ContinueOnError)
	flag.CommandLine.SetOutput(bytes.NewBuffer(nil))
	os.Args = args

	for k, v := range env {
		t.Setenv(k, v)
	}
	if _, set := env["NAMESPACE"]; !set {
		t.Setenv("NAMESPACE", "")
	}

	t.Cleanup(func() {
		os.Args = origArgs
		flag.CommandLine = origCmdLine
	})
}

func TestCertVerify(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		ca      string
		crt     string
		wantErr bool
	}{
		{name: "valid certificates", key: validPrivateKey, ca: validCACert, crt: validPeerCert},
		{name: "valid rsa private key", key: "-----BEGIN RSA PRIVATE KEY-----\nkey\n-----END RSA PRIVATE KEY-----", ca: validCACert, crt: validPeerCert},
		{name: "empty key", key: "", ca: validCACert, crt: validPeerCert, wantErr: true},
		{name: "invalid key", key: "key", ca: validCACert, crt: validPeerCert, wantErr: true},
		{name: "invalid ca", key: validPrivateKey, ca: "ca", crt: validPeerCert, wantErr: true},
		{name: "invalid peer cert", key: validPrivateKey, ca: validCACert, crt: "cert", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := certVerify(tt.key, tt.ca, tt.crt)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestBuildEtcdSecret(t *testing.T) {
	secret := buildEtcdSecret("monitoring", "kube-etcd-client-certs", certificateData{
		key: validPrivateKey,
		ca:  validCACert,
		crt: validPeerCert,
	})

	if secret.Name != "kube-etcd-client-certs" {
		t.Fatalf("got secret name %q", secret.Name)
	}
	if secret.Namespace != "monitoring" {
		t.Fatalf("got namespace %q", secret.Namespace)
	}
	if string(secret.Data["etcd-client.key"]) != validPrivateKey {
		t.Fatal("secret key data was not set")
	}
	if string(secret.Data["etcd-client-ca.crt"]) != validCACert {
		t.Fatal("secret CA data was not set")
	}
	if string(secret.Data["etcd-client.crt"]) != validPeerCert {
		t.Fatal("secret cert data was not set")
	}
}

func TestBuildEtcdService(t *testing.T) {
	svc := buildEtcdService()

	if svc.Namespace != "kube-system" {
		t.Fatalf("got namespace %q, want kube-system", svc.Namespace)
	}
	if svc.Name != "etcd" {
		t.Fatalf("got name %q, want etcd", svc.Name)
	}
	if svc.Spec.Selector["component"] != "etcd" {
		t.Fatalf("got selector %#v", svc.Spec.Selector)
	}
	if svc.Spec.ClusterIP != corev1.ClusterIPNone {
		t.Fatalf("got ClusterIP %q, want None", svc.Spec.ClusterIP)
	}
	if len(svc.Spec.Ports) != 1 || svc.Spec.Ports[0].Port != 2379 {
		t.Fatalf("got ports %#v", svc.Spec.Ports)
	}
}

func TestCreateOrUpdateSecret(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	log := testLogger()
	ctx := context.TODO()

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd-certs", Namespace: "monitoring"},
		Data:       map[string][]byte{"value": []byte("old")},
	}

	if err := createOrUpdateSecret(ctx, clientset, secret, log); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	got, err := clientset.CoreV1().Secrets("monitoring").Get(ctx, "etcd-certs", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("get created secret failed: %v", err)
	}
	if string(got.Data["value"]) != "old" {
		t.Fatalf("got data %q", got.Data["value"])
	}

	secret.Data = map[string][]byte{"value": []byte("new")}
	if err := createOrUpdateSecret(ctx, clientset, secret, log); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	got, err = clientset.CoreV1().Secrets("monitoring").Get(ctx, "etcd-certs", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("get updated secret failed: %v", err)
	}
	if string(got.Data["value"]) != "new" {
		t.Fatalf("got data %q, want new", got.Data["value"])
	}
}

func TestCreateOrUpdateSecretGetError(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	clientset.PrependReactor("get", "secrets", func(action ktesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("boom")
	})

	err := createOrUpdateSecret(context.TODO(), clientset, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd-certs", Namespace: "monitoring"},
	}, testLogger())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateOrUpdateSecretCreateError(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	clientset.PrependReactor("create", "secrets", func(action ktesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("create boom")
	})

	err := createOrUpdateSecret(context.TODO(), clientset, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd-certs", Namespace: "monitoring"},
	}, testLogger())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateOrUpdateSecretUpdateError(t *testing.T) {
	existing := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd-certs", Namespace: "monitoring"},
	}
	clientset := fake.NewSimpleClientset(existing)
	clientset.PrependReactor("update", "secrets", func(action ktesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("update boom")
	})

	err := createOrUpdateSecret(context.TODO(), clientset, existing, testLogger())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateOrUpdateEtcdService(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	log := testLogger()
	ctx := context.TODO()

	if err := createOrUpdateEtcdService(ctx, clientset, log); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	got, err := clientset.CoreV1().Services("kube-system").Get(ctx, "etcd", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("get created service failed: %v", err)
	}
	if got.Spec.Selector["component"] != "etcd" {
		t.Fatalf("got selector %#v", got.Spec.Selector)
	}

	got.Labels["stale"] = "true"
	if _, err := clientset.CoreV1().Services("kube-system").Update(ctx, got, metav1.UpdateOptions{}); err != nil {
		t.Fatalf("seed update failed: %v", err)
	}

	if err := createOrUpdateEtcdService(ctx, clientset, log); err != nil {
		t.Fatalf("update failed: %v", err)
	}
	got, err = clientset.CoreV1().Services("kube-system").Get(ctx, "etcd", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("get updated service failed: %v", err)
	}
	if got.Labels["stale"] != "true" {
		t.Fatal("existing labels should be preserved on update")
	}
}

func TestCreateOrUpdateEtcdServiceGetError(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	clientset.PrependReactor("get", "services", func(action ktesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("get boom")
	})

	err := createOrUpdateEtcdService(context.TODO(), clientset, testLogger())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateOrUpdateEtcdServiceCreateError(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	clientset.PrependReactor("create", "services", func(action ktesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("create boom")
	})

	err := createOrUpdateEtcdService(context.TODO(), clientset, testLogger())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateOrUpdateEtcdServiceUpdateError(t *testing.T) {
	clientset := fake.NewSimpleClientset(&corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd", Namespace: "kube-system"},
	})
	clientset.PrependReactor("update", "services", func(action ktesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("update boom")
	})

	err := createOrUpdateEtcdService(context.TODO(), clientset, testLogger())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRunningEtcdPod(t *testing.T) {
	runningPod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd-running"},
		Status:     corev1.PodStatus{Phase: corev1.PodRunning},
	}
	pendingPod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd-pending"},
		Status:     corev1.PodStatus{Phase: corev1.PodPending},
	}

	got, err := runningEtcdPod([]corev1.Pod{pendingPod, runningPod})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != runningPod.Name {
		t.Fatalf("got pod %q, want %q", got.Name, runningPod.Name)
	}

	if _, err := runningEtcdPod([]corev1.Pod{pendingPod}); err == nil {
		t.Fatal("expected error when there are no running pods")
	}
}

func TestCertPathsFromArgs(t *testing.T) {
	peerKey, caCrt, peerCrt := certPathsFromArgs([]string{
		"--peer-key-file=/custom/peer.key",
		"--peer-trusted-ca-file=/custom/ca.crt",
		"--peer-cert-file=/custom/peer.crt",
		"--ignored",
	}, "default-key", "default-ca", "default-crt")

	if peerKey != "/custom/peer.key" {
		t.Fatalf("got peer key %q", peerKey)
	}
	if caCrt != "/custom/ca.crt" {
		t.Fatalf("got ca cert %q", caCrt)
	}
	if peerCrt != "/custom/peer.crt" {
		t.Fatalf("got peer cert %q", peerCrt)
	}
}

func TestCertPathsFromPodUsesDefaultsWhenEtcdContainerIsMissing(t *testing.T) {
	peerKey, caCrt, peerCrt := certPathsFromPod(corev1.Pod{
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{Name: "sidecar"}},
		},
	}, "default-key", "default-ca", "default-crt")

	if peerKey != "default-key" || caCrt != "default-ca" || peerCrt != "default-crt" {
		t.Fatalf("got paths %q %q %q", peerKey, caCrt, peerCrt)
	}
}

func TestCertPathsFromPodReadsEtcdContainer(t *testing.T) {
	pod := corev1.Pod{Spec: corev1.PodSpec{Containers: []corev1.Container{
		{Name: "sidecar"},
		{Name: "etcd", Command: []string{
			"etcd",
			"--peer-key-file=/x/peer.key",
			"--peer-trusted-ca-file=/x/ca.crt",
			"--peer-cert-file=/x/peer.crt",
		}},
	}}}

	peerKey, caCrt, peerCrt := certPathsFromPod(pod, "default-key", "default-ca", "default-crt")
	if peerKey != "/x/peer.key" || caCrt != "/x/ca.crt" || peerCrt != "/x/peer.crt" {
		t.Fatalf("got paths %q %q %q", peerKey, caCrt, peerCrt)
	}
}

func TestPodListError(t *testing.T) {
	wantErr := errors.New("boom")
	if got := podListError(wantErr); got != wantErr {
		t.Fatalf("got %v, want passthrough %v", got, wantErr)
	}
	if got := podListError(nil); got == nil {
		t.Fatal("expected synthesized error when input is nil")
	}
}

func TestExtractEtcdCertPaths(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd-0", Namespace: "kube-system", Labels: map[string]string{"component": "etcd"}},
		Spec: corev1.PodSpec{Containers: []corev1.Container{{
			Name: "etcd",
			Command: []string{
				"etcd",
				"--peer-key-file=/custom/peer.key",
				"--peer-trusted-ca-file=/custom/ca.crt",
				"--peer-cert-file=/custom/peer.crt",
			},
		}}},
		Status: corev1.PodStatus{Phase: corev1.PodRunning},
	}
	clientset := fake.NewSimpleClientset(pod)

	peerKey, caCrt, peerCrt, err := extractEtcdCertPaths(context.TODO(), clientset, testLogger(), "kube-system", "component=etcd")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if peerKey != "/custom/peer.key" || caCrt != "/custom/ca.crt" || peerCrt != "/custom/peer.crt" {
		t.Fatalf("got paths %q %q %q", peerKey, caCrt, peerCrt)
	}
}

func TestExtractEtcdCertPathsNoPods(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	_, _, _, err := extractEtcdCertPaths(context.TODO(), clientset, testLogger(), "kube-system", "component=etcd")
	if err == nil {
		t.Fatal("expected error when no etcd pods exist")
	}
}

func TestExtractEtcdCertPathsListError(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	clientset.PrependReactor("list", "pods", func(action ktesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("list boom")
	})
	_, _, _, err := extractEtcdCertPaths(context.TODO(), clientset, testLogger(), "kube-system", "component=etcd")
	if err == nil {
		t.Fatal("expected error when list fails")
	}
}

func TestExtractEtcdCertPathsNoRunningPod(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd-0", Namespace: "kube-system", Labels: map[string]string{"component": "etcd"}},
		Status:     corev1.PodStatus{Phase: corev1.PodPending},
	}
	clientset := fake.NewSimpleClientset(pod)

	_, _, _, err := extractEtcdCertPaths(context.TODO(), clientset, testLogger(), "kube-system", "component=etcd")
	if err == nil {
		t.Fatal("expected error when no running etcd pod exists")
	}
}

func TestReadFileToString(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.txt")
	if err := os.WriteFile(path, []byte("hello"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	got, err := readFileToString(testLogger(), path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hello" {
		t.Fatalf("got %q, want %q", got, "hello")
	}

	if _, err := readFileToString(testLogger(), filepath.Join(dir, "missing")); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestGetCertsFromConfigmapAndSecret(t *testing.T) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd-metric-serving-ca", Namespace: "openshift-etcd-operator"},
		Data:       map[string]string{"ca-bundle.crt": validCACert},
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd-metric-client", Namespace: "openshift-etcd-operator"},
		Data: map[string][]byte{
			"tls.key": []byte(validPrivateKey),
			"tls.crt": []byte(validPeerCert),
		},
	}
	clientset := fake.NewSimpleClientset(cm, secret)

	got, err := getCertsFromConfigmapAndSecret(context.TODO(), clientset, testLogger(),
		"openshift-etcd-operator", "etcd-metric-serving-ca", "etcd-metric-client")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.key != validPrivateKey || got.ca != validCACert || got.crt != validPeerCert {
		t.Fatalf("got cert data %#v", got)
	}
}

func TestGetCertsFromConfigmapAndSecretMissingConfigMap(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	_, err := getCertsFromConfigmapAndSecret(context.TODO(), clientset, testLogger(),
		"openshift-etcd-operator", "etcd-metric-serving-ca", "etcd-metric-client")
	if err == nil {
		t.Fatal("expected error when configmap is missing")
	}
}

func TestGetCertsFromConfigmapAndSecretMissingSecret(t *testing.T) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd-metric-serving-ca", Namespace: "openshift-etcd-operator"},
		Data:       map[string]string{"ca-bundle.crt": validCACert},
	}
	clientset := fake.NewSimpleClientset(cm)
	_, err := getCertsFromConfigmapAndSecret(context.TODO(), clientset, testLogger(),
		"openshift-etcd-operator", "etcd-metric-serving-ca", "etcd-metric-client")
	if err == nil {
		t.Fatal("expected error when secret is missing")
	}
}

func TestGetCertsFromHostpath(t *testing.T) {
	dir := t.TempDir()
	peerKey := filepath.Join(dir, "peer.key")
	caCrt := filepath.Join(dir, "ca.crt")
	peerCrt := filepath.Join(dir, "peer.crt")
	if err := os.WriteFile(peerKey, []byte(validPrivateKey), 0o600); err != nil {
		t.Fatalf("write key: %v", err)
	}
	if err := os.WriteFile(caCrt, []byte(validCACert), 0o600); err != nil {
		t.Fatalf("write ca: %v", err)
	}
	if err := os.WriteFile(peerCrt, []byte(validPeerCert), 0o600); err != nil {
		t.Fatalf("write crt: %v", err)
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd-0", Namespace: "kube-system", Labels: map[string]string{"component": "etcd"}},
		Spec: corev1.PodSpec{Containers: []corev1.Container{{
			Name: "etcd",
			Command: []string{
				"etcd",
				"--peer-key-file=" + peerKey,
				"--peer-trusted-ca-file=" + caCrt,
				"--peer-cert-file=" + peerCrt,
			},
		}}},
		Status: corev1.PodStatus{Phase: corev1.PodRunning},
	}
	clientset := fake.NewSimpleClientset(pod)

	got, err := getCertsFromHostpath(context.TODO(), clientset, testLogger(), "kube-system", "component=etcd")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.key != validPrivateKey || got.ca != validCACert || got.crt != validPeerCert {
		t.Fatalf("got cert data %#v", got)
	}
}

func TestGetCertsFromHostpathExtractFails(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	_, err := getCertsFromHostpath(context.TODO(), clientset, testLogger(), "kube-system", "component=etcd")
	if err == nil {
		t.Fatal("expected error when no etcd pods exist")
	}
}

func TestGetCertsFromHostpathMissingFiles(t *testing.T) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd-0", Namespace: "kube-system", Labels: map[string]string{"component": "etcd"}},
		Spec: corev1.PodSpec{Containers: []corev1.Container{{
			Name: "etcd",
			Command: []string{
				"etcd",
				"--peer-key-file=/does/not/exist/peer.key",
				"--peer-trusted-ca-file=/does/not/exist/ca.crt",
				"--peer-cert-file=/does/not/exist/peer.crt",
			},
		}}},
		Status: corev1.PodStatus{Phase: corev1.PodRunning},
	}
	clientset := fake.NewSimpleClientset(pod)

	if _, err := getCertsFromHostpath(context.TODO(), clientset, testLogger(), "kube-system", "component=etcd"); err == nil {
		t.Fatal("expected error when cert files don't exist")
	}
}

func TestLoadEtcdCertsOpenshiftV4(t *testing.T) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd-metric-serving-ca", Namespace: "openshift-etcd-operator"},
		Data:       map[string]string{"ca-bundle.crt": validCACert},
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd-metric-client", Namespace: "openshift-etcd-operator"},
		Data: map[string][]byte{
			"tls.key": []byte(validPrivateKey),
			"tls.crt": []byte(validPeerCert),
		},
	}
	clientset := fake.NewSimpleClientset(cm, secret)

	got, err := loadEtcdCerts(context.TODO(), clientset, true, testLogger())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.crt != validPeerCert {
		t.Fatalf("got %#v", got)
	}
}

func TestLoadEtcdCertsOpenshiftV4Error(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	if _, err := loadEtcdCerts(context.TODO(), clientset, true, testLogger()); err == nil {
		t.Fatal("expected error when configmap/secret missing")
	}
}

func TestLoadEtcdCertsHostpathError(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	if _, err := loadEtcdCerts(context.TODO(), clientset, false, testLogger()); err == nil {
		t.Fatal("expected error when no etcd pod exists")
	}
}

func TestLogEtcdCertsError(t *testing.T) {
	log := testLogger()
	forbidden := apierrors.NewForbidden(schema.GroupResource{Resource: "secrets"}, "etcd", errors.New("nope"))
	logEtcdCertsError(log, forbidden, "msg")
	logEtcdCertsError(log, errors.New("plain"), "msg")
}

func TestHasOpenshiftSecurityAPI(t *testing.T) {
	tests := []struct {
		name      string
		resources []*metav1.APIResourceList
		want      bool
	}{
		{
			name: "present",
			resources: []*metav1.APIResourceList{
				{GroupVersion: "v1"},
				{GroupVersion: "security.openshift.io/v1"},
			},
			want: true,
		},
		{
			name:      "absent",
			resources: []*metav1.APIResourceList{{GroupVersion: "v1"}},
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := &fakediscovery.FakeDiscovery{Fake: &ktesting.Fake{Resources: tt.resources}}
			got, err := hasOpenshiftSecurityAPI(dc)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasOpenshiftSecurityAPIError(t *testing.T) {
	dc := &fakediscovery.FakeDiscovery{Fake: &ktesting.Fake{}}
	dc.Fake.PrependReactor("get", "group", func(action ktesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("api boom")
	})

	if _, err := hasOpenshiftSecurityAPI(dc); err == nil {
		t.Fatal("expected error")
	}
}

func TestRunOpenshiftV4(t *testing.T) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd-metric-serving-ca", Namespace: "openshift-etcd-operator"},
		Data:       map[string]string{"ca-bundle.crt": validCACert},
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd-metric-client", Namespace: "openshift-etcd-operator"},
		Data: map[string][]byte{
			"tls.key": []byte(validPrivateKey),
			"tls.crt": []byte(validPeerCert),
		},
	}
	clientset := fake.NewSimpleClientset(cm, secret)
	clients := &kubeClients{
		clientset: clientset,
		discoveryClient: &fakediscovery.FakeDiscovery{Fake: &ktesting.Fake{Resources: []*metav1.APIResourceList{
			{GroupVersion: "security.openshift.io/v1"},
		}}},
	}

	err := run(context.TODO(), clients, appOptions{namespace: "monitoring", secretName: "kube-etcd-client-certs"}, testLogger())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := clientset.CoreV1().Secrets("monitoring").Get(context.TODO(), "kube-etcd-client-certs", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("get secret: %v", err)
	}
	if string(got.Data["etcd-client.crt"]) != validPeerCert {
		t.Fatalf("got cert %q", got.Data["etcd-client.crt"])
	}

	if _, err := clientset.CoreV1().Services("kube-system").Get(context.TODO(), "etcd", metav1.GetOptions{}); err == nil {
		t.Fatal("etcd Service should not be created on OpenShift v4")
	}
}

func TestRunKubernetes(t *testing.T) {
	dir := t.TempDir()
	peerKey := filepath.Join(dir, "peer.key")
	caCrt := filepath.Join(dir, "ca.crt")
	peerCrt := filepath.Join(dir, "peer.crt")
	for path, data := range map[string]string{peerKey: validPrivateKey, caCrt: validCACert, peerCrt: validPeerCert} {
		if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd-0", Namespace: "kube-system", Labels: map[string]string{"component": "etcd"}},
		Spec: corev1.PodSpec{Containers: []corev1.Container{{
			Name: "etcd",
			Command: []string{
				"etcd",
				"--peer-key-file=" + peerKey,
				"--peer-trusted-ca-file=" + caCrt,
				"--peer-cert-file=" + peerCrt,
			},
		}}},
		Status: corev1.PodStatus{Phase: corev1.PodRunning},
	}
	clientset := fake.NewSimpleClientset(pod)
	clients := &kubeClients{
		clientset:       clientset,
		discoveryClient: &fakediscovery.FakeDiscovery{Fake: &ktesting.Fake{Resources: []*metav1.APIResourceList{{GroupVersion: "v1"}}}},
	}

	err := run(context.TODO(), clients, appOptions{namespace: "monitoring", secretName: "kube-etcd-client-certs"}, testLogger())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := clientset.CoreV1().Secrets("monitoring").Get(context.TODO(), "kube-etcd-client-certs", metav1.GetOptions{}); err != nil {
		t.Fatalf("expected secret to exist: %v", err)
	}
	if _, err := clientset.CoreV1().Services("kube-system").Get(context.TODO(), "etcd", metav1.GetOptions{}); err != nil {
		t.Fatalf("expected etcd Service to exist on Kubernetes: %v", err)
	}
}

func TestRunCertsLoadError(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	clients := &kubeClients{
		clientset:       clientset,
		discoveryClient: &fakediscovery.FakeDiscovery{Fake: &ktesting.Fake{Resources: []*metav1.APIResourceList{{GroupVersion: "v1"}}}},
	}
	err := run(context.TODO(), clients, appOptions{namespace: "monitoring", secretName: "etcd"}, testLogger())
	if err == nil {
		t.Fatal("expected error when cert load fails")
	}
}

func TestRunSecretWriteError(t *testing.T) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd-metric-serving-ca", Namespace: "openshift-etcd-operator"},
		Data:       map[string]string{"ca-bundle.crt": validCACert},
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd-metric-client", Namespace: "openshift-etcd-operator"},
		Data:       map[string][]byte{"tls.key": []byte(validPrivateKey), "tls.crt": []byte(validPeerCert)},
	}
	clientset := fake.NewSimpleClientset(cm, secret)
	clientset.PrependReactor("create", "secrets", func(action ktesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("create boom")
	})
	clients := &kubeClients{
		clientset: clientset,
		discoveryClient: &fakediscovery.FakeDiscovery{Fake: &ktesting.Fake{Resources: []*metav1.APIResourceList{
			{GroupVersion: "security.openshift.io/v1"},
		}}},
	}

	err := run(context.TODO(), clients, appOptions{namespace: "monitoring", secretName: "etcd"}, testLogger())
	if err == nil {
		t.Fatal("expected error when secret write fails")
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}
