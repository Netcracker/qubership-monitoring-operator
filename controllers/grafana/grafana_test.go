package grafana

import (
	"context"
	"fmt"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/labelsassert"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

var (
	cr              *monv1.PlatformMonitoring
	labelKey        = "label.key"
	labelValue      = "label-value"
	annotationKey   = "annotation.key"
	annotationValue = "annotation-value"
)

func TestGrafanaManifests(t *testing.T) {
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{
				Annotations: map[string]string{annotationKey: annotationValue},
				Labels:      map[string]string{labelKey: labelValue},
			},
		},
	}
	t.Run("Test Grafana manifest", func(t *testing.T) {
		m, err := grafana(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Grafana manifest should not be empty")
		assert.NotNil(t, m.Spec.Client)
		assert.True(t, m.Spec.Client.UseKubeAuth)
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		// In grafana-operator v5, Labels and Annotations are in Deployment.Spec.Template
		if m.Spec.Deployment != nil && m.Spec.Deployment.Spec.Template != nil {
			if m.Spec.Deployment.Spec.Template.Labels != nil {
				assert.Equal(t, labelValue, m.Spec.Deployment.Spec.Template.Labels[labelKey])
			}
			if m.Spec.Deployment.Spec.Template.Annotations != nil {
				assert.Equal(t, annotationValue, m.Spec.Deployment.Spec.Template.Annotations[annotationKey])
			}
		}
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{},
		},
	}
	// Disabled for v5: in v5 labels/annotations live in Deployment.Spec.Template, not Deployment
	//t.Run("Test Grafana manifest with nil annotation", func(t *testing.T) {
	//	m, err := grafana(cr)
	//	...
	//})
	t.Run("Test GrafanaDatasource manifest", func(t *testing.T) {
		m, err := grafanaDataSource(cr, nil, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "GrafanaDatasource manifest should not be empty")
	})
	t.Run("Test GrafanaPromxyDatasource manifest", func(t *testing.T) {
		m, err := grafanaPromxyDataSource(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "GrafanaPromxyDatasource manifest should not be empty")
		assert.Equal(t, "platform-monitoring-promxy", m.GetName())
		if m.Spec.Datasource != nil {
			assert.Contains(t, m.Spec.Datasource.URL, "promxy")
		}
	})
	t.Run("Test Ingress v1 manifest", func(t *testing.T) {
		m, err := grafanaIngressV1(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Ingress v1 manifest should not be empty")
	})
	t.Run("Test PodMonitor manifest", func(t *testing.T) {
		crWithLabels := &monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring", Labels: map[string]string{labelKey: labelValue}},
			Spec: monv1.PlatformMonitoringSpec{
				Grafana: &monv1.Grafana{Labels: map[string]string{labelKey: labelValue}},
			},
		}
		m, err := grafanaPodMonitor(crWithLabels)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "PodMonitor manifest should not be empty")
		labelsassert.AssertCRLabels(t, m.GetLabels(), utils.GrafanaComponentName, "victoriametrics-operator", map[string]string{labelKey: labelValue})
	})
}

func TestComputeAdminPasswordChecksum(t *testing.T) {
	base := map[string][]byte{
		"GF_SECURITY_ADMIN_USER":     []byte("admin"),
		"GF_SECURITY_ADMIN_PASSWORD": []byte("secret"),
	}
	renamedUser := map[string][]byte{
		"GF_SECURITY_ADMIN_USER":     []byte("other-admin"),
		"GF_SECURITY_ADMIN_PASSWORD": []byte("secret"),
	}
	changedPassword := map[string][]byte{
		"GF_SECURITY_ADMIN_USER":     []byte("admin"),
		"GF_SECURITY_ADMIN_PASSWORD": []byte("new-secret"),
	}

	assert.Equal(t, computeAdminPasswordChecksum(base), computeAdminPasswordChecksum(renamedUser))
	assert.NotEqual(t, computeAdminPasswordChecksum(base), computeAdminPasswordChecksum(changedPassword))
	assert.Empty(t, computeAdminPasswordChecksum(map[string][]byte{"GF_SECURITY_ADMIN_USER": []byte("admin")}))
}

