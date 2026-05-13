package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	etcdSourceNamespaceOpenshiftV4 = "openshift-etcd-operator"
	etcdSourceConfigmapOpenshiftV4 = "etcd-metric-serving-ca"
	etcdSourceSecretOpenshiftV4    = "etcd-metric-client"

	etcdServiceName                = "etcd"
	etcdServiceNamespace           = "kube-system"
	etcdPodLabelSelector           = "component=etcd"
	openshiftSecurityAPIGroup      = "security.openshift.io"
)

type appOptions struct {
	namespace  string
	secretName string
	logLevel   slog.Level
}

type kubeClients struct {
	clientset       kubernetes.Interface
	discoveryClient discovery.DiscoveryInterface
}

type certificateData struct {
	key string
	ca  string
	crt string
}

func main() {
	opts, err := parseOptions()
	log := newLogger(opts.logLevel)
	slog.SetDefault(log)

	if err != nil {
		log.Error("invalid command line options", "error", err)
		os.Exit(1)
	}

	clients, err := newKubeClients()
	if err != nil {
		log.Error("failed to initialize Kubernetes clients", "error", err)
		os.Exit(1)
	}

	if err := run(context.TODO(), clients, opts, log); err != nil {
		log.Error("etcd certificates synchronization failed", "error", err)
		os.Exit(1)
	}
}

func parseOptions() (appOptions, error) {
	var secretName string
	var namespace string
	var logLevel string
	flag.StringVar(&secretName, "secret", "kube-etcd-client-certs", "Name of the secret to create/update")
	flag.StringVar(&namespace, "namespace", "", "Namespace in which to create the Secret (defaults to NAMESPACE env)")
	flag.StringVar(&logLevel, "log-level", "info", "Log level: debug, info, warn, or error")
	flag.Parse()

	if namespace == "" {
		namespace = os.Getenv("NAMESPACE")
	}
	if namespace == "" {
		namespace = "monitoring"
	}

	level, err := parseLogLevel(logLevel)
	if err != nil {
		return appOptions{}, err
	}

	return appOptions{
		namespace:  namespace,
		secretName: secretName,
		logLevel:   level,
	}, nil
}

func parseLogLevel(level string) (slog.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("unsupported log level %q", level)
	}
}

func newLogger(level slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
}

func run(ctx context.Context, clients *kubeClients, opts appOptions, log *slog.Logger) error {
	isOpenshiftV4, err := hasOpenshiftSecurityAPI(clients.discoveryClient)
	if err != nil {
		return fmt.Errorf("couldn't probe OpenShift API: %w", err)
	}
	log.Info("Detected cluster flavor", "openshiftV4", isOpenshiftV4)

	certs, err := loadEtcdCerts(ctx, clients.clientset, isOpenshiftV4, log)
	if err != nil {
		return err
	}

	if err := certVerify(certs.key, certs.ca, certs.crt); err != nil {
		return fmt.Errorf("failed to verify etcd certificates: %w", err)
	}

	secret := buildEtcdSecret(opts.namespace, opts.secretName, certs)
	if err := createOrUpdateSecret(ctx, clients.clientset, secret, log); err != nil {
		return fmt.Errorf("secret operation failed for %s: %w", secret.Name, err)
	}

	if isOpenshiftV4 {
		log.Info("Skipping etcd Service creation on OpenShift v4 (openshift-etcd/etcd already exists)")
		return nil
	}
	if err := createOrUpdateEtcdService(ctx, clients.clientset, log); err != nil {
		return fmt.Errorf("failed to create/update etcd Service: %w", err)
	}
	return nil
}

func newKubeClients() (*kubeClients, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("couldn't get in-cluster config: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("couldn't get clientset: %w", err)
	}
	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("couldn't get discovery client: %w", err)
	}
	return &kubeClients{clientset: clientset, discoveryClient: dc}, nil
}

func hasOpenshiftSecurityAPI(dc discovery.DiscoveryInterface) (bool, error) {
	groups, err := dc.ServerGroups()
	if err != nil {
		return false, fmt.Errorf("failed to list API groups: %w", err)
	}
	for _, g := range groups.Groups {
		if g.Name == openshiftSecurityAPIGroup {
			return true, nil
		}
	}
	return false, nil
}

func loadEtcdCerts(ctx context.Context, clientset kubernetes.Interface, isOpenshiftV4 bool, log *slog.Logger) (certificateData, error) {
	if isOpenshiftV4 {
		certs, err := getCertsFromConfigmapAndSecret(ctx, clientset, log, etcdSourceNamespaceOpenshiftV4, etcdSourceConfigmapOpenshiftV4, etcdSourceSecretOpenshiftV4)
		if err != nil {
			logEtcdCertsError(log, err, "Failed to get etcd certificates from configmap and secret")
			return certificateData{}, err
		}
		log.Info("Extracted etcd certificates from configmap and secret (OpenShift v4)",
			"etcdNamespace", etcdSourceNamespaceOpenshiftV4,
			"configmap", etcdSourceConfigmapOpenshiftV4,
			"secret", etcdSourceSecretOpenshiftV4)
		return certs, nil
	}

	certs, err := getCertsFromHostpath(ctx, clientset, log, etcdServiceNamespace, etcdPodLabelSelector)
	if err != nil {
		logEtcdCertsError(log, err, "Failed to get etcd certificates from hostpath")
		return certificateData{}, err
	}
	return certs, nil
}

