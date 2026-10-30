package grafana

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	grafanaGenericOAuthSection = "auth.generic_oauth"
	grafanaJWTSection          = "auth.jwt"
	grafanaSecuritySection     = "security"
	grafanaServerSection       = "server"

	grafanaGenericOAuthScopes = "openid profile"
	grafanaTLSMountBase       = "/etc/grafana-tls"
)

// grafanaReservedConfigKeys lists spec.grafana.config keys that are not copied to grafana.ini.
// Helm writes config.security.admin_user and admin_password into the admin credentials Secret,
// which Grafana reads through $__file{}. Copying them would put the admin credentials in plain
// text into the Grafana custom resource. The whole auth.jwt section is skipped separately because
// grafana-operator reaches the Grafana API through it.
var grafanaReservedConfigKeys = map[string]map[string]bool{
	grafanaSecuritySection: {"admin_user": true, "admin_password": true},
}

// grafanaRootURL returns the Grafana root_url for the given protocol and host.
func grafanaRootURL(protocol string, host string) string {
	if protocol == "" {
		protocol = "http"
	}
	return fmt.Sprintf("%v://%v/", protocol, host)
}

// grafanaTLSVolumeName returns a valid volume name for a Secret name. Secret names may contain
// dots, which volume names (DNS labels) do not allow, so the name is derived from a hash.
func grafanaTLSVolumeName(secretName string) string {
	sum := sha256.Sum256([]byte(secretName))
	return fmt.Sprintf("grafana-tls-%x", sum[:8])
}

// hasGrafanaOAuthConfig reports whether spec.auth configures at least one IdP endpoint.
func hasGrafanaOAuthConfig(auth *monv1.Auth) bool {
	if auth == nil {
		return false
	}
	return strings.TrimSpace(auth.LoginURL) != "" ||
		strings.TrimSpace(auth.TokenURL) != "" ||
		strings.TrimSpace(auth.UserInfoURL) != ""
}

// configureGrafanaGenericOAuth maps spec.auth to the auth.generic_oauth section of grafana.ini
// and mounts the TLS Secrets that the section references.
func configureGrafanaGenericOAuth(graf *grafv1.Grafana, podSpec *grafv1.DeploymentV1PodSpec, container *corev1.Container, auth *monv1.Auth) {
	authSection := ensureGrafanaConfigSection(graf, grafanaGenericOAuthSection)
	authSection["enabled"] = "true"
	authSection["auth_url"] = auth.LoginURL
	authSection["token_url"] = auth.TokenURL
	authSection["api_url"] = auth.UserInfoURL
	authSection["scopes"] = grafanaGenericOAuthScopes

	tlsConfig := auth.TLSConfig
	if tlsConfig == nil {
		return
	}
	if tlsConfig.InsecureSkipVerify != nil {
		authSection["tls_skip_verify_insecure"] = strconv.FormatBool(*tlsConfig.InsecureSkipVerify)
	}
	tlsFiles := []struct {
		key      string
		selector *corev1.SecretKeySelector
	}{
		{key: "tls_client_ca", selector: tlsConfig.CASecret},
		{key: "tls_client_cert", selector: tlsConfig.CertSecret},
		{key: "tls_client_key", selector: tlsConfig.KeySecret},
	}
	for _, tlsFile := range tlsFiles {
		if tlsFile.selector == nil || tlsFile.selector.Name == "" || tlsFile.selector.Key == "" {
			continue
		}
		mountPath := mountGrafanaTLSSecret(podSpec, container, tlsFile.selector.Name)
		authSection[tlsFile.key] = mountPath + "/" + tlsFile.selector.Key
	}
}

// mountGrafanaTLSSecret mounts the Secret once, read-only, and returns its mount path.
func mountGrafanaTLSSecret(podSpec *grafv1.DeploymentV1PodSpec, container *corev1.Container, secretName string) string {
	volumeName := grafanaTLSVolumeName(secretName)
	mountPath := grafanaTLSMountBase + "/" + secretName

	if !hasVolume(podSpec.Volumes, volumeName) {
		podSpec.Volumes = append(podSpec.Volumes, corev1.Volume{
			Name: volumeName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{SecretName: secretName},
			},
		})
	}
	if !hasVolumeMount(container.VolumeMounts, volumeName) {
		container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
			Name:      volumeName,
			MountPath: mountPath,
			ReadOnly:  true,
		})
	}
	return mountPath
}

func hasVolume(volumes []corev1.Volume, name string) bool {
	for _, volume := range volumes {
		if volume.Name == name {
			return true
		}
	}
	return false
}

func hasVolumeMount(mounts []corev1.VolumeMount, name string) bool {
	for _, mount := range mounts {
		if mount.Name == name {
			return true
		}
	}
	return false
}

// applyGrafanaRootURL sets server.root_url from the Grafana Ingress host. The value is a default:
// spec.grafana.config can override it.
func applyGrafanaRootURL(graf *grafv1.Grafana, ingress *monv1.Ingress) {
	if ingress == nil || !ingress.IsInstall() || ingress.Host == "" {
		return
	}
	protocol := "http"
	if len(ingress.TLS) > 0 || ingress.TLSSecretName != "" {
		protocol = "https"
	}
	ensureGrafanaConfigSection(graf, grafanaServerSection)["root_url"] = grafanaRootURL(protocol, ingress.Host)
}

// applyGrafanaUserConfig merges spec.grafana.config into the Grafana custom resource spec.config.
// User values override operator defaults, except for the auth.jwt section and the keys in
// grafanaReservedConfigKeys. Scalar values are converted to strings because grafana.ini has no
// types; nested objects and arrays are rejected.
func applyGrafanaUserConfig(graf *grafv1.Grafana, config *runtime.RawExtension) error {
	if config == nil || len(config.Raw) == 0 {
		return nil
	}
	decoder := json.NewDecoder(bytes.NewReader(config.Raw))
	decoder.UseNumber()
	var userConfig map[string]map[string]any
	if err := decoder.Decode(&userConfig); err != nil {
		return fmt.Errorf("cannot decode spec.grafana.config: each section must be an object of scalar values: %w", err)
	}
	for section, values := range userConfig {
		if section == grafanaJWTSection {
			continue
		}
		for key, value := range values {
			if grafanaReservedConfigKeys[section][key] {
				continue
			}
			stringValue, ok, err := grafanaConfigValue(value)
			if err != nil {
				return fmt.Errorf("cannot use spec.grafana.config %s.%s: %w", section, key, err)
			}
			if ok {
				ensureGrafanaConfigSection(graf, section)[key] = stringValue
			}
		}
	}
	return nil
}

// grafanaConfigValue converts a decoded JSON scalar to its grafana.ini form. It returns false
// for null values, which leave the key unset.
func grafanaConfigValue(value any) (string, bool, error) {
	switch v := value.(type) {
	case nil:
		return "", false, nil
	case string:
		return v, true, nil
	case bool:
		return strconv.FormatBool(v), true, nil
	case json.Number:
		return v.String(), true, nil
	default:
		return "", false, fmt.Errorf("value must be a string, number, or boolean, got %T", value)
	}
}
