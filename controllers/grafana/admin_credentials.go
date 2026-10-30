package grafana

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"maps"
	"strings"
	"time"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	grafanaAdminPasswordKey = "GF_SECURITY_ADMIN_PASSWORD"

	// adminSecretChecksumAnnotation on the Grafana pod template holds the admin password checksum,
	// so grafana-operator restarts Grafana when the password changes.
	adminSecretChecksumAnnotation = "checksum/admin-secret"

	// adminPasswordSyncedAnnotation on the Grafana resource holds the checksum of the password that was
	// last written to the Grafana database. It is set only after a successful reset, so a failed reset
	// is retried on the next reconcile.
	adminPasswordSyncedAnnotation = "monitoring.netcracker.com/grafana-admin-password-checksum"

	grafanaContainerName = "grafana"
)

// Grafana pod readiness polling for the admin password reset. Variables so tests can shorten them.
var (
	grafanaPodReadyTimeout = 5 * time.Minute
	grafanaPodPollInterval = 5 * time.Second
)

// adminPasswordChecksum returns the SHA-256 hex digest of the admin password, or "" when the
// Secret has no password. Only the password is hashed: Grafana CLI can reset the password in an
// existing database, but it cannot rename the admin user.
func adminPasswordChecksum(secret *corev1.Secret) string {
	password := secret.Data[grafanaAdminPasswordKey]
	if len(password) == 0 {
		return ""
	}
	sum := sha256.Sum256(password)
	return hex.EncodeToString(sum[:])
}

// readAdminPasswordChecksum returns the checksum of the password in the admin credentials Secret.
// A missing Secret or password yields "" without an error.
func (r *GrafanaReconciler) readAdminPasswordChecksum(cr *monv1.PlatformMonitoring) (string, error) {
	secret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{
		Name:      getGrafanaAdminSecretName(cr),
		Namespace: grafanaNamespace(cr),
	}}
	if err := r.GetResource(secret); err != nil {
		if errors.IsNotFound(err) {
			return "", nil
		}
		return "", fmt.Errorf("cannot read Grafana admin credentials Secret %s/%s: %w",
			secret.Namespace, secret.Name, err)
	}
	return adminPasswordChecksum(secret), nil
}

// addAdminPasswordChecksum sets the admin password checksum on the desired pod template. When the
// Secret has no password, it keeps the checksum from the existing pod template to avoid a restart.
// On creation, it also marks the password as synced: Grafana creates the admin user from the
// mounted Secret on first start, so no reset is needed.
func (r *GrafanaReconciler) addAdminPasswordChecksum(
	cr *monv1.PlatformMonitoring,
	manifest *grafv1.Grafana,
	existingAnnotations map[string]string,
	creating bool,
) error {
	checksum, err := r.readAdminPasswordChecksum(cr)
	if err != nil {
		return err
	}
	if checksum == "" {
		checksum = existingAnnotations[adminSecretChecksumAnnotation]
	}
	if checksum == "" {
		return nil
	}
	template := manifest.Spec.Deployment.Spec.Template
	if template.Annotations == nil {
		template.Annotations = make(map[string]string)
	}
	template.Annotations[adminSecretChecksumAnnotation] = checksum
	if creating {
		// Copy first: the manifest may share its annotation map with the PlatformMonitoring spec.
		annotations := maps.Clone(manifest.GetAnnotations())
		if annotations == nil {
			annotations = make(map[string]string, 1)
		}
		annotations[adminPasswordSyncedAnnotation] = checksum
		manifest.SetAnnotations(annotations)
	}
	return nil
}

// syncGrafanaAdminPassword writes the password from the admin credentials Secret to the Grafana
// database when it differs from the last synced password. Grafana applies admin_password only
// when it creates the admin user, so a changed Secret alone does not change the password of an
// existing database, for example one on a persistent volume.
func (r *GrafanaReconciler) syncGrafanaAdminPassword(cr *monv1.PlatformMonitoring) error {
	checksum, err := r.readAdminPasswordChecksum(cr)
	if err != nil || checksum == "" {
		return err
	}

	current := &grafv1.Grafana{}
	current.SetName(grafanaName(cr))
	current.SetNamespace(grafanaNamespace(cr))
	current.SetGroupVersionKind(schema.GroupVersionKind{Group: "grafana.integreatly.org", Version: "v1beta1", Kind: "Grafana"})
	if err = r.GetResource(current); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if current.GetAnnotations()[adminPasswordSyncedAnnotation] == checksum {
		return nil
	}

	r.Log.Info("Grafana admin password changed; resetting it in the Grafana database",
		"grafana", current.GetName(), "namespace", current.GetNamespace())
	if err = r.resetGrafanaCredentials(cr); err != nil {
		return err
	}

	patch := client.MergeFrom(current.DeepCopy())
	metav1.SetMetaDataAnnotation(&current.ObjectMeta, adminPasswordSyncedAnnotation, checksum)
	if err = r.Client.Patch(context.TODO(), current, patch); err != nil {
		return fmt.Errorf("cannot record synced Grafana admin password checksum: %w", err)
	}
	return nil
}

