package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NamespaceConfig holds the config for a specific namespace
type NamespaceConfig struct {
	Namespace      string   `json:"namespace"`
	SecretNames    []string `json:"secretNames"`
	ConfigMapNames []string `json:"configMapNames"`
}

type Config struct {
	SourceNamespaces     []string
	TargetNamespace      string
	ConfigMapName        string
	ConfigMapNamespace   string
	KeyNamespaces        string
	GlobalSecretNames    []string          `json:"globalSecretNames"`
	GlobalConfigMapNames []string          `json:"globalConfigMapNames"`
	NamespaceConfigs     []NamespaceConfig `json:"namespaceConfigs"`
}

func main() {
	// Base config with ConfigMap location
	baseConfig := Config{
		TargetNamespace:    "monitoring",
		ConfigMapName:      "namespace-list",
		ConfigMapNamespace: "monitoring",
		KeyNamespaces:      "namespaces",
	}

	clientset, err := newClientset()
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	// Load full config from ConfigMap
	config, err := loadConfigFromConfigMap(clientset, baseConfig)
	if err != nil {
		log.Fatalf("Failed to load config from ConfigMap: %v", err)
	}

	// Get source namespaces
	config.SourceNamespaces, err = getNamespacesFromConfigMap(clientset, config)
	if err != nil {
		log.Fatalf("Failed to get namespaces from ConfigMap: %v", err)
	}

	fmt.Printf("Found namespaces to process: %v\n", config.SourceNamespaces)
	fmt.Printf("Global secrets to copy: %v\n", config.GlobalSecretNames)
	fmt.Printf("Global configmaps to copy: %v\n", config.GlobalConfigMapNames)
	fmt.Printf("Per-namespace configs: %d\n", len(config.NamespaceConfigs))

	for _, ns := range config.SourceNamespaces {
		fmt.Printf("\n=== Processing namespace: %s ===\n", ns)

		// Get secret and configmap names for this namespace
		secretNames := getSecretNamesForNamespace(config, ns)
		configMapNames := getConfigMapNamesForNamespace(config, ns)

		fmt.Printf("  Secrets to copy: %v\n", secretNames)
		fmt.Printf("  ConfigMaps to copy: %v\n", configMapNames)

		if err := copySecrets(clientset, ns, config.TargetNamespace, secretNames); err != nil {
			log.Printf("Error copying secrets from %s: %v", ns, err)
		}

		if err := copyConfigMaps(clientset, ns, config.TargetNamespace, configMapNames); err != nil {
			log.Printf("Error copying configmaps from %s: %v", ns, err)
		}
	}

	fmt.Println("\n=== Processing complete ===")
}

func newClientset() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			home, _ := os.UserHomeDir()
			kubeconfig = filepath.Join(home, ".kube", "config")
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return clientset, nil
}

func loadConfigFromConfigMap(clientset *kubernetes.Clientset, baseConfig Config) (Config, error) {
	cm, err := clientset.CoreV1().ConfigMaps(baseConfig.ConfigMapNamespace).Get(
		context.Background(),
		baseConfig.ConfigMapName,
		metav1.GetOptions{},
	)
	if err != nil {
		return baseConfig, fmt.Errorf("failed to get ConfigMap %s/%s: %w", baseConfig.ConfigMapNamespace, baseConfig.ConfigMapName, err)
	}

	configJSON, ok := cm.Data["config"]
	if !ok {
		return baseConfig, nil
	}

	var config Config
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return baseConfig, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	if config.TargetNamespace == "" {
		config.TargetNamespace = baseConfig.TargetNamespace
	}
	if config.ConfigMapName == "" {
		config.ConfigMapName = baseConfig.ConfigMapName
	}
	if config.ConfigMapNamespace == "" {
		config.ConfigMapNamespace = baseConfig.ConfigMapNamespace
	}
	if config.KeyNamespaces == "" {
		config.KeyNamespaces = baseConfig.KeyNamespaces
	}

	return config, nil
}

func getNamespacesFromConfigMap(clientset *kubernetes.Clientset, config Config) ([]string, error) {
	cm, err := clientset.CoreV1().ConfigMaps(config.ConfigMapNamespace).Get(
		context.Background(),
		config.ConfigMapName,
		metav1.GetOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get ConfigMap %s/%s: %w", config.ConfigMapNamespace, config.ConfigMapName, err)
	}

	namespacesStr, ok := cm.Data[config.KeyNamespaces]
	if !ok {
		return nil, fmt.Errorf("key %q not found in ConfigMap", config.KeyNamespaces)
	}

	return parseNamespaces(namespacesStr), nil
}