func TestGrafanaManifestAdminPasswordChecksumAnnotation(t *testing.T) {
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{},
		},
	}

	m, err := grafanaWithAdminPasswordChecksum(cr, "checksum-value")
	require.NoError(t, err)
	require.NotNil(t, m.Spec.Deployment)
	require.NotNil(t, m.Spec.Deployment.Spec.Template)
	assert.Equal(t, "checksum-value", m.Spec.Deployment.Spec.Template.Annotations[adminSecretChecksumAnnotation])

	m, err = grafana(cr)
	require.NoError(t, err)
	require.NotNil(t, m.Spec.Deployment)
	require.NotNil(t, m.Spec.Deployment.Spec.Template)
	assert.NotContains(t, m.Spec.Deployment.Spec.Template.Annotations, adminSecretChecksumAnnotation)
}

func TestGrafanaManifestCustomNameNamespaceAndSecretMounts(t *testing.T) {
	disableDefaultSecret := true
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Auth: &monv1.Auth{},
			Grafana: &monv1.Grafana{
				Name:                      "custom-grafana",
				Namespace:                 "grafana-ns",
				DisableDefaultAdminSecret: &disableDefaultSecret,
			},
		},
	}

	m, err := grafanaWithAdminPasswordChecksum(cr, "password-checksum")
	require.NoError(t, err)
	assert.Equal(t, "custom-grafana", m.GetName())
	assert.Equal(t, "grafana-ns", m.GetNamespace())
	require.NotNil(t, m.Spec.Deployment)
	require.NotNil(t, m.Spec.Deployment.Spec.Template)
	podSpec := m.Spec.Deployment.Spec.Template.Spec
	require.NotNil(t, podSpec)
	require.NotEmpty(t, podSpec.Containers)

	assert.Equal(t, "password-checksum", m.Spec.Deployment.Spec.Template.Annotations[adminSecretChecksumAnnotation])
	assert.Equal(t, "$__file{/etc/grafana-admin/GF_SECURITY_ADMIN_USER}", m.Spec.Config["security"]["admin_user"])
	assert.Equal(t, "$__file{/etc/grafana-admin/GF_SECURITY_ADMIN_PASSWORD}", m.Spec.Config["security"]["admin_password"])
	assert.Equal(t, "$__file{/etc/grafana-oauth/GF_AUTH_GENERIC_OAUTH_CLIENT_SECRET}", m.Spec.Config["auth.generic_oauth"]["client_secret"])

	adminVolume := findVolume(t, podSpec.Volumes, "grafana-admin-secret")
	require.NotNil(t, adminVolume.Secret)
	assert.Equal(t, "custom-grafana-admin-credentials", adminVolume.Secret.SecretName)
	require.NotNil(t, adminVolume.Secret.Optional)
	assert.True(t, *adminVolume.Secret.Optional)

	oauthVolume := findVolume(t, podSpec.Volumes, "grafana-oauth-secret")
	require.NotNil(t, oauthVolume.Secret)
	assert.Equal(t, "grafana-oauth-client-secret", oauthVolume.Secret.SecretName)
	require.NotNil(t, oauthVolume.Secret.Optional)
	assert.True(t, *oauthVolume.Secret.Optional)

	container := podSpec.Containers[0]
	assert.Contains(t, container.VolumeMounts, corev1.VolumeMount{
		Name:      "grafana-admin-secret",
		MountPath: "/etc/grafana-admin",
		ReadOnly:  true,
	})
	assert.Contains(t, container.VolumeMounts, corev1.VolumeMount{
		Name:      "grafana-oauth-secret",
		MountPath: "/etc/grafana-oauth",
		ReadOnly:  true,
	})
}