// resetGrafanaCredentials runs `grafana cli admin reset-admin-password` in a ready Grafana pod
// with the password from the admin credentials Secret.
func (r *GrafanaReconciler) resetGrafanaCredentials(cr *monv1.PlatformMonitoring) error {
	namespace := grafanaNamespace(cr)
	secret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: getGrafanaAdminSecretName(cr), Namespace: namespace}}
	if err := r.GetResource(secret); err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("Grafana admin credentials Secret not found; skipping password reset", "secret", secret.Name)
			return nil
		}
		return err
	}
	password := string(secret.Data[grafanaAdminPasswordKey])
	if password == "" {
		r.Log.Info("Grafana admin credentials Secret has no password; skipping password reset", "secret", secret.Name)
		return nil
	}

	podName, err := r.waitForReadyGrafanaPod(namespace, grafanaName(cr))
	if err != nil {
		return err
	}
	r.Log.Info("Resetting Grafana admin password", "pod", podName)

	// The password goes through stdin so it does not appear in the exec request URL or API server audit logs.
	command := []string{
		"grafana", "cli",
		"--homepath", "/usr/share/grafana",
		"--config", "/etc/grafana/grafana.ini",
		"admin", "reset-admin-password", "--password-from-stdin",
	}
	runCommand := r.runPodCommand
	if runCommand == nil {
		runCommand = r.execInPod
	}
	stdout, stderr, err := runCommand(namespace, podName, grafanaContainerName, command, password+"\n")
	if err != nil {
		return fmt.Errorf("grafana cli reset-admin-password failed in pod %s: %w; stdout: %s; stderr: %s",
			podName, err, stdout, stderr)
	}
	r.Log.Info("Grafana admin password reset", "pod", podName)
	return nil
}

// podCommandRunner runs a command in a pod container with the given stdin and returns its output.
type podCommandRunner func(namespace, pod, container string, command []string, stdin string) (string, string, error)

// execInPod runs a command in a pod container through the pods/exec subresource.
func (r *GrafanaReconciler) execInPod(namespace, pod, container string, command []string, stdin string) (string, string, error) {
	req := r.KubeClient.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(pod).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command:   command,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
		}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(r.config, "POST", req.URL())
	if err != nil {
		return "", "", fmt.Errorf("cannot create exec client for pod %s: %w", pod, err)
	}
	var stdout, stderr bytes.Buffer
	err = executor.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdin:  strings.NewReader(stdin),
		Stdout: &stdout,
		Stderr: &stderr,
	})
	return stdout.String(), stderr.String(), err
}

// waitForReadyGrafanaPod returns the name of a ready, non-terminating Grafana pod, polling until
// grafanaPodReadyTimeout expires.
func (r *GrafanaReconciler) waitForReadyGrafanaPod(namespace, name string) (string, error) {
	selector := "app.kubernetes.io/name=" + utils.TruncLabel(name)
	deadline := time.Now().Add(grafanaPodReadyTimeout)
	for {
		pods, err := r.KubeClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: selector})
		if err != nil {
			return "", fmt.Errorf("cannot list Grafana pods: %w", err)
		}
		for i := range pods.Items {
			if pods.Items[i].DeletionTimestamp == nil && isPodReady(&pods.Items[i]) {
				return pods.Items[i].Name, nil
			}
		}
		if time.Now().After(deadline) {
			return "", fmt.Errorf("no ready Grafana pod with label %s in namespace %s after %s",
				selector, namespace, grafanaPodReadyTimeout)
		}
		time.Sleep(grafanaPodPollInterval)
	}
}

func isPodReady(pod *corev1.Pod) bool {
	if pod.Status.Phase != corev1.PodRunning {
		return false
	}
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}
