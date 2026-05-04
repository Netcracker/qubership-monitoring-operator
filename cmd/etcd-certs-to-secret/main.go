package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"strconv"
	"strings"

	monitoringv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

//go:embed  assets/*.yaml
var assets embed.FS

const (
	nameLabelKey      = "name"
	appNameLabelKey   = "app.kubernetes.io/name"
	instanceLabelKey  = "app.kubernetes.io/instance"
	componentLabelKey = "app.kubernetes.io/component"
	etcdComponent     = "monitoring-etcd"
)

type appOptions struct {
	namespace  string
	secretName string
	logLevel   slog.Level
}

type kubeClients struct {
	client          client.Client
	clientset       *kubernetes.Clientset
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

	if err := run(context.TODO(), opts, log); err != nil {
		log.Error("etcd certificates synchronization failed", "error", err)
		os.Exit(1)
	}
}

func parseOptions() (appOptions, error) {
	var secretName string
	var logLevel string
	flag.StringVar(&secretName, "secret", "kube-etcd-client-certs", "Name of the secret to create/update")
	flag.StringVar(&logLevel, "log-level", "info", "Log level: debug, info, warn, or error")
	flag.Parse()

	namespace := "monitoring"
	if value, found := os.LookupEnv("WATCH_NAMESPACE"); found {
		namespace = value
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

func run(ctx context.Context, opts appOptions, log *slog.Logger) error {
	namespacedName := types.NamespacedName{
		Namespace: opts.namespace,
		Name:      "platformmonitoring",
	}

	clients, err := newKubeClients()
	if err != nil {
		return err
	}

	isOpenshift, err := hasRouteApi(clients.clientset)
	if err != nil {
		return fmt.Errorf("couldn't check if cluster has Route API: %w", err)
	}

	isOpenshiftV4, err := isOpenshiftV4(clients.discoveryClient, isOpenshift, log)
	if err != nil {
		return fmt.Errorf("couldn't check cluster version: %w", err)
	}

	certs, err := loadEtcdCerts(clients.clientset, isOpenshift, isOpenshiftV4, log)
	if err != nil {
		return err
	}

	if err := certVerify(certs.key, certs.ca, certs.crt, log); err != nil {
		return fmt.Errorf("failed to verify etcd certificates: %w", err)
	}

	customResourceInstance := &monitoringv1.PlatformMonitoring{}
	if err := clients.client.Get(ctx, namespacedName, customResourceInstance); err != nil {
		return fmt.Errorf("failed to get PlatformMonitoring custom resource: %w", err)
	}

	secret, err := buildEtcdSecret(customResourceInstance, opts.namespace, opts.secretName, certs)
	if err != nil {
		return err
	}

	return syncEtcdResources(customResourceInstance, clients, secret, opts.namespace, isOpenshift, isOpenshiftV4, log)
}

func newKubeClients() (*kubeClients, error) {
	scheme := runtime.NewScheme()

	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("failed to add core Kubernetes types to scheme: %w", err)
	}
	if err := monitoringv1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("failed to add PlatformMonitoring to scheme: %w", err)
	}
	if err := promv1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("failed to add promv1 to scheme: %w", err)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("couldn't get config: %w", err)
	}
	cl, err := client.New(cfg, client.Options{
		Scheme: scheme})
	if err != nil {
		return nil, fmt.Errorf("couldn't get client: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("couldn't get clientset: %w", err)
	}

	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("couldn't get discovery client: %w", err)
	}

	return &kubeClients{
		client:          cl,
		clientset:       clientset,
		discoveryClient: dc,
	}, nil
}

