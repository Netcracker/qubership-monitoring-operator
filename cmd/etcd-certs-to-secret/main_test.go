package main

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	qubershiporgv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
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

func TestCreateOrUpdateServiceMonitor(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := qubershiporgv1.AddToScheme(scheme); err != nil {
		t.Fatalf("add PlatformMonitoring scheme failed: %v", err)
	}
	if err := promv1.AddToScheme(scheme); err != nil {
		t.Fatalf("add ServiceMonitor scheme failed: %v", err)
	}

	cl := ctrlfake.NewClientBuilder().WithScheme(scheme).Build()
	cr := platformMonitoring()
	log := testLogger()

	if err := createOrUpdateServiceMonitor(cr, cl, "monitoring", false, log); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	key := client.ObjectKey{Name: "monitoring-etcd-service-monitor", Namespace: "monitoring"}
	got := &promv1.ServiceMonitor{}
	if err := cl.Get(context.TODO(), key, got); err != nil {
		t.Fatalf("get created ServiceMonitor failed: %v", err)
	}
	if got.Spec.NamespaceSelector.MatchNames[0] != utils.EtcdServiceComponentNamespace {
		t.Fatalf("got namespace selector %#v", got.Spec.NamespaceSelector.MatchNames)
	}
	if got.Labels["team"] != "monitoring" {
		t.Fatal("custom labels were not copied")
	}

	if err := createOrUpdateServiceMonitor(cr, cl, "monitoring", true, log); err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if err := cl.Get(context.TODO(), key, got); err != nil {
		t.Fatalf("get updated ServiceMonitor failed: %v", err)
	}
	if got.Spec.NamespaceSelector.MatchNames[0] != utils.EtcdServiceComponentNamespaceOpenshiftV4 {
		t.Fatalf("got namespace selector %#v", got.Spec.NamespaceSelector.MatchNames)
	}
	if got.Spec.Endpoints[0].Port != "etcd-metrics" {
		t.Fatalf("got endpoint port %q, want etcd-metrics", got.Spec.Endpoints[0].Port)
	}
}

func platformMonitoring() *qubershiporgv1.PlatformMonitoring {
	return &qubershiporgv1.PlatformMonitoring{
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
