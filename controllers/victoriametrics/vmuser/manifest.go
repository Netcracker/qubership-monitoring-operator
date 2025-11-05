package vmuser

import (
	"embed"
	"maps"

	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/victoriametrics"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

//go:embed  assets/*.yaml
var assets embed.FS

func vmUser(cr *v1alpha1.PlatformMonitoring) (*vmetricsv1b1.VMUser, error) {
	vmuser := vmetricsv1b1.VMUser{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmUserAsset), 100).Decode(&vmuser); err != nil {
		return nil, err
	}

	vmuser.SetNamespace(cr.GetNamespace())
	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmUser.IsInstall() && cr.Spec.Victoriametrics.VmAuth.IsInstall() {

		// Set labels
		vmuser.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(vmuser.GetName(), vmuser.GetNamespace())
		vmuser.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmUser.Image)

		if cr.Spec.Victoriametrics.VmUser.Name != nil {
			vmuser.Spec.Name = cr.Spec.Victoriametrics.VmUser.Name
		}
		if cr.Spec.Victoriametrics.VmUser.UserName != nil {
			vmuser.Spec.UserName = cr.Spec.Victoriametrics.VmUser.UserName
		}
		if cr.Spec.Victoriametrics.VmUser.PasswordRef != nil {
			vmuser.Spec.PasswordRef = cr.Spec.Victoriametrics.VmUser.PasswordRef
		}
		if cr.Spec.Victoriametrics.VmUser.TokenRef != nil {
			vmuser.Spec.TokenRef = cr.Spec.Victoriametrics.VmUser.TokenRef
		}
		if cr.Spec.Victoriametrics.VmUser.GeneratePassword {
			vmuser.Spec.GeneratePassword = true
		}
		if cr.Spec.Victoriametrics.VmUser.BearerToken != nil {
			vmuser.Spec.BearerToken = cr.Spec.Victoriametrics.VmUser.BearerToken
		}
		if cr.Spec.Victoriametrics.VmUser.TargetRefs != nil {
			vmuser.Spec.TargetRefs = cr.Spec.Victoriametrics.VmUser.TargetRefs
		} else {
			var vmSPaths []string

			if cr.Spec.Victoriametrics.VmSingle.IsInstall() {
				targetRef := vmetricsv1b1.TargetRef{
					CRD: &vmetricsv1b1.CRDRef{
						Kind:      "VMSingle",
						Name:      "k8s",
						Namespace: cr.GetNamespace(),
					},
					Paths: []string{""},
				}
				vmuser.Spec.TargetRefs = append(vmuser.Spec.TargetRefs, targetRef)
			} else {
				if cr.Spec.Victoriametrics.VmAgent.IsInstall() {
					targetRef := vmetricsv1b1.TargetRef{
						CRD: &vmetricsv1b1.CRDRef{
							Kind:      "VMAgent",
							Name:      "k8s",
							Namespace: cr.GetNamespace(),
						},
						Paths: []string{""},
					}
					vmuser.Spec.TargetRefs = append(vmuser.Spec.TargetRefs, targetRef)
				}
			}

			if cr.Spec.Victoriametrics.VmCluster.IsInstall() {
				if cr.Spec.Victoriametrics.VmCluster.VmSelect != nil {
					targetRef := vmetricsv1b1.TargetRef{
						CRD: &vmetricsv1b1.CRDRef{
							Kind:      "VMCluster/vmselect",
							Name:      "k8s",
							Namespace: cr.GetNamespace(),
						},
						Paths: vmSelectPaths(),
					}
					vmuser.Spec.TargetRefs = append(vmuser.Spec.TargetRefs, targetRef)
				}
				if cr.Spec.Victoriametrics.VmCluster.VmStorage != nil {
					targetRef := vmetricsv1b1.TargetRef{
						CRD: &vmetricsv1b1.CRDRef{
							Kind:      "VMCluster/vmstorage",
							Name:      "k8s",
							Namespace: cr.GetNamespace(),
						},
						Paths: []string{""},
					}
					vmuser.Spec.TargetRefs = append(vmuser.Spec.TargetRefs, targetRef)
				}
				if cr.Spec.Victoriametrics.VmCluster.VmInsert != nil {
					targetRef := vmetricsv1b1.TargetRef{
						CRD: &vmetricsv1b1.CRDRef{
							Kind:      "VMCluster/vminsert",
							Name:      "k8s",
							Namespace: cr.GetNamespace(),
						},
						Paths: vmInsertPaths(),
					}
					vmuser.Spec.TargetRefs = append(vmuser.Spec.TargetRefs, targetRef)
				}
			}

			if cr.Spec.Victoriametrics.VmAlert.IsInstall() {
				vmAlert := vmetricsv1b1.VMAlert{}
				vmAlert.SetName(utils.VmComponentName)
				vmAlert.SetNamespace(cr.GetNamespace())
				if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
					vmAlert.Spec.ExtraArgs = make(map[string]string)
					maps.Copy(vmAlert.Spec.ExtraArgs, map[string]string{"tls": "true"})
				}
				if cr.Spec.Victoriametrics.VmAlert.Port != "" {
					vmAlert.Spec.Port = cr.Spec.Victoriametrics.VmAlert.Port
				} else {
					vmAlert.Spec.Port = "8080"
				}
				targetRef := vmetricsv1b1.TargetRef{
					CRD: &vmetricsv1b1.CRDRef{
						Kind:      "VMAlert",
						Name:      "k8s",
						Namespace: cr.GetNamespace(),
					},
					Static: &vmetricsv1b1.StaticRef{
						URL: vmAlert.AsURL(),
					},
					Paths: vmAlertPaths(),
				}
				vmSPaths = vmSinglePaths()
				vmuser.Spec.TargetRefs = append(vmuser.Spec.TargetRefs, targetRef)
			}

			if cr.Spec.Victoriametrics.VmAlertManager.IsInstall() {
				targetRef := vmetricsv1b1.TargetRef{
					CRD: &vmetricsv1b1.CRDRef{
						Kind:      "VMAlertmanager",
						Name:      "k8s",
						Namespace: cr.GetNamespace(),
					},
					Paths: vmAlertManagerPaths(),
				}
				vmuser.Spec.TargetRefs = append(vmuser.Spec.TargetRefs, targetRef)
			}

			if cr.Spec.Victoriametrics.VmAgent.IsInstall() {
				targetRef := vmetricsv1b1.TargetRef{
					CRD: &vmetricsv1b1.CRDRef{
						Kind:      "VMAgent",
						Name:      "k8s",
						Namespace: cr.GetNamespace(),
					},
					Paths: vmAgentPaths(),
				}
				vmuser.Spec.TargetRefs = append(vmuser.Spec.TargetRefs, targetRef)
			}

			if cr.Spec.Victoriametrics.VmSingle.IsInstall() {
				if len(vmSPaths) == 0 {
					vmSPaths = vmSingleFullPaths()
				}
				targetRef := vmetricsv1b1.TargetRef{
					CRD: &vmetricsv1b1.CRDRef{
						Kind:      "VMSingle",
						Name:      "k8s",
						Namespace: cr.GetNamespace(),
					},
					Paths: vmSPaths,
				}
				vmuser.Spec.TargetRefs = append(vmuser.Spec.TargetRefs, targetRef)
			}
		}
	}

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
		vmuser.Spec.UserConfigOption = vmetricsv1b1.UserConfigOption{
			TLSConfig: &vmetricsv1b1.TLSConfig{
				InsecureSkipVerify: false,
				CAFile:             "/etc/vm/secrets/" + victoriametrics.GetVmauthTLSSecretName(cr.Spec.Victoriametrics.VmAuth) + "/ca.crt",
				CertFile:           "/etc/vm/secrets/" + victoriametrics.GetVmauthTLSSecretName(cr.Spec.Victoriametrics.VmAuth) + "/tls.crt",
				KeyFile:            "/etc/vm/secrets/" + victoriametrics.GetVmauthTLSSecretName(cr.Spec.Victoriametrics.VmAuth) + "/tls.key",
			},
		}
	}

	return &vmuser, nil
}