func loadEtcdCerts(clientset *kubernetes.Clientset, isOpenshift bool, isOpenshiftV4 bool, log *slog.Logger) (certificateData, error) {
	if isOpenshiftV4 {
		keyData, caData, crtData, err := getCertsFromConfigmapAndSecret(clientset, log, utils.EtcdCertificatesSourceNamespaceOpenshiftV4, utils.EtcdCertificatesSourceConfigmapOpenshiftV4, utils.EtcdCertificatesSourceSecretOpenshiftV4)
		if err != nil {
			logEtcdCertsError(log, err, "Failed to get etcd certificates from configmap and secret")
			return certificateData{}, err
		}
		log.Info("Extracting etcd certificates from configmap and secret (Openshift v4)", "etcdNamespace", utils.EtcdCertificatesSourceNamespaceOpenshiftV4, "EtcdCertsSourceConfigmap", utils.EtcdCertificatesSourceConfigmapOpenshiftV4, "etcdCertsSourceSecret", utils.EtcdCertificatesSourceSecretOpenshiftV4)
		return certificateData{key: keyData, ca: caData, crt: crtData}, nil
	}

	keyData, caData, crtData, err := getCertsFromHostpath(clientset, log, utils.EtcdServiceComponentNamespace, utils.EtcdPodLabelSelector, isOpenshift)
	if err != nil {
		logEtcdCertsError(log, err, "Failed to get etcd certificates from hostpath")
		return certificateData{}, err
	}
	return certificateData{key: keyData, ca: caData, crt: crtData}, nil
}

func logEtcdCertsError(log *slog.Logger, err error, message string) {
	if apierrors.IsForbidden(err) {
		log.Error("Unable to update etcd certificates due to a lack of permission to access the requested etcd resource.", "error", err)
		return
	}
	log.Error(message, "error", err)
}

func buildEtcdSecret(cr *monitoringv1.PlatformMonitoring, namespace string, secretName string, certs certificateData) (*corev1.Secret, error) {
	certData := make(map[string][]byte)
	certData["etcd-client-ca.crt"] = []byte(certs.ca)
	certData["etcd-client.crt"] = []byte(certs.crt)
	certData["etcd-client.key"] = []byte(certs.key)

	secret, err := etcdSecret(cr, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Data: certData,
	})
	if err != nil {
		return nil, fmt.Errorf("failed creating Secret manifest: %w", err)
	}
	return secret, nil
}

func syncEtcdResources(cr *monitoringv1.PlatformMonitoring, clients *kubeClients, secret *corev1.Secret, namespace string, isOpenshift bool, isOpenshiftV4 bool, log *slog.Logger) error {
	if err := createOrUpdateSecret(clients.clientset, secret, log); err != nil {
		return fmt.Errorf("secret operation failed for %s: %w", secret.Name, err)
	}

	if err := createOrUpdateServiceMonitor(cr, clients.client, namespace, isOpenshiftV4, log); err != nil {
		return fmt.Errorf("failed to create/update etcd ServiceMonitor: %w", err)
	}

	if err := CreateOrUpdateService(cr, clients.clientset, isOpenshift, isOpenshiftV4, log); err != nil {
		return fmt.Errorf("failed to create/update etcd Service: %w", err)
	}
	return nil
}

