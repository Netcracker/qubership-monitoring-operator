package grafana

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func secretKey(name, key string) *corev1.SecretKeySelector {
	return &corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: name}, Key: key}
}

func volumeByName(volumes []corev1.Volume, name string) *corev1.Volume {
	for i := range volumes {
		if volumes[i].Name == name {
			return &volumes[i]
		}
	}
	return nil
}

func countVolumesByName(volumes []corev1.Volume, name string) int {
	count := 0
	for _, volume := range volumes {
		if volume.Name == name {
			count++
		}
	}
	return count
}

func TestGrafanaManifestGenericOAuthConfig(t *testing.T) {
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Auth: &monv1.Auth{
				LoginURL:    "https://idp.example.com/auth",
				TokenURL:    "https://idp.example.com/token",
				UserInfoURL: "https://idp.example.com/userinfo",
				TLSConfig: &monv1.TLSConfig{
					InsecureSkipVerify: ptr(false),
					CASecret:           secretKey("idp.ca.secret", "ca.crt"),
					CertSecret:         secretKey("idp-client-secret", "tls.crt"),
					KeySecret:          secretKey("idp-client-secret", "tls.key"),
				},
			},
			Grafana: &monv1.Grafana{},
		},
	}

	m, err := grafanaWithDefaultSources(cr)
	require.NoError(t, err)

	authSection := m.Spec.Config[grafanaGenericOAuthSection]
	require.NotNil(t, authSection)
	assert.Equal(t, map[string]string{
		"enabled":                  "true",
		"auth_url":                 "https://idp.example.com/auth",
		"token_url":                "https://idp.example.com/token",
		"api_url":                  "https://idp.example.com/userinfo",
		"scopes":                   "openid profile",
		"client_secret":            "$__file{/etc/grafana-oauth/GF_AUTH_GENERIC_OAUTH_CLIENT_SECRET}",
		"tls_skip_verify_insecure": "false",
		"tls_client_ca":            "/etc/grafana-tls/idp.ca.secret/ca.crt",
		"tls_client_cert":          "/etc/grafana-tls/idp-client-secret/tls.crt",
		"tls_client_key":           "/etc/grafana-tls/idp-client-secret/tls.key",
	}, authSection)

	podSpec := m.Spec.Deployment.Spec.Template.Spec
	require.NotNil(t, podSpec)

	caVolumeName := grafanaTLSVolumeName("idp.ca.secret")
	clientVolumeName := grafanaTLSVolumeName("idp-client-secret")
	assert.NotContains(t, caVolumeName, ".")
	assert.Len(t, caVolumeName, len("grafana-tls-")+16)

	caVolume := volumeByName(podSpec.Volumes, caVolumeName)
	require.NotNil(t, caVolume)
	require.NotNil(t, caVolume.Secret)
	assert.Equal(t, "idp.ca.secret", caVolume.Secret.SecretName)

	clientVolume := volumeByName(podSpec.Volumes, clientVolumeName)
	require.NotNil(t, clientVolume)
	require.NotNil(t, clientVolume.Secret)
	assert.Equal(t, "idp-client-secret", clientVolume.Secret.SecretName)
	assert.Equal(t, 1, countVolumesByName(podSpec.Volumes, clientVolumeName))

	mounts := podSpec.Containers[0].VolumeMounts
	assert.Contains(t, mounts, corev1.VolumeMount{
		Name: caVolumeName, MountPath: "/etc/grafana-tls/idp.ca.secret", ReadOnly: true,
	})
	assert.Contains(t, mounts, corev1.VolumeMount{
		Name: clientVolumeName, MountPath: "/etc/grafana-tls/idp-client-secret", ReadOnly: true,
	})
}

