package main

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
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

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}
