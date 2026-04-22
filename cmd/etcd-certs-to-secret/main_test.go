package main

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	monitoringv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
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

func TestCertVerify(t *testing.T) {
	log := testLogger()

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
			err := certVerify(tt.key, tt.ca, tt.crt, log)
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
	cr := platformMonitoring()
	secret, err := buildEtcdSecret(cr, "monitoring", "kube-etcd-client-certs", certificateData{
		key: validPrivateKey,
		ca:  validCACert,
		crt: validPeerCert,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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
	if secret.Labels["team"] != "monitoring" {
		t.Fatal("custom labels were not copied")
	}
	if secret.Annotations["owner"] != "platform" {
		t.Fatal("custom annotations were not copied")
	}
	if len(secret.OwnerReferences) != 1 || secret.OwnerReferences[0].Name != cr.Name {
		t.Fatalf("unexpected owner references: %#v", secret.OwnerReferences)
	}
}

func TestCreateOrUpdateSecret(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	log := testLogger()

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "etcd-certs", Namespace: "monitoring"},
		Data:       map[string][]byte{"value": []byte("old")},
	}

	if err := createOrUpdateSecret(clientset, secret, log); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	got, err := clientset.CoreV1().Secrets("monitoring").Get(context.TODO(), "etcd-certs", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("get created secret failed: %v", err)
	}
	if string(got.Data["value"]) != "old" {
		t.Fatalf("got data %q", got.Data["value"])
	}

	secret.ResourceVersion = got.ResourceVersion
	secret.Data = map[string][]byte{"value": []byte("new")}
	if err := createOrUpdateSecret(clientset, secret, log); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	got, err = clientset.CoreV1().Secrets("monitoring").Get(context.TODO(), "etcd-certs", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("get updated secret failed: %v", err)
	}
	if string(got.Data["value"]) != "new" {
		t.Fatalf("got data %q, want new", got.Data["value"])
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

func TestEtcdService(t *testing.T) {
	tests := []struct {
		name            string
		isOpenshift     bool
		isOpenshiftV4   bool
		wantNamespace   string
		wantSelectorKey string
		wantPorts       int
	}{
		{name: "kubernetes", wantNamespace: "kube-system", wantSelectorKey: "component", wantPorts: 1},
		{name: "openshift v3", isOpenshift: true, wantNamespace: "kube-system", wantSelectorKey: "openshift.io/component", wantPorts: 1},
		{name: "openshift v4", isOpenshift: true, isOpenshiftV4: true, wantNamespace: "openshift-etcd", wantSelectorKey: "etcd", wantPorts: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := etcdService(tt.isOpenshift, tt.wantNamespace, tt.isOpenshiftV4)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if service.Namespace != tt.wantNamespace {
				t.Fatalf("got namespace %q, want %q", service.Namespace, tt.wantNamespace)
			}
			if _, ok := service.Spec.Selector[tt.wantSelectorKey]; !ok {
				t.Fatalf("selector %q not found in %#v", tt.wantSelectorKey, service.Spec.Selector)
			}
			if len(service.Spec.Ports) != tt.wantPorts {
				t.Fatalf("got %d ports, want %d", len(service.Spec.Ports), tt.wantPorts)
			}
			if service.Spec.ClusterIP != "" {
				t.Fatalf("got ClusterIP %q, want empty", service.Spec.ClusterIP)
			}
		})
	}
}

func TestCreateOrUpdateService(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	cr := platformMonitoring()
	log := testLogger()

	if err := CreateOrUpdateService(cr, clientset, false, false, log); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	got, err := clientset.CoreV1().Services(utils.EtcdServiceComponentNamespace).Get(context.TODO(), utils.EtcdServiceComponentName, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("get created service failed: %v", err)
	}
	if got.Labels["team"] != "monitoring" {
		t.Fatal("custom labels were not copied")
	}
	if got.Annotations["owner"] != "platform" {
		t.Fatal("custom annotations were not copied")
	}
	if _, ok := got.Spec.Selector["component"]; !ok {
		t.Fatalf("expected Kubernetes selector, got %#v", got.Spec.Selector)
	}

	got.Labels["stale"] = "true"
	if _, err := clientset.CoreV1().Services(utils.EtcdServiceComponentNamespace).Update(context.TODO(), got, metav1.UpdateOptions{}); err != nil {
		t.Fatalf("seed update failed: %v", err)
	}

	if err := CreateOrUpdateService(cr, clientset, true, false, log); err != nil {
		t.Fatalf("update failed: %v", err)
	}
	got, err = clientset.CoreV1().Services(utils.EtcdServiceComponentNamespace).Get(context.TODO(), utils.EtcdServiceComponentName, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("get updated service failed: %v", err)
	}
	if got.Labels["stale"] != "true" {
		t.Fatal("existing labels should be preserved")
	}
	if _, ok := got.Spec.Selector["openshift.io/component"]; !ok {
		t.Fatalf("expected OpenShift selector, got %#v", got.Spec.Selector)
	}
}

func TestCreateOrUpdateServiceOpenShiftV4Namespace(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	cr := platformMonitoring()

	if err := CreateOrUpdateService(cr, clientset, true, true, testLogger()); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	got, err := clientset.CoreV1().Services(utils.EtcdServiceComponentNamespaceOpenshiftV4).Get(context.TODO(), utils.EtcdServiceComponentName, metav1.GetOptions{})
	if err != nil {
		t.Fatalf("get created service failed: %v", err)
	}
	if _, ok := got.Spec.Selector["etcd"]; !ok {
		t.Fatalf("expected OpenShift v4 selector, got %#v", got.Spec.Selector)
	}
}

func TestCreateOrUpdateServiceMonitor(t *testing.T) {
	cl := newServiceMonitorClient(t)
	cr := platformMonitoring()
	log := testLogger()

	if err := createOrUpdateServiceMonitor(cr, cl, "monitoring", false, log); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	got := getServiceMonitor(t, cl, "created")
	assertServiceMonitorNamespace(t, got, utils.EtcdServiceComponentNamespace)
	assertCustomLabelsCopied(t, got.Labels)

	if err := createOrUpdateServiceMonitor(cr, cl, "monitoring", true, log); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	got = getServiceMonitor(t, cl, "updated")
	assertServiceMonitorNamespace(t, got, utils.EtcdServiceComponentNamespaceOpenshiftV4)
	if got.Spec.Endpoints[0].Port != "etcd-metrics" {
		t.Fatalf("got endpoint port %q, want etcd-metrics", got.Spec.Endpoints[0].Port)
	}
}

func newServiceMonitorClient(t *testing.T) client.Client {
	t.Helper()

	scheme := runtime.NewScheme()
	if err := monitoringv1.AddToScheme(scheme); err != nil {
		t.Fatalf("add PlatformMonitoring scheme failed: %v", err)
	}
	if err := promv1.AddToScheme(scheme); err != nil {
		t.Fatalf("add ServiceMonitor scheme failed: %v", err)
	}
	return ctrlfake.NewClientBuilder().WithScheme(scheme).Build()
}

func getServiceMonitor(t *testing.T, cl client.Client, operation string) *promv1.ServiceMonitor {
	t.Helper()

	key := client.ObjectKey{Name: "monitoring-etcd-service-monitor", Namespace: "monitoring"}
	sm := &promv1.ServiceMonitor{}
	if err := cl.Get(context.TODO(), key, sm); err != nil {
		t.Fatalf("get %s ServiceMonitor failed: %v", operation, err)
	}
	return sm
}

func assertServiceMonitorNamespace(t *testing.T, sm *promv1.ServiceMonitor, want string) {
	t.Helper()

	if sm.Spec.NamespaceSelector.MatchNames[0] != want {
		t.Fatalf("got namespace selector %#v", sm.Spec.NamespaceSelector.MatchNames)
	}
}

func assertCustomLabelsCopied(t *testing.T, labels map[string]string) {
	t.Helper()

	if labels["team"] != "monitoring" {
		t.Fatal("custom labels were not copied")
	}
}

func platformMonitoring() *monitoringv1.PlatformMonitoring {
	return &monitoringv1.PlatformMonitoring{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "monitoring.netcracker.com/v1",
			Kind:       "PlatformMonitoring",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        "platformmonitoring",
			Namespace:   "monitoring",
			UID:         types.UID("test-uid"),
			Labels:      map[string]string{"team": "monitoring"},
			Annotations: map[string]string{"owner": "platform"},
		},
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}