func logEtcdCertsError(log *slog.Logger, err error, message string) {
	if apierrors.IsForbidden(err) {
		log.Error("Unable to update etcd certificates due to a lack of permission to access the requested etcd resource.", "error", err)
		return
	}
	log.Error(message, "error", err)
}

func buildEtcdSecret(namespace string, secretName string, certs certificateData) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       secretName,
				"app.kubernetes.io/component":  "monitoring-etcd",
				"app.kubernetes.io/managed-by": "etcd-certs-to-secret",
			},
		},
		Data: map[string][]byte{
			"etcd-client-ca.crt": []byte(certs.ca),
			"etcd-client.crt":    []byte(certs.crt),
			"etcd-client.key":    []byte(certs.key),
		},
	}
}

func buildEtcdService() *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{Kind: "Service", APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      etcdServiceName,
			Namespace: etcdServiceNamespace,
			Labels: map[string]string{
				"k8s-app":                      "etcd",
				"app.kubernetes.io/name":       etcdServiceName,
				"app.kubernetes.io/component":  "etcd",
				"app.kubernetes.io/part-of":    "monitoring",
				"app.kubernetes.io/managed-by": "etcd-certs-to-secret",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "metrics",
					Protocol:   corev1.ProtocolTCP,
					Port:       2379,
					TargetPort: intstr.FromInt32(2379),
				},
			},
			Selector:  map[string]string{"component": "etcd"},
			ClusterIP: corev1.ClusterIPNone,
			Type:      corev1.ServiceTypeClusterIP,
		},
	}
}

// Retrieve etcd certificates paths from etcd pods' command line arguments.
func extractEtcdCertPaths(ctx context.Context, clientset kubernetes.Interface, log *slog.Logger, etcdNamespace string, etcdPodLabel string) (peerKey, caCrt, peerCrt string, err error) {
	peerKey = "/etc/kubernetes/pki/etcd/peer.key"
	caCrt = "/etc/kubernetes/pki/etcd/ca.crt"
	peerCrt = "/etc/kubernetes/pki/etcd/peer.crt"

	pods, err := clientset.CoreV1().Pods(etcdNamespace).List(ctx, metav1.ListOptions{
		LabelSelector: etcdPodLabel,
	})
	if err != nil || len(pods.Items) == 0 {
		log.Error("Failed to retrieve pods to get etcd certificates", "error", err)
		return "", "", "", podListError(err)
	}

	etcdPod, err := runningEtcdPod(pods.Items)
	if err != nil {
		log.Error("Failed to find etcd pods among pods to get etcd certificates", "error", err)
		return "", "", "", err
	}

	peerKey, caCrt, peerCrt = certPathsFromPod(etcdPod, peerKey, caCrt, peerCrt)
	return peerKey, caCrt, peerCrt, nil
}

func podListError(err error) error {
	if err != nil {
		return err
	}
	return fmt.Errorf("no etcd pods found")
}

func runningEtcdPod(pods []corev1.Pod) (corev1.Pod, error) {
	for _, pod := range pods {
		if pod.Status.Phase == corev1.PodRunning {
			return pod, nil
		}
	}
	return corev1.Pod{}, fmt.Errorf("no running etcd pods found")
}

func certPathsFromPod(pod corev1.Pod, peerKey string, caCrt string, peerCrt string) (string, string, string) {
	for _, container := range pod.Spec.Containers {
		if container.Name == "etcd" {
			return certPathsFromArgs(container.Command, peerKey, caCrt, peerCrt)
		}
	}
	return peerKey, caCrt, peerCrt
}

func certPathsFromArgs(args []string, peerKey string, caCrt string, peerCrt string) (string, string, string) {
	for _, arg := range args {
		name, value, found := strings.Cut(arg, "=")
		if !found {
			continue
		}
		switch name {
		case "--peer-key-file":
			peerKey = value
		case "--peer-trusted-ca-file":
			caCrt = value
		case "--peer-cert-file":
			peerCrt = value
		}
	}
	return peerKey, caCrt, peerCrt
}

func getCertsFromConfigmapAndSecret(ctx context.Context, clientset kubernetes.Interface, log *slog.Logger, etcdNamespace string, configmapName string, etcdCertsSourceSecret string) (certificateData, error) {
	configMap, err := clientset.CoreV1().ConfigMaps(etcdNamespace).Get(ctx, configmapName, metav1.GetOptions{})
	if err != nil {
		log.Error("Failed to get configmap", "namespace", etcdNamespace, "configmap", configmapName, "error", err)
		return certificateData{}, err
	}
	secret, err := clientset.CoreV1().Secrets(etcdNamespace).Get(ctx, etcdCertsSourceSecret, metav1.GetOptions{})
	if err != nil {
		log.Error("Failed to get secret", "namespace", etcdNamespace, "secret", etcdCertsSourceSecret, "error", err)
		return certificateData{}, err
	}
	return certificateData{
		key: string(secret.Data["tls.key"]),
		ca:  configMap.Data["ca-bundle.crt"],
		crt: string(secret.Data["tls.crt"]),
	}, nil
}