func TestGrafanaManifestDeploymentOptions(t *testing.T) {
	replicas := int32(2)
	runAsUser := int64(10001)
	runAsGroup := int64(10002)
	fsGroup := int64(10003)
	disableDefaultSecret := false
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Auth: &monv1.Auth{},
			Grafana: &monv1.Grafana{
				Image:                "ghcr.io/grafana/grafana:11.2.0",
				GrafanaHomeDashboard: true,
				Operator: monv1.GrafanaOperator{
					InitContainerImage: "ghcr.io/netcracker/grafana-plugins:latest",
				},
				Replicas: &replicas,
				SecurityContext: &monv1.SecurityContext{
					RunAsUser:  &runAsUser,
					RunAsGroup: &runAsGroup,
					FSGroup:    &fsGroup,
				},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{corev1.ResourceCPU: resource.MustParse("100m")},
					Limits:   corev1.ResourceList{corev1.ResourceMemory: resource.MustParse("256Mi")},
				},
				Tolerations: []corev1.Toleration{{
					Key:      "monitoring",
					Operator: corev1.TolerationOpExists,
				}},
				NodeSelector: map[string]string{"node-role": "monitoring"},
				Affinity: &corev1.Affinity{
					NodeAffinity: &corev1.NodeAffinity{},
				},
				Labels: map[string]string{
					"dashboards": "custom",
				},
				Annotations: map[string]string{
					"grafana.example.com/restarted-by": "test",
				},
				PriorityClassName:         "monitoring-priority",
				DisableDefaultAdminSecret: &disableDefaultSecret,
			},
		},
	}

	m, err := grafanaWithAdminPasswordChecksum(cr, "password-checksum")
	require.NoError(t, err)
	require.NotNil(t, m.Spec.Deployment)
	assert.Equal(t, "ghcr.io/grafana/grafana:11.2.0", m.Spec.Version)
	assert.Equal(t, &replicas, m.Spec.Deployment.Spec.Replicas)
	assert.Equal(t, "custom", m.Labels["dashboards"])
	assert.Equal(t, "custom", m.Spec.Deployment.Spec.Template.Labels["dashboards"])
	assert.Equal(t, "test", m.Spec.Deployment.Spec.Template.Annotations["grafana.example.com/restarted-by"])

	podSpec := m.Spec.Deployment.Spec.Template.Spec
	require.NotNil(t, podSpec)
	require.NotEmpty(t, podSpec.Containers)
	assert.Equal(t, "monitoring-priority", podSpec.PriorityClassName)
	assert.Equal(t, map[string]string{"node-role": "monitoring"}, podSpec.NodeSelector)
	assert.Equal(t, cr.Spec.Grafana.Tolerations, podSpec.Tolerations)
	assert.Equal(t, cr.Spec.Grafana.Affinity, podSpec.Affinity)
	require.NotNil(t, podSpec.SecurityContext)
	assert.Equal(t, &runAsUser, podSpec.SecurityContext.RunAsUser)
	assert.Equal(t, &runAsGroup, podSpec.SecurityContext.RunAsGroup)
	assert.Equal(t, &fsGroup, podSpec.SecurityContext.FSGroup)

	container := podSpec.Containers[0]
	assert.Equal(t, cr.Spec.Grafana.Resources, container.Resources)
	assert.Contains(t, container.EnvFrom, corev1.EnvFromSource{
		ConfigMapRef: &corev1.ConfigMapEnvSource{
			LocalObjectReference: corev1.LocalObjectReference{Name: "grafana-extra-vars"},
		},
	})
	assert.Contains(t, container.EnvFrom, corev1.EnvFromSource{
		SecretRef: &corev1.SecretEnvSource{
			LocalObjectReference: corev1.LocalObjectReference{Name: "grafana-extra-vars-secret"},
		},
	})
	assert.Contains(t, container.VolumeMounts, corev1.VolumeMount{
		Name:      "configmap-grafana-home-dashboard",
		MountPath: "/etc/grafana-configmaps/grafana-home-dashboard",
		ReadOnly:  true,
	})
	assert.Contains(t, container.VolumeMounts, corev1.VolumeMount{
		Name:      "grafana-plugins",
		MountPath: "/var/lib/grafana/plugins",
	})

	homeDashboardVolume := findVolume(t, podSpec.Volumes, "configmap-grafana-home-dashboard")
	require.NotNil(t, homeDashboardVolume.ConfigMap)
	assert.Equal(t, "grafana-home-dashboard", homeDashboardVolume.ConfigMap.Name)
	pluginsVolume := findVolume(t, podSpec.Volumes, "grafana-plugins")
	assert.NotNil(t, pluginsVolume.EmptyDir)
	adminVolume := findVolume(t, podSpec.Volumes, "grafana-admin-secret")
	require.NotNil(t, adminVolume.Secret)
	require.NotNil(t, adminVolume.Secret.Optional)
	assert.False(t, *adminVolume.Secret.Optional)

	require.Len(t, podSpec.InitContainers, 1)
	assert.Equal(t, "grafana-plugins-init", podSpec.InitContainers[0].Name)
	assert.Equal(t, "ghcr.io/netcracker/grafana-plugins:latest", podSpec.InitContainers[0].Image)
}