func TestGrafanaManifestGenericOAuthSkipsIncompleteTLSSelectors(t *testing.T) {
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Auth: &monv1.Auth{
				LoginURL: "https://idp.example.com/auth",
				TLSConfig: &monv1.TLSConfig{
					CASecret:   secretKey("idp-ca", ""),
					CertSecret: secretKey("", "tls.crt"),
				},
			},
			Grafana: &monv1.Grafana{},
		},
	}

	m, err := grafana(cr, grafanaCredentialSources{AdminSecretPresent: true}, false)
	require.NoError(t, err)

	authSection := m.Spec.Config[grafanaGenericOAuthSection]
	assert.Equal(t, "true", authSection["enabled"])
	assert.NotContains(t, authSection, "client_secret")
	assert.NotContains(t, authSection, "tls_skip_verify_insecure")
	assert.NotContains(t, authSection, "tls_client_ca")
	assert.NotContains(t, authSection, "tls_client_cert")
	assert.Nil(t, volumeByName(m.Spec.Deployment.Spec.Template.Spec.Volumes, grafanaTLSVolumeName("idp-ca")))
}

func TestGrafanaManifestWithoutOAuthEndpoints(t *testing.T) {
	for name, auth := range map[string]*monv1.Auth{
		"nil auth":   nil,
		"empty auth": {},
		"blank urls": {LoginURL: " ", TokenURL: "\t", UserInfoURL: ""},
	} {
		t.Run(name, func(t *testing.T) {
			cr := &monv1.PlatformMonitoring{
				ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
				Spec:       monv1.PlatformMonitoringSpec{Auth: auth, Grafana: &monv1.Grafana{}},
			}
			m, err := grafana(cr, grafanaCredentialSources{AdminSecretPresent: true}, false)
			require.NoError(t, err)
			assert.Nil(t, m.Spec.Config[grafanaGenericOAuthSection])
		})
	}
}

func TestGrafanaManifestRootURL(t *testing.T) {
	tests := []struct {
		name    string
		ingress *monv1.Ingress
		want    string
	}{
		{name: "no ingress"},
		{name: "ingress not installed", ingress: &monv1.Ingress{Install: ptr(false), Host: "grafana.example.com"}},
		{name: "ingress without host", ingress: &monv1.Ingress{Install: ptr(true)}},
		{
			name:    "plain ingress",
			ingress: &monv1.Ingress{Install: ptr(true), Host: "grafana.example.com"},
			want:    "http://grafana.example.com/",
		},
		{
			name:    "ingress with TLS secret",
			ingress: &monv1.Ingress{Install: ptr(true), Host: "grafana.example.com", TLSSecretName: "grafana-tls"},
			want:    "https://grafana.example.com/",
		},
		{
			name: "ingress with TLS list",
			ingress: &monv1.Ingress{
				Install: ptr(true),
				Host:    "grafana.example.com",
				TLS:     []monv1.IngressTLS{{Hosts: []string{"grafana.example.com"}}},
			},
			want: "https://grafana.example.com/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := &monv1.PlatformMonitoring{
				ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
				Spec:       monv1.PlatformMonitoringSpec{Grafana: &monv1.Grafana{Ingress: tt.ingress}},
			}
			m, err := grafanaWithDefaultSources(cr)
			require.NoError(t, err)
			if tt.want == "" {
				assert.NotContains(t, m.Spec.Config[grafanaServerSection], "root_url")
				return
			}
			assert.Equal(t, tt.want, m.Spec.Config[grafanaServerSection]["root_url"])
		})
	}
}