func getCertsFromHostpath(ctx context.Context, clientset kubernetes.Interface, log *slog.Logger, etcdNamespace string, etcdPodLabel string) (certificateData, error) {
	peerKey, caCrt, peerCrt, err := extractEtcdCertPaths(ctx, clientset, log, etcdNamespace, etcdPodLabel)
	if err != nil {
		log.Error("Failed to get etcd certificates paths from etcd pods arguments", "error", err)
		return certificateData{}, err
	}
	log.Info("Using etcd certificates paths from etcd pods", "peerKey", peerKey, "caCrt", caCrt, "peerCrt", peerCrt)

	caData, err := readFileToString(log, caCrt)
	if err != nil {
		return certificateData{}, err
	}
	keyData, err := readFileToString(log, peerKey)
	if err != nil {
		return certificateData{}, err
	}
	crtData, err := readFileToString(log, peerCrt)
	if err != nil {
		return certificateData{}, err
	}
	return certificateData{key: keyData, ca: caData, crt: crtData}, nil
}

func readFileToString(log *slog.Logger, filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Error("Failed to read file", "file", filename, "error", err)
		return "", err
	}
	return string(data), nil
}

func createOrUpdateSecret(ctx context.Context, clientset kubernetes.Interface, secret *corev1.Secret, log *slog.Logger) error {
	existing, err := clientset.CoreV1().Secrets(secret.Namespace).Get(ctx, secret.Name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			if _, err := clientset.CoreV1().Secrets(secret.Namespace).Create(ctx, secret, metav1.CreateOptions{}); err != nil {
				return fmt.Errorf("failed to create secret: %w", err)
			}
			log.Info("Secret created", "secret", secret.Name)
			return nil
		}
		return fmt.Errorf("failed to check secret existence: %w", err)
	}
	secret.ResourceVersion = existing.ResourceVersion
	if _, err := clientset.CoreV1().Secrets(secret.Namespace).Update(ctx, secret, metav1.UpdateOptions{}); err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}
	log.Info("Secret updated", "secret", secret.Name)
	return nil
}

func createOrUpdateEtcdService(ctx context.Context, clientset kubernetes.Interface, log *slog.Logger) error {
	desired := buildEtcdService()
	existing, err := clientset.CoreV1().Services(desired.Namespace).Get(ctx, desired.Name, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return fmt.Errorf("failed to check etcd service existence: %w", err)
		}
		if _, err := clientset.CoreV1().Services(desired.Namespace).Create(ctx, desired, metav1.CreateOptions{}); err != nil {
			return fmt.Errorf("failed to create etcd service: %w", err)
		}
		log.Info("Service created", "service", desired.Name)
		return nil
	}
	existing.Spec.Ports = desired.Spec.Ports
	existing.Spec.Selector = desired.Spec.Selector
	if existing.Labels == nil {
		existing.Labels = map[string]string{}
	}
	for k, v := range desired.Labels {
		existing.Labels[k] = v
	}
	if _, err := clientset.CoreV1().Services(desired.Namespace).Update(ctx, existing, metav1.UpdateOptions{}); err != nil {
		return fmt.Errorf("failed to update etcd service: %w", err)
	}
	log.Info("Service updated", "service", desired.Name)
	return nil
}

func certVerify(keyData string, caData string, crtData string) error {
	if len(keyData) == 0 || len(caData) == 0 || len(crtData) == 0 {
		return fmt.Errorf("failed to get etcd certificates content, empty certificate data")
	}

	peerKeyBeginIndex := strings.Index(keyData, "-----BEGIN PRIVATE KEY-----")
	if peerKeyBeginIndex == -1 {
		peerKeyBeginIndex = strings.Index(keyData, "-----BEGIN RSA PRIVATE KEY-----")
	}

	caCertBeginIndex := strings.Index(caData, "-----BEGIN CERTIFICATE-----")
	caCertEndIndex := strings.Index(caData, "-----END CERTIFICATE-----")

	peerCertBeginIndex := strings.Index(crtData, "-----BEGIN CERTIFICATE-----")
	peerCertEndIndex := strings.Index(crtData, "-----END CERTIFICATE-----")

	if peerKeyBeginIndex == -1 {
		return fmt.Errorf("failed to find private key header")
	}
	if caCertBeginIndex == -1 || caCertEndIndex == -1 || caCertBeginIndex > caCertEndIndex {
		return fmt.Errorf("invalid CA certificate format")
	}
	if peerCertBeginIndex == -1 || peerCertEndIndex == -1 || peerCertBeginIndex > peerCertEndIndex {
		return fmt.Errorf("invalid peer certificate format")
	}
	return nil
}