func TestHandleGrafanaCreatesResource(t *testing.T) {
	cr := &monv1.PlatformMonitoring{
		TypeMeta: metav1.TypeMeta{
			APIVersion: monv1.SchemeGroupVersion.String(),
			Kind:       "PlatformMonitoring",
		},
		ObjectMeta: metav1.ObjectMeta{Name: "platform-monitoring", Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{
				Name:      "custom-grafana",
				Namespace: "grafana-ns",
			},
		},
	}
	scheme := newGrafanaTestScheme(t)
	reconciler := &GrafanaReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: ctrlfake.NewClientBuilder().WithScheme(scheme).Build(),
			Scheme: scheme,
			Log:    utils.Logger("grafana_test"),
		},
	}

	require.NoError(t, reconciler.handleGrafana(cr, "password-checksum"))

	created := &grafv1.Grafana{}
	require.NoError(t, reconciler.Client.Get(t.Context(), client.ObjectKey{Name: "custom-grafana", Namespace: "grafana-ns"}, created))
	require.NotNil(t, created.Spec.Deployment)
	require.NotNil(t, created.Spec.Deployment.Spec.Template)
	assert.Equal(t, "password-checksum", created.Spec.Deployment.Spec.Template.Annotations[adminSecretChecksumAnnotation])
	assert.Equal(t, "custom-grafana", created.GetName())
	assert.Equal(t, "grafana-ns", created.GetNamespace())
}

func TestHandleGrafanaReturnsGetError(t *testing.T) {
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Name: "platform-monitoring", Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{},
		},
	}
	scheme := newGrafanaTestScheme(t)
	reconciler := &GrafanaReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: ctrlfake.NewClientBuilder().
				WithScheme(scheme).
				WithInterceptorFuncs(interceptor.Funcs{
					Get: func(ctx context.Context, c client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
						return fmt.Errorf("get failed")
					},
				}).
				Build(),
			Scheme: scheme,
			Log:    utils.Logger("grafana_test"),
		},
	}

	require.ErrorContains(t, reconciler.handleGrafana(cr, ""), "get failed")
}

func TestHandleGrafanaReturnsCreateError(t *testing.T) {
	cr := &monv1.PlatformMonitoring{
		TypeMeta: metav1.TypeMeta{
			APIVersion: monv1.SchemeGroupVersion.String(),
			Kind:       "PlatformMonitoring",
		},
		ObjectMeta: metav1.ObjectMeta{Name: "platform-monitoring", Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{},
		},
	}
	scheme := newGrafanaTestScheme(t)
	reconciler := &GrafanaReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: ctrlfake.NewClientBuilder().
				WithScheme(scheme).
				WithInterceptorFuncs(interceptor.Funcs{
					Create: func(ctx context.Context, c client.WithWatch, obj client.Object, opts ...client.CreateOption) error {
						return fmt.Errorf("create failed")
					},
				}).
				Build(),
			Scheme: scheme,
			Log:    utils.Logger("grafana_test"),
		},
	}

	require.ErrorContains(t, reconciler.handleGrafana(cr, ""), "create failed")
}

func TestGrafanaCredentialHelpers(t *testing.T) {
	defaultCR := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
	}
	assert.Equal(t, "grafana-admin-credentials", getGrafanaAdminSecretName(defaultCR))
	assert.Equal(t, "monitoring", getGrafanaNamespace(defaultCR))
	assert.Equal(t, "grafana", getGrafanaName(defaultCR))

	customCR := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{
				Name:      "custom-grafana",
				Namespace: "grafana-ns",
			},
		},
	}
	assert.Equal(t, "custom-grafana-admin-credentials", getGrafanaAdminSecretName(customCR))
	assert.Equal(t, "grafana-ns", getGrafanaNamespace(customCR))
	assert.Equal(t, "custom-grafana", getGrafanaName(customCR))
}

func TestHandleGrafanaCredentialsSecretUsesCustomNameAndNamespace(t *testing.T) {
	passwordChecksum := computeAdminPasswordChecksum(map[string][]byte{
		"GF_SECURITY_ADMIN_PASSWORD": []byte("secret"),
	})
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{
				Name:      "custom-grafana",
				Namespace: "grafana-ns",
			},
		},
	}
	scheme := newGrafanaTestScheme(t)
	reconciler := &GrafanaReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: ctrlfake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{Name: "custom-grafana-admin-credentials", Namespace: "grafana-ns"},
					Data: map[string][]byte{
						"GF_SECURITY_ADMIN_PASSWORD": []byte("secret"),
					},
				}).
				Build(),
			Scheme: scheme,
			Log:    utils.Logger("grafana_test"),
		},
	}

	state, err := reconciler.handleGrafanaCredentialsSecret(cr)
	require.NoError(t, err)
	assert.Equal(t, passwordChecksum, state.adminPasswordChecksum)
	assert.False(t, state.resetPassword)
}

