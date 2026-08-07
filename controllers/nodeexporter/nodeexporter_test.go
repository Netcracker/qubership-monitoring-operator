package nodeexporter

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/labelsassert"
	secv1 "github.com/openshift/api/security/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	cr              *monv1.PlatformMonitoring
	labelKey        = "label.key"
	labelValue      = "label-value"
	annotationKey   = "annotation.key"
	annotationValue = "annotation-value"
)

func TestNodeExporterManifests(t *testing.T) {
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			NodeExporter: &monv1.NodeExporter{
				Annotations: map[string]string{annotationKey: annotationValue},
				Labels:      map[string]string{labelKey: labelValue},
			},
		},
	}
	t.Run("Test DaemonSet manifest", func(t *testing.T) {
		m, err := nodeExporterDaemonSet(cr, false)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "DaemonSet manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		assert.NotNil(t, m.Spec.Template.Labels)
		assert.Equal(t, labelValue, m.Spec.Template.Labels[labelKey])
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
		assert.NotNil(t, m.Spec.Template.Annotations)
		assert.Equal(t, annotationValue, m.Spec.Template.Annotations[annotationKey])
		assertNodeExporterHardening(t, m, false)
		assertNodeExporterHostAccess(t, m)
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			NodeExporter: &monv1.NodeExporter{},
		},
	}
	t.Run("Test DaemonSet manifest with nil labels and annotation", func(t *testing.T) {
		m, err := nodeExporterDaemonSet(cr, false)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "DaemonSet manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.NotNil(t, m.Spec.Template.Labels)
		assert.Nil(t, m.GetAnnotations())
		assert.Nil(t, m.Spec.Template.Annotations)
	})
	t.Run("Test OpenShift DaemonSet security context", func(t *testing.T) {
		m, err := nodeExporterDaemonSet(cr, true)
		require.NoError(t, err)
		assertNodeExporterHardening(t, m, true)
		assertNodeExporterHostAccess(t, m)
	})
	t.Run("Test configured IDs are preserved", func(t *testing.T) {
		configuredCR := cr.DeepCopy()
		configuredCR.Spec.NodeExporter.SecurityContext = &monv1.SecurityContext{
			RunAsUser:  ptr.To(int64(3000)),
			RunAsGroup: ptr.To(int64(3001)),
			FSGroup:    ptr.To(int64(3002)),
		}

		m, err := nodeExporterDaemonSet(configuredCR, false)
		require.NoError(t, err)
		assert.Equal(t, int64(3000), *m.Spec.Template.Spec.SecurityContext.RunAsUser)
		assert.Equal(t, int64(3001), *m.Spec.Template.Spec.SecurityContext.RunAsGroup)
		assert.Equal(t, int64(3002), *m.Spec.Template.Spec.SecurityContext.FSGroup)
		assertNodeExporterHardening(t, m, false)

		openShiftDaemonSet, err := nodeExporterDaemonSet(configuredCR, true)
		require.NoError(t, err)
		assert.Nil(t, openShiftDaemonSet.Spec.Template.Spec.SecurityContext.RunAsUser)
		assert.Nil(t, openShiftDaemonSet.Spec.Template.Spec.SecurityContext.RunAsGroup)
		assert.Nil(t, openShiftDaemonSet.Spec.Template.Spec.SecurityContext.FSGroup)
	})
	t.Run("Test existing temporary volume and mount are replaced", func(t *testing.T) {
		daemonSet := &appsv1.DaemonSet{
			Spec: appsv1.DaemonSetSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						Volumes: []corev1.Volume{{
							Name: "tmp",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						}},
						Containers: []corev1.Container{{
							Name:         "node-exporter",
							VolumeMounts: []corev1.VolumeMount{{Name: "other-tmp", MountPath: "/tmp"}},
						}},
					},
				},
			},
		}

		applyNodeExporterHardening(daemonSet, false, nil)

		assertNodeExporterHardening(t, daemonSet, false)
		assert.NotContains(t, daemonSet.Spec.Template.Spec.Containers[0].VolumeMounts,
			corev1.VolumeMount{Name: "other-tmp", MountPath: "/tmp"})
	})
	t.Run("Test ClusterRole manifest", func(t *testing.T) {
		m, err := nodeExporterClusterRole(cr, true, false)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ClusterRole manifest should not be empty")
	})
	t.Run("Test ClusterRoleBinding manifest", func(t *testing.T) {
		m, err := nodeExporterClusterRoleBinding(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ClusterRoleBinding manifest should not be empty")
	})
	t.Run("Test Service manifest", func(t *testing.T) {
		m, err := nodeExporterService(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Service manifest should not be empty")
	})
	t.Run("Test ServiceMonitor manifest", func(t *testing.T) {
		crWithLabels := &monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring", Labels: map[string]string{labelKey: labelValue}},
			Spec: monv1.PlatformMonitoringSpec{
				NodeExporter: &monv1.NodeExporter{
					Annotations: map[string]string{annotationKey: annotationValue},
					Labels:      map[string]string{labelKey: labelValue},
				},
			},
		}
		m, err := nodeExporterServiceMonitor(crWithLabels)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ServiceMonitor manifest should not be empty")
		labelsassert.AssertCRLabels(t, m.GetLabels(), utils.NodeExporterComponentName, "victoriametrics-operator", map[string]string{labelKey: labelValue})
	})
	t.Run("Test ServiceAccount manifest", func(t *testing.T) {
		crWithSALabels := &monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
			Spec: monv1.PlatformMonitoringSpec{
				NodeExporter: &monv1.NodeExporter{
					ServiceAccount: &monv1.EmbeddedObjectMetadata{
						Labels: map[string]string{labelKey: labelValue},
					},
				},
			},
		}
		m, err := nodeExporterServiceAccount(crWithSALabels)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ServiceAccount manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey], "ServiceAccount.Labels should be merged")
	})
}