// Retrieve etcd certificates paths from etcd pods' command line arguments
func extractEtcdCertPaths(ctx context.Context, clientset *kubernetes.Clientset, log *slog.Logger, etcdNamespace string, etcdPodLabel string) (peerKey, caCrt, peerCrt string, err error) {
	// Default cert paths
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

func hasRouteApi(clientset *kubernetes.Clientset) (bool, error) {

	resources, err := clientset.Discovery().ServerPreferredResources()
	if err != nil {
		return false, fmt.Errorf("failed to get server api resources: %v", err)
	}
	for _, resourceList := range resources {
		if resourceList.GroupVersion == "route.openshift.io/v1" {
			return true, nil
		}
	}
	return false, nil
}

func isOpenshiftV4(dc discovery.DiscoveryInterface, isOpenshift bool, log *slog.Logger) (bool, error) {
	serverVersion, err := dc.ServerVersion()
	if err != nil {
		return false, fmt.Errorf("failed to get server version: %v", err)
	}
	log.Info("Server version", "minor", serverVersion.Minor)
	minor, err := strconv.Atoi(serverVersion.Minor)
	if err != nil {
		return false, fmt.Errorf("failed to convert minor server version %s to integer: %v", serverVersion.Minor, err)
	}
	return minor >= 18 && isOpenshift, nil
}

func getCertsFromConfigmapAndSecret(clientset *kubernetes.Clientset, log *slog.Logger, etcdNamespace string, configmapName string, etcdCertsSourceSecret string) (string, string, string, error) {
	configMap, err := clientset.CoreV1().ConfigMaps(etcdNamespace).Get(context.TODO(), configmapName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsForbidden(err) {
			log.Error("Failed to get configmap due to insufficient permissions", "namespace", etcdNamespace, "configmap", configmapName, "error", err)
		} else {
			log.Error("Failed to get configmap", "namespace", etcdNamespace, "configmap", configmapName, "error", err)
		}
		return "", "", "", err
	}
	caData := configMap.Data["ca-bundle.crt"]

	secret, err := clientset.CoreV1().Secrets(etcdNamespace).Get(context.TODO(), etcdCertsSourceSecret, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsForbidden(err) {
			log.Error(fmt.Sprintf("Failed to get secret %s (namespace: %s) to get etcd certificates", etcdCertsSourceSecret, etcdNamespace), "error", err)
		}
		return "", "", "", err
	}

	secretData := secret.Data
	crtData := string(secretData["tls.crt"])
	keyData := string(secretData["tls.key"])
	return keyData, caData, crtData, nil
}

func getCertsFromHostpath(clientset *kubernetes.Clientset, log *slog.Logger, etcdNamespace string, etcdPodLabel string, isOpenshift bool) (keyData string, caData string, crtData string, err error) {
	var peerKey, caCrt, peerCrt string
	if isOpenshift {
		peerKey = "/etc/etcd/peer.key"
		caCrt = "/etc/etcd/ca.crt"
		peerCrt = "/etc/etcd/peer.crt"
		log.Info("Using default Openshift prior to v4 etcd certificates paths", "peerKey", peerKey, "caCrt", caCrt, "peerCrt", peerCrt)
	} else {
		peerKey, caCrt, peerCrt, err = extractEtcdCertPaths(context.TODO(), clientset, log, etcdNamespace, etcdPodLabel)
		if err != nil {
			log.Error("Failed to get etcd certificates paths from etcd pods arguments", "error", err)
			return "", "", "", err
		}
		log.Info("Using etcd certificates paths from etcd pods", "peerKey", peerKey, "caCrt", caCrt, "peerCrt", peerCrt)
	}

	caData, err = readFileToString(log, caCrt)
	if err != nil {
		return "", "", "", err
	}

	keyData, err = readFileToString(log, peerKey)
	if err != nil {
		return "", "", "", err
	}

	crtData, err = readFileToString(log, peerCrt)
	if err != nil {
		return "", "", "", err
	}

	return keyData, caData, crtData, nil
}

func readFileToString(log *slog.Logger, filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Error("Failed to read file", "file", filename, "error", err)
		return "", err
	}
	return string(data), nil
}

func etcdServiceMonitor(cr *monitoringv1.PlatformMonitoring, namespace string, isOpenshiftV4 bool) (*promv1.ServiceMonitor, error) {
	sm := promv1.ServiceMonitor{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.EtcdServiceMonitorAsset), 100).Decode(&sm); err != nil {
		return nil, err
	}

	//Set parameters
	sm.SetGroupVersionKind(schema.GroupVersionKind{Group: "monitoring.coreos.com", Version: "v1", Kind: "ServiceMonitor"})
	sm.SetName(namespace + "-" + "etcd-service-monitor")
	sm.SetNamespace(namespace)
	if isOpenshiftV4 {
		sm.Spec.NamespaceSelector.MatchNames = []string{utils.EtcdServiceComponentNamespaceOpenshiftV4}
		// Port "etcd-metrics" is used in OpenShift v4.x
		sm.Spec.Endpoints[0].Port = "etcd-metrics"
	} else {
		sm.Spec.NamespaceSelector.MatchNames = []string{utils.EtcdServiceComponentNamespace}
	}

	if cr.Spec.KubernetesMonitors != nil {
		monitor, ok := cr.Spec.KubernetesMonitors[utils.EtcdServiceMonitorName]
		if ok && monitor.IsInstall() {
			monitor.OverrideServiceMonitor(&sm)
		}
	}

	if cr.GetLabels() != nil {
		maps.Copy(sm.Labels, cr.GetLabels())
	}

	sm.Labels[nameLabelKey] = utils.TruncLabel(sm.GetName())
	sm.Labels[appNameLabelKey] = utils.TruncLabel(sm.GetName())
	sm.Labels[instanceLabelKey] = utils.GetInstanceLabel(sm.GetName(), sm.GetNamespace())
	sm.Labels[componentLabelKey] = etcdComponent

	if sm.Annotations == nil && cr.GetAnnotations() != nil {
		sm.SetAnnotations(cr.GetAnnotations())
	} else {
		maps.Copy(sm.Annotations, cr.GetAnnotations())
	}
	sm.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: "monitoring.netcracker.com/v1",
			Kind:       "PlatformMonitoring",
			Name:       cr.Name,
			UID:        cr.UID,
			Controller: ptr.To(true),
		},
	}
	return &sm, nil
}