func parseNamespaces(namespacesStr string) []string {
	var namespaces []string
	lines := strings.Split(namespacesStr, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				namespaces = append(namespaces, part)
			}
		}
	}

	return namespaces
}

func getSecretNamesForNamespace(config Config, ns string) []string {
	for _, nc := range config.NamespaceConfigs {
		if nc.Namespace == ns {
			if len(nc.SecretNames) > 0 {
				return nc.SecretNames
			}
		}
	}
	return config.GlobalSecretNames
}

func getConfigMapNamesForNamespace(config Config, ns string) []string {
	for _, nc := range config.NamespaceConfigs {
		if nc.Namespace == ns {
			if len(nc.ConfigMapNames) > 0 {
				return nc.ConfigMapNames
			}
		}
	}
	return config.GlobalConfigMapNames
}

func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if strings.TrimSpace(s) == str {
			return true
		}
	}
	return false
}

func copySecrets(clientset *kubernetes.Clientset, sourceNs, targetNs string, secretNames []string) error {
	secrets, err := clientset.CoreV1().Secrets(sourceNs).List(
		context.Background(),
		metav1.ListOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to list secrets: %w", err)
	}

	fmt.Printf("  Found %d secrets in %s\n", len(secrets.Items), sourceNs)

	for _, secret := range secrets.Items {
		if len(secretNames) > 0 && !containsString(secretNames, secret.Name) {
			fmt.Printf("    Skipping (not in list): %s\n", secret.Name)
			continue
		}

		if shouldSkipSecret(secret) {
			fmt.Printf("    Skipping: %s\n", secret.Name)
			continue
		}

		newSecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secret.Name,
				Namespace: targetNs,
			},
			Type: secret.Type,
			Data: secret.Data,
		}

		_, err := clientset.CoreV1().Secrets(targetNs).Create(
			context.Background(),
			newSecret,
			metav1.CreateOptions{},
		)
		if err != nil {
			_, updateErr := clientset.CoreV1().Secrets(targetNs).Update(
				context.Background(),
				newSecret,
				metav1.UpdateOptions{},
			)
			if updateErr != nil {
				fmt.Printf("    Error copying secret %s: %v\n", secret.Name, updateErr)
				continue
			}
			fmt.Printf("    Updated: %s\n", secret.Name)
			continue
		}
		fmt.Printf("    Copied: %s\n", secret.Name)
	}

	return nil
}

func copyConfigMaps(clientset *kubernetes.Clientset, sourceNs, targetNs string, configMapNames []string) error {
	configmaps, err := clientset.CoreV1().ConfigMaps(sourceNs).List(
		context.Background(),
		metav1.ListOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to list configmaps: %w", err)
	}

	fmt.Printf("  Found %d configmaps in %s\n", len(configmaps.Items), sourceNs)

	for _, cm := range configmaps.Items {
		if len(configMapNames) > 0 && !containsString(configMapNames, cm.Name) {
			fmt.Printf("    Skipping (not in list): %s\n", cm.Name)
			continue
		}

		if shouldSkipConfigMap(cm) {
			fmt.Printf("    Skipping: %s\n", cm.Name)
			continue
		}

		newCM := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cm.Name,
				Namespace: targetNs,
			},
			Data:       cm.Data,
			BinaryData: cm.BinaryData,
		}

		_, err := clientset.CoreV1().ConfigMaps(targetNs).Create(
			context.Background(),
			newCM,
			metav1.CreateOptions{},
		)
		if err != nil {
			_, updateErr := clientset.CoreV1().ConfigMaps(targetNs).Update(
				context.Background(),
				newCM,
				metav1.UpdateOptions{},
			)
			if updateErr != nil {
				fmt.Printf("    Error copying configmap %s: %v\n", cm.Name, updateErr)
				continue
			}
			fmt.Printf("    Updated: %s\n", cm.Name)
			continue
		}
		fmt.Printf("    Copied: %s\n", cm.Name)
	}

	return nil
}

func shouldSkipSecret(secret corev1.Secret) bool {
	if len(secret.Data) == 0 && len(secret.StringData) == 0 {
		return true
	}
	if secret.Type == corev1.SecretTypeServiceAccountToken {
		return true
	}
	return false
}

func shouldSkipConfigMap(cm corev1.ConfigMap) bool {
	if len(cm.Data) == 0 && len(cm.BinaryData) == 0 {
		return true
	}
	if cm.Name == "kube-root-ca.crt" {
		return true
	}
	return false
}