func assertNodeExporterHardening(t *testing.T, daemonSet *appsv1.DaemonSet, isOpenShift bool) {
	t.Helper()

	podSecurityContext := daemonSet.Spec.Template.Spec.SecurityContext
	require.NotNil(t, podSecurityContext)
	assert.Equal(t, ptr.To(true), podSecurityContext.RunAsNonRoot)
	require.NotNil(t, podSecurityContext.SeccompProfile)
	assert.Equal(t, corev1.SeccompProfileTypeRuntimeDefault, podSecurityContext.SeccompProfile.Type)
	if isOpenShift {
		assert.Nil(t, podSecurityContext.RunAsUser)
		assert.Nil(t, podSecurityContext.RunAsGroup)
		assert.Nil(t, podSecurityContext.FSGroup)
	} else {
		require.NotNil(t, podSecurityContext.RunAsUser)
		require.NotNil(t, podSecurityContext.RunAsGroup)
		require.NotNil(t, podSecurityContext.FSGroup)
	}

	require.Len(t, daemonSet.Spec.Template.Spec.Containers, 1)
	container := daemonSet.Spec.Template.Spec.Containers[0]
	require.NotNil(t, container.SecurityContext)
	assert.Equal(t, ptr.To(false), container.SecurityContext.AllowPrivilegeEscalation)
	assert.Equal(t, ptr.To(true), container.SecurityContext.ReadOnlyRootFilesystem)
	require.NotNil(t, container.SecurityContext.Capabilities)
	assert.Equal(t, []corev1.Capability{"ALL"}, container.SecurityContext.Capabilities.Drop)
	assert.Contains(t, container.VolumeMounts, utils.TmpVolumeMount())

	tmpVolumes := 0
	for _, volume := range daemonSet.Spec.Template.Spec.Volumes {
		if volume.Name == "tmp" {
			tmpVolumes++
			require.NotNil(t, volume.EmptyDir)
			require.NotNil(t, volume.EmptyDir.SizeLimit)
			assert.Equal(t, "100Mi", volume.EmptyDir.SizeLimit.String())
		}
	}
	assert.Equal(t, 1, tmpVolumes)
}