func createOrUpdateServiceMonitor(cr *monitoringv1.PlatformMonitoring, cl client.Client, namespace string, isOpenshiftV4 bool, log *slog.Logger) error {
	sm, err := etcdServiceMonitor(cr, namespace, isOpenshiftV4)
	if err != nil {
		log.Error("Failed creating ServiceMonitor manifest", "error", err)
	}

	existingSM := &promv1.ServiceMonitor{}
	err = cl.Get(context.TODO(), client.ObjectKey{Name: sm.Name, Namespace: sm.Namespace}, existingSM)
	if err != nil {
		if apierrors.IsNotFound(err) {

			if err := cl.Create(context.TODO(), sm); err != nil {
				return fmt.Errorf("failed to create etcd ServiceMonitor: %w", err)
			}
			log.Info("ServiceMonitor created", "name", sm.Name)
		} else {
			return fmt.Errorf("failed to check etcd servicemonitor existence: %w", err)
		}
	} else {

		sm.SetResourceVersion(existingSM.GetResourceVersion())
		if err := cl.Update(context.TODO(), sm); err != nil {
			return fmt.Errorf("failed to update etcd ServiceMonitor: %w", err)
		}
		log.Info("ServiceMonitor updated", "name", sm.Name)
	}
	return nil
}

func CreateOrUpdateService(cr *monitoringv1.PlatformMonitoring, clientset kubernetes.Interface, isOpenshift bool, isOpenshiftV4 bool, log *slog.Logger) error {
	etcdServiceNamespace := etcdServiceNamespace(isOpenshiftV4)

	m, err := etcdService(isOpenshift, etcdServiceNamespace, isOpenshiftV4)
	if err != nil {
		log.Error("Failed creating Service manifest", "error", err)
		return err
	}

	e, exists, err := getExistingEtcdService(clientset, etcdServiceNamespace, m.Name)
	if err != nil {
		return err
	}

	applyEtcdService(cr, e, m)
	return persistEtcdService(clientset, etcdServiceNamespace, e, exists, log)
}

func etcdServiceNamespace(isOpenshiftV4 bool) string {
	if isOpenshiftV4 {
		return utils.EtcdServiceComponentNamespaceOpenshiftV4
	}
	return utils.EtcdServiceComponentNamespace
}

func getExistingEtcdService(clientset kubernetes.Interface, namespace string, name string) (*corev1.Service, bool, error) {
	service, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		}, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("failed to get check if etcd service exists: %w", err)
	}
	return service, true, nil
}

func applyEtcdService(cr *monitoringv1.PlatformMonitoring, e *corev1.Service, m *corev1.Service) {
	e.TypeMeta = m.TypeMeta
	e.Spec.Ports = m.Spec.Ports
	e.Spec.Selector = m.Spec.Selector
	e.Spec.ClusterIP = m.Spec.ClusterIP

	if e.Labels == nil {
		e.Labels = make(map[string]string)
	}
	if cr.GetLabels() != nil {
		maps.Copy(e.Labels, cr.GetLabels())
	}

	if e.Annotations == nil {
		e.Annotations = make(map[string]string)
	}
	if cr.GetAnnotations() != nil {
		maps.Copy(e.Annotations, cr.GetAnnotations())
	}
	maps.Copy(e.Labels, m.Labels)
}

func persistEtcdService(clientset kubernetes.Interface, namespace string, service *corev1.Service, exists bool, log *slog.Logger) error {
	if !exists {
		return createEtcdService(clientset, namespace, service, log)
	}
	return updateEtcdService(clientset, namespace, service, log)
}

func createEtcdService(clientset kubernetes.Interface, namespace string, service *corev1.Service, log *slog.Logger) error {
	if _, err := clientset.CoreV1().Services(namespace).Create(context.TODO(), service, metav1.CreateOptions{}); err != nil {
		log.Error("Failed to create etcd service", "error", err, "service", service.Name)
		return err
	}
	log.Info("Service created", "service", service.Name)
	return nil
}