func TestHandleGrafanaCredentialsSecretPasswordOnlySync(t *testing.T) {
	passwordChecksum := computeAdminPasswordChecksum(map[string][]byte{
		"GF_SECURITY_ADMIN_PASSWORD": []byte("secret"),
	})
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{},
		},
	}
	scheme := newGrafanaTestScheme(t)
	reconciler := &GrafanaReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: ctrlfake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(
					&corev1.Secret{
						ObjectMeta: metav1.ObjectMeta{Name: "grafana-admin-credentials", Namespace: "monitoring"},
						Data: map[string][]byte{
							"GF_SECURITY_ADMIN_USER":     []byte("renamed-admin"),
							"GF_SECURITY_ADMIN_PASSWORD": []byte("secret"),
						},
					},
					grafanaWithChecksum("grafana", "monitoring", passwordChecksum),
				).
				Build(),
			Scheme: scheme,
			Log:    utils.Logger("grafana_test"),
		},
	}

	state, err := reconciler.handleGrafanaCredentialsSecret(cr)
	require.NoError(t, err)
	assert.Equal(t, passwordChecksum, state.adminPasswordChecksum)
	assert.False(t, state.resetPassword)
}

func TestHandleGrafanaCredentialsSecretDetectsPasswordChange(t *testing.T) {
	oldChecksum := computeAdminPasswordChecksum(map[string][]byte{
		"GF_SECURITY_ADMIN_PASSWORD": []byte("old-secret"),
	})
	newChecksum := computeAdminPasswordChecksum(map[string][]byte{
		"GF_SECURITY_ADMIN_PASSWORD": []byte("new-secret"),
	})
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{},
		},
	}
	scheme := newGrafanaTestScheme(t)
	reconciler := &GrafanaReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: ctrlfake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(
					&corev1.Secret{
						ObjectMeta: metav1.ObjectMeta{Name: "grafana-admin-credentials", Namespace: "monitoring"},
						Data: map[string][]byte{
							"GF_SECURITY_ADMIN_USER":     []byte("admin"),
							"GF_SECURITY_ADMIN_PASSWORD": []byte("new-secret"),
						},
					},
					grafanaWithChecksum("grafana", "monitoring", oldChecksum),
				).
				Build(),
			Scheme: scheme,
			Log:    utils.Logger("grafana_test"),
		},
	}

	state, err := reconciler.handleGrafanaCredentialsSecret(cr)
	require.NoError(t, err)
	assert.Equal(t, newChecksum, state.adminPasswordChecksum)
	assert.True(t, state.resetPassword)
}

func TestHandleGrafanaCredentialsSecretFirstInstallDoesNotReset(t *testing.T) {
	newChecksum := computeAdminPasswordChecksum(map[string][]byte{
		"GF_SECURITY_ADMIN_PASSWORD": []byte("new-secret"),
	})
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{},
		},
	}
	scheme := newGrafanaTestScheme(t)
	reconciler := &GrafanaReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: ctrlfake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{Name: "grafana-admin-credentials", Namespace: "monitoring"},
					Data: map[string][]byte{
						"GF_SECURITY_ADMIN_USER":     []byte("admin"),
						"GF_SECURITY_ADMIN_PASSWORD": []byte("new-secret"),
					},
				}).
				Build(),
			Scheme: scheme,
			Log:    utils.Logger("grafana_test"),
		},
	}

	state, err := reconciler.handleGrafanaCredentialsSecret(cr)
	require.NoError(t, err)
	assert.Equal(t, newChecksum, state.adminPasswordChecksum)
	assert.False(t, state.resetPassword)
}

func TestHandleGrafanaCredentialsSecretMissingPassword(t *testing.T) {
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{},
		},
	}
	scheme := newGrafanaTestScheme(t)
	reconciler := &GrafanaReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: ctrlfake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{Name: "grafana-admin-credentials", Namespace: "monitoring"},
					Data: map[string][]byte{
						"GF_SECURITY_ADMIN_USER": []byte("admin"),
					},
				}).
				Build(),
			Scheme: scheme,
			Log:    utils.Logger("grafana_test"),
		},
	}

	state, err := reconciler.handleGrafanaCredentialsSecret(cr)
	require.NoError(t, err)
	assert.Empty(t, state.adminPasswordChecksum)
	assert.False(t, state.resetPassword)
}