func assertNodeExporterHostAccess(t *testing.T, daemonSet *appsv1.DaemonSet) {
	t.Helper()

	podSpec := daemonSet.Spec.Template.Spec
	assert.True(t, podSpec.HostNetwork)
	assert.True(t, podSpec.HostPID)
	assert.False(t, podSpec.HostIPC)

	hostPathVolumes := 0
	for _, volume := range podSpec.Volumes {
		if volume.HostPath != nil {
			hostPathVolumes++
		}
	}
	assert.Equal(t, 5, hostPathVolumes)

	require.Len(t, podSpec.Containers, 1)
	require.NotEmpty(t, podSpec.Containers[0].VolumeMounts)
	require.NotNil(t, podSpec.Containers[0].VolumeMounts[0].MountPropagation)
	assert.Equal(t, corev1.MountPropagationHostToContainer,
		*podSpec.Containers[0].VolumeMounts[0].MountPropagation)
}

func TestNodeExporterSecurityContextConstraintsManifest(t *testing.T) {
	m, err := nodeExporterSecurityContextConstraints()
	require.NoError(t, err)

	assert.Equal(t, utils.NodeExporterComponentName, m.GetName())
	assert.True(t, m.AllowHostNetwork)
	assert.True(t, m.AllowHostPID)
	assert.False(t, m.AllowHostIPC)
	assert.True(t, m.AllowHostDirVolumePlugin)
	assert.False(t, m.AllowPrivilegedContainer)
	assert.True(t, m.ReadOnlyRootFilesystem)
	assert.Equal(t, []corev1.Capability{"ALL"}, m.RequiredDropCapabilities)
}

func newNodeExporterTestReconciler(objects ...runtime.Object) *NodeExporterReconciler {
	testScheme := scheme.Scheme
	_ = secv1.AddToScheme(testScheme)
	kubeClient := fake.NewClientBuilder().WithScheme(testScheme).WithRuntimeObjects(objects...).Build()
	return &NodeExporterReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: kubeClient,
			Scheme: testScheme,
			Log:    utils.Logger("nodeexporter_test"),
		},
	}
}

func TestHandleSecurityContextConstraintsCreatesWhenMissing(t *testing.T) {
	reconciler := newNodeExporterTestReconciler()
	testCR := &monv1.PlatformMonitoring{ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"}}

	require.NoError(t, reconciler.handleSecurityContextConstraints(testCR))

	scc := &secv1.SecurityContextConstraints{}
	require.NoError(t, reconciler.Client.Get(t.Context(),
		client.ObjectKey{Name: utils.NodeExporterComponentName}, scc))
	assert.True(t, scc.AllowHostNetwork)
}

func TestHandleSecurityContextConstraintsUpdatesExisting(t *testing.T) {
	existing := &secv1.SecurityContextConstraints{
		ObjectMeta:               metav1.ObjectMeta{Name: utils.NodeExporterComponentName},
		AllowHostDirVolumePlugin: false,
		Volumes:                  []secv1.FSType{"*"},
	}
	reconciler := newNodeExporterTestReconciler(existing)
	testCR := &monv1.PlatformMonitoring{ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"}}

	require.NoError(t, reconciler.handleSecurityContextConstraints(testCR))

	scc := &secv1.SecurityContextConstraints{}
	require.NoError(t, reconciler.Client.Get(t.Context(),
		client.ObjectKey{Name: utils.NodeExporterComponentName}, scc))
	assert.True(t, scc.AllowHostDirVolumePlugin, "the policy fields must be brought in line with the desired manifest")
	assert.NotEqual(t, []secv1.FSType{"*"}, scc.Volumes)
}

func TestDeleteSecurityContextConstraintsRemovesExisting(t *testing.T) {
	existing := &secv1.SecurityContextConstraints{ObjectMeta: metav1.ObjectMeta{Name: utils.NodeExporterComponentName}}
	reconciler := newNodeExporterTestReconciler(existing)
	testCR := &monv1.PlatformMonitoring{ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"}}

	require.NoError(t, reconciler.deleteSecurityContextConstraints(testCR))

	err := reconciler.Client.Get(t.Context(), client.ObjectKey{Name: utils.NodeExporterComponentName}, &secv1.SecurityContextConstraints{})
	assert.True(t, errors.IsNotFound(err))
}

func TestDeleteSecurityContextConstraintsNoopWhenMissing(t *testing.T) {
	reconciler := newNodeExporterTestReconciler()
	testCR := &monv1.PlatformMonitoring{ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"}}

	assert.NoError(t, reconciler.deleteSecurityContextConstraints(testCR))
}