func updateEtcdService(clientset kubernetes.Interface, namespace string, service *corev1.Service, log *slog.Logger) error {
	if _, err := clientset.CoreV1().Services(namespace).Update(context.TODO(), service, metav1.UpdateOptions{}); err != nil {
		if apierrors.IsForbidden(err) {
			log.Info("Failed to update etcd service", "error", err, "service", service.Name)
			return err
		}
		log.Error("Failed to update etcd service", "error", err, "service", service.Name)
		return fmt.Errorf("failed to update etcd service: %w", err)
	}
	log.Info("Service updated", "service", service.Name)
	return nil
}

func etcdService(isOpenshift bool, etcdServiceNamespace string, isOpenshiftV4 bool) (*corev1.Service, error) {
	service := corev1.Service{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.EtcdServiceComponentAsset), 100).Decode(&service); err != nil {
		return nil, err
	}

	service.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Service"})
	service.SetName(utils.EtcdServiceComponentName)
	service.SetNamespace(etcdServiceNamespace)

	// Kubernetes uses "component: etcd" selector
	// OpenShift v3.x uses "openshift.io/component: etcd" selector
	if isOpenshift && !isOpenshiftV4 {
		service.Spec.Selector = map[string]string{"openshift.io/component": "etcd"}
	}
	// OpenShift v4.x uses "etcd: 'true'" selector
	if isOpenshiftV4 {
		service.Spec.Selector = map[string]string{"etcd": "true"}
	}
	// If cluster is not OpenShift v4.x, remove port "etcd-metrics"
	if !isOpenshiftV4 {
		service.Spec.Ports = service.Spec.Ports[:1]
	}
	service.Spec.ClusterIP = ""
	service.Labels[nameLabelKey] = utils.TruncLabel(service.GetName())
	service.Labels[appNameLabelKey] = utils.TruncLabel(service.GetName())
	service.Labels[instanceLabelKey] = utils.GetInstanceLabel(service.GetName(), service.GetNamespace())
	service.Labels[componentLabelKey] = etcdComponent

	return &service, nil
}

func createOrUpdateSecret(clientset kubernetes.Interface, secret *corev1.Secret, log *slog.Logger) error {
	_, err := clientset.CoreV1().Secrets(secret.Namespace).Get(context.TODO(), secret.Name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsForbidden(err) {
			return fmt.Errorf("failed to check secret existence due to insufficient permissions: %w", err)
		}
		if apierrors.IsNotFound(err) {

			_, err = clientset.CoreV1().Secrets(secret.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("failed to create secret: %w", err)
			}
			log.Info("Secret created", "secret", secret.Name)
			return nil
		}
		return fmt.Errorf("failed to check secret existence: %w", err)
	}

	_, err = clientset.CoreV1().Secrets(secret.Namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}
	log.Info("Secret updated", "secret", secret.Name)
	return nil
}

func certVerify(keyData string, caData string, crtData string, log *slog.Logger) error {
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
func etcdSecret(cr *monitoringv1.PlatformMonitoring, secret *corev1.Secret) (*corev1.Secret, error) {
	secret.Labels = make(map[string]string)
	secret.Annotations = make(map[string]string)
	//Set parameters
	if cr.GetLabels() != nil {
		for k, v := range cr.GetLabels() {
			if _, ok := secret.Labels[k]; !ok {
				secret.Labels[k] = v
			}
		}
	}
	secret.Labels[nameLabelKey] = secret.Name
	secret.Labels[appNameLabelKey] = utils.TruncLabel(secret.Name)
	secret.Labels[instanceLabelKey] = utils.GetInstanceLabel(secret.Name, secret.Namespace)
	secret.Labels[componentLabelKey] = etcdComponent
	if secret.Annotations == nil && cr.GetAnnotations() != nil {
		secret.SetAnnotations(cr.GetAnnotations())
	} else {
		maps.Copy(secret.Annotations, cr.GetAnnotations())
	}
	secret.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: "monitoring.netcracker.com/v1",
			Kind:       "PlatformMonitoring",
			Name:       cr.Name,
			UID:        cr.UID,
			Controller: ptr.To(true),
		},
	}
	return secret, nil
}