func TestGrafanaManifestUserConfig(t *testing.T) {
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Auth: &monv1.Auth{LoginURL: "https://idp.example.com/auth"},
			Grafana: &monv1.Grafana{
				Ingress: &monv1.Ingress{Install: ptr(true), Host: "grafana.example.com"},
				// The chart default config mixes strings, booleans, and numbers.
				Config: &runtime.RawExtension{Raw: []byte(`{
					"server": {"root_url": "https://custom.example.com/grafana"},
					"auth": {"disable_login_form": false, "disable_signout_menu": true},
					"auth.generic_oauth": {"scopes": "openid email profile"},
					"live": {"max_connections": 0},
					"dataproxy": {"timeout": 1.5, "unset": null},
					"log": {"level": "warn"},
					"security": {"admin_user": "plain", "admin_password": "plain", "cookie_secure": true},
					"auth.jwt": {"enabled": false}
				}`)},
			},
		},
	}

	m, err := grafanaWithDefaultSources(cr)
	require.NoError(t, err)

	assert.Equal(t, "https://custom.example.com/grafana", m.Spec.Config[grafanaServerSection]["root_url"])
	assert.Equal(t, "false", m.Spec.Config["auth"]["disable_login_form"])
	assert.Equal(t, "true", m.Spec.Config["auth"]["disable_signout_menu"])
	assert.Equal(t, "openid email profile", m.Spec.Config[grafanaGenericOAuthSection]["scopes"])
	assert.Equal(t, "https://idp.example.com/auth", m.Spec.Config[grafanaGenericOAuthSection]["auth_url"])
	assert.Equal(t, "0", m.Spec.Config["live"]["max_connections"])
	assert.Equal(t, map[string]string{"timeout": "1.5"}, m.Spec.Config["dataproxy"])
	assert.Equal(t, "warn", m.Spec.Config["log"]["level"])

	// Reserved keys keep the operator-managed values.
	security := m.Spec.Config[grafanaSecuritySection]
	assert.Equal(t, "$__file{/etc/grafana-admin/GF_SECURITY_ADMIN_USER}", security["admin_user"])
	assert.Equal(t, "$__file{/etc/grafana-admin/GF_SECURITY_ADMIN_PASSWORD}", security["admin_password"])
	assert.Equal(t, "true", security["cookie_secure"])
	assert.Equal(t, "true", m.Spec.Config[grafanaJWTSection]["enabled"])
}

func TestGrafanaManifestRejectsInvalidUserConfig(t *testing.T) {
	for name, raw := range map[string]string{
		"not an object":  `["server"]`,
		"section scalar": `{"server": "root_url"}`,
		"nested object":  `{"server": {"root_url": {"value": "x"}}}`,
		"array value":    `{"server": {"domains": ["a", "b"]}}`,
	} {
		t.Run(name, func(t *testing.T) {
			cr := &monv1.PlatformMonitoring{
				ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
				Spec: monv1.PlatformMonitoringSpec{
					Grafana: &monv1.Grafana{Config: &runtime.RawExtension{Raw: []byte(raw)}},
				},
			}
			_, err := grafanaWithDefaultSources(cr)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "spec.grafana.config")
		})
	}
}

func TestApplyGrafanaUserConfigIgnoresEmptyConfig(t *testing.T) {
	m, err := grafanaWithDefaultSources(&monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec:       monv1.PlatformMonitoringSpec{Grafana: &monv1.Grafana{}},
	})
	require.NoError(t, err)
	before := len(m.Spec.Config)

	require.NoError(t, applyGrafanaUserConfig(m, nil))
	require.NoError(t, applyGrafanaUserConfig(m, &runtime.RawExtension{}))
	assert.Len(t, m.Spec.Config, before)
}

func TestMountGrafanaTLSSecretIsIdempotent(t *testing.T) {
	m, err := grafanaWithDefaultSources(&monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec:       monv1.PlatformMonitoringSpec{Grafana: &monv1.Grafana{}},
	})
	require.NoError(t, err)
	podSpec := m.Spec.Deployment.Spec.Template.Spec
	container := &podSpec.Containers[0]

	first := mountGrafanaTLSSecret(podSpec, container, "idp-ca")
	volumes, mounts := len(podSpec.Volumes), len(container.VolumeMounts)
	second := mountGrafanaTLSSecret(podSpec, container, "idp-ca")

	assert.Equal(t, "/etc/grafana-tls/idp-ca", first)
	assert.Equal(t, first, second)
	assert.Len(t, podSpec.Volumes, volumes)
	assert.Len(t, container.VolumeMounts, mounts)
}

func TestHasGrafanaOAuthConfig(t *testing.T) {
	assert.False(t, hasGrafanaOAuthConfig(nil))
	assert.False(t, hasGrafanaOAuthConfig(&monv1.Auth{}))
	assert.True(t, hasGrafanaOAuthConfig(&monv1.Auth{TokenURL: "https://idp.example.com/token"}))
	assert.True(t, hasGrafanaOAuthConfig(&monv1.Auth{UserInfoURL: "https://idp.example.com/userinfo"}))
}

func TestGrafanaRootURLDefaultsToHTTP(t *testing.T) {
	assert.Equal(t, "http://grafana.example.com/", grafanaRootURL("", "grafana.example.com"))
}