func TestHandleGrafanaCredentialsSecretMissingSecret(t *testing.T) {
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{},
		},
	}
	scheme := newGrafanaTestScheme(t)
	reconciler := &GrafanaReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: ctrlfake.NewClientBuilder().WithScheme(scheme).Build(),
			Scheme: scheme,
			Log:    utils.Logger("grafana_test"),
		},
	}

	state, err := reconciler.handleGrafanaCredentialsSecret(cr)
	require.NoError(t, err)
	assert.Empty(t, state.adminPasswordChecksum)
	assert.False(t, state.resetPassword)
}

func TestResetGrafanaCredentialsSkipsMissingSecret(t *testing.T) {
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{},
		},
	}
	scheme := newGrafanaTestScheme(t)
	reconciler := &GrafanaReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: ctrlfake.NewClientBuilder().WithScheme(scheme).Build(),
			Scheme: scheme,
			Log:    utils.Logger("grafana_test"),
		},
	}

	require.NoError(t, reconciler.resetGrafanaCredentials(cr))
}

func TestResetGrafanaCredentialsSkipsEmptyPassword(t *testing.T) {
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{},
		},
	}
	scheme := newGrafanaTestScheme(t)
	reconciler := &GrafanaReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: ctrlfake.NewClientBuilder().
				WithScheme(scheme).
				WithObjects(&corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{Name: "grafana-admin-credentials", Namespace: "monitoring"},
					Data: map[string][]byte{
						"GF_SECURITY_ADMIN_PASSWORD": []byte(""),
					},
				}).
				Build(),
			Scheme: scheme,
			Log:    utils.Logger("grafana_test"),
		},
	}

	require.NoError(t, reconciler.resetGrafanaCredentials(cr))
}

func TestWaitForReadyGrafanaPod(t *testing.T) {
	reconciler := &GrafanaReconciler{
		KubeClient: k8sfake.NewSimpleClientset(&corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "custom-grafana-pod",
				Namespace: "grafana-ns",
				Labels: map[string]string{
					"app.kubernetes.io/name": "custom-grafana",
				},
			},
			Status: corev1.PodStatus{
				Phase: corev1.PodRunning,
				Conditions: []corev1.PodCondition{
					{Type: corev1.PodReady, Status: corev1.ConditionTrue},
				},
			},
		}),
	}

	podName, err := reconciler.waitForReadyGrafanaPod("grafana-ns", "custom-grafana")
	require.NoError(t, err)
	assert.Equal(t, "custom-grafana-pod", podName)
}

func TestWaitForReadyGrafanaPodListError(t *testing.T) {
	kubeClient := k8sfake.NewSimpleClientset()
	kubeClient.PrependReactor("list", "pods", func(action k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, fmt.Errorf("list failed")
	})
	reconciler := &GrafanaReconciler{KubeClient: kubeClient}

	_, err := reconciler.waitForReadyGrafanaPod("grafana-ns", "custom-grafana")
	require.ErrorContains(t, err, "cannot list Grafana pods")
}

func TestIsGrafanaPodReady(t *testing.T) {
	assert.True(t, isGrafanaPodReady(&corev1.Pod{
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
			Conditions: []corev1.PodCondition{
				{Type: corev1.PodReady, Status: corev1.ConditionTrue},
			},
		},
	}))
	assert.False(t, isGrafanaPodReady(&corev1.Pod{
		Status: corev1.PodStatus{Phase: corev1.PodPending},
	}))
	assert.False(t, isGrafanaPodReady(&corev1.Pod{
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
			Conditions: []corev1.PodCondition{
				{Type: corev1.PodReady, Status: corev1.ConditionFalse},
			},
		},
	}))
}

func newGrafanaTestScheme(t *testing.T) *runtime.Scheme {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, grafv1.AddToScheme(scheme))
	require.NoError(t, monv1.AddToScheme(scheme))
	return scheme
}

func grafanaWithChecksum(name, namespace, checksum string) *grafv1.Grafana {
	return &grafv1.Grafana{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: grafv1.GrafanaSpec{
			Deployment: &grafv1.DeploymentV1{
				Spec: grafv1.DeploymentV1Spec{
					Template: &grafv1.DeploymentV1PodTemplateSpec{
						ObjectMeta: grafv1.ObjectMeta{
							Annotations: map[string]string{adminSecretChecksumAnnotation: checksum},
						},
					},
				},
			},
		},
	}
}

func findVolume(t *testing.T, volumes []corev1.Volume, name string) corev1.Volume {
	t.Helper()
	for _, volume := range volumes {
		if volume.Name == name {
			return volume
		}
	}
	t.Fatalf("volume %s not found", name)
	return corev1.Volume{}
}