func vmSelectPaths() []string {
	return []string{"/select/.*"}
}

func vmInsertPaths() []string {
	return []string{"/insert/.*"}
}

func vmAlertPaths() []string {
	return []string{
		"/api/v1/rules",
		"/api/v1/alerts",
		"/api/v1/alert.*",
		"/vmalert.*",
	}
}

func vmAlertManagerPaths() []string {
	return []string{
		"/api/v2/alerts.*",
		"/api/v2/receivers.*",
		"/api/v2/silences.*",
		"/api/v2/status.*",
	}
}

func vmAgentPaths() []string {
	return []string{
		"/config.*",
		"/target.*",
		"/service-discovery.*",
		"/static.*",
		"/api/v1/write",
		"/api/v1/import.*",
		"/api/v1/target.*",
	}
}

func vmSingleFullPaths() []string {
	return []string{"/vmui.*",
		"/config.*",
		"/graph.*",
		"/api/v1/label.*",
		"/api/v1/query.*",
		"/api/v1/rules",
		"/api/v1/alerts",
		"/api/v1/metadata.*",
		"/api/v1/format.*",
		"/api/v1/series.*",
		"/api/v1/status.*",
		"/api/v1/export.*",
		"/api/v1/admin/tsdb.*",
		"/prometheus/graph.*",
		"/prometheus/api/v1/label.*",
		"/graphite.*",
		"/prometheus/api/v1/query.*",
		"/prometheus/api/v1/rules",
		"/prometheus/api/v1/alerts",
		"/prometheus/api/v1/metadata",
		"/prometheus/api/v1/series.*",
		"/prometheus/api/v1/status.*",
		"/prometheus/api/v1/export.*",
		"/prometheus/federate",
		"/prometheus/api/v1/admin/tsdb.*",
	}
}

func vmSinglePaths() []string {
	return []string{"/vmui.*",
		"/graph.*",
		"/api/v1/label.*",
		"/api/v1/query.*",
		"/api/v1/metadata.*",
		"/api/v1/format.*",
		"/api/v1/series.*",
		"/api/v1/status.*",
		"/api/v1/export.*",
		"/api/v1/admin/tsdb.*",
		"/prometheus/graph.*",
		"/prometheus/api/v1/label.*",
		"/graphite.*",
		"/prometheus/api/v1/query.*",
		"/prometheus/api/v1/rules",
		"/prometheus/api/v1/alerts",
		"/prometheus/api/v1/metadata",
		"/prometheus/api/v1/series.*",
		"/prometheus/api/v1/status.*",
		"/prometheus/api/v1/export.*",
		"/prometheus/federate",
		"/prometheus/api/v1/admin/tsdb.*",
	}
}
