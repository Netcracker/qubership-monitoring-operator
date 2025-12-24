package utils

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	sprig "github.com/go-task/slim-sprig"

	routev1 "github.com/openshift/api/route/v1"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// MustAssetReader loads and return the asset for the given name as bytes reader.
// Panics when the asset loading would return an error.
func MustAssetReader(assets embed.FS, asset string) io.Reader {
	content, _ := assets.ReadFile(asset)
	return bytes.NewReader(content)
}

// MustAssetReaderToString loads and return the asset for the given name as a string.
// Panics when the asset loading would return an error.
func MustAssetReaderToString(assets embed.FS, asset string) string {
	content, _ := assets.ReadFile(asset)
	return string(content)
}

// HasIngressV1beta1Api checks that the cluster API has networking.k8s.io.v1beta1.Ingress API.
// It helps to identify whether the cluster is a Kubernetes cluster.
// This API unavailable in Kubernetes v1.22+.
func (r *ComponentReconciler) HasIngressV1beta1Api() bool {
	return r.HasApi(v1beta1.SchemeGroupVersion, "Ingress")
}

// HasIngressV1Api checks that the cluster API has networking.k8s.io.v1.Ingress API.
// It helps to identify whether the cluster is a Kubernetes cluster.
// This API available from Kubernetes v1.19+.
func (r *ComponentReconciler) HasIngressV1Api() bool {
	return r.HasApi(networkingv1.SchemeGroupVersion, "Ingress")
}

// HasRouteApi checks that the cluster API has v1.Route API.
// It helps to identify whether the cluster is an Openshift cluster.
func (r *ComponentReconciler) HasRouteApi() bool {
	return r.HasApi(routev1.GroupVersion, "Route")
}

// HasApi checks that cluster API has specified API.
// Return true if API exists.
func (r *ComponentReconciler) HasApi(groupVersion schema.GroupVersion, kind string) bool {
	hasApi, err := ResourceExists(r.Dc, groupVersion.String(), kind)
	if err != nil {
		r.Log.Error(err, "Error while check hasAPI")
	}
	return hasApi
}

// ResourceExists returns true if the given resource kind exists
// in the given api groupversion
func ResourceExists(dc discovery.DiscoveryInterface, apiGroupVersion, kind string) (bool, error) {
	_, apiLists, err := dc.ServerGroupsAndResources()
	if err != nil {
		return false, err
	}
	for _, apiList := range apiLists {
		if apiList.GroupVersion == apiGroupVersion {
			for _, r := range apiList.APIResources {
				if r.Kind == kind {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func (r *ComponentReconciler) CreateResource(cr *v1alpha1.PlatformMonitoring, o K8sResource, setRefOptional ...bool) error {
	res := o.GetObjectKind().GroupVersionKind().Kind
	setRef := true
	if len(setRefOptional) > 0 {
		setRef = setRefOptional[0]
	}
	if setRef {
		if err := controllerutil.SetControllerReference(cr, o, r.Scheme); err != nil {
			if !(strings.Contains(err.Error(), "cluster-scoped resource must not have a namespace-scoped owner") ||
				strings.Contains(err.Error(), "cross-namespace owner references are disallowed")) {
				return err
			}
		}
	}
	if err := r.Client.Create(context.TODO(), o); err != nil {
		return err
	}
	r.Log.Info("Successful creating", ResourceKey, res)
	return nil
}

func (r *ComponentReconciler) GetResource(o K8sResource) error {
	objectKey := client.ObjectKeyFromObject(o)
	if err := r.Client.Get(context.TODO(), objectKey, o); err != nil {
		return err
	}
	return nil
}

func (r *ComponentReconciler) UpdateResource(o K8sResource) error {
	// Update object
	if err := r.Client.Update(context.TODO(), o); err != nil {
		return err
	}
	r.Log.Info("Successful updating", ResourceKey, o.GetObjectKind().GroupVersionKind().Kind)
	return nil
}

func (r *ComponentReconciler) DeleteResource(o K8sResource) error {
	err := r.Client.Delete(context.TODO(), o)
	if err != nil {
		return err
	}
	r.Log.Info("Successful deleting", ResourceKey, o.GetObjectKind().GroupVersionKind().Kind)
	return nil
}

// GetServerVersion allows to recognize OpenShift v4.5 or higher.
func (r *ComponentReconciler) IsOpenShiftV4() (bool, error) {
	isOpenShift := r.HasRouteApi()
	minorServerVersion, err := r.GetMinorServerVersion()
	if err != nil {
		r.Log.Error(err, "Failed to recognize OpenShift V4")
		return false, err
	}
	// We support OpenShift since v4.5 that contains Kubernetes v1.18
	if !isOpenShift || minorServerVersion < 18 {
		return false, nil
	}
	return true, nil
}

// GetServerVersion allows to recognize OpenShift 3.11
func (r *ComponentReconciler) IsOpenShiftV3() (bool, error) {
	isOpenShift := r.HasRouteApi()
	minorServerVersion, err := r.GetMinorServerVersion()
	if err != nil {
		r.Log.Error(err, "Failed to recognize OpenShift V3.11")
		return false, err
	}
	if isOpenShift && minorServerVersion == 11 {
		return true, nil
	}
	return false, nil
}

// GetServerVersion tries to get info about Kubernetes server.
func (r *ComponentReconciler) GetServerInfo() (*version.Info, error) {
	serverVersion, err := r.Dc.ServerVersion()
	if err != nil {
		r.Log.Error(err, "Failed to get server version of Kubernetes")
		return nil, err
	}
	return serverVersion, nil
}

// GetServerVersion allows to recognise minor version of Kubernetes server.
// It can be useful because different K8s versions have a different API.
func (r *ComponentReconciler) GetMinorServerVersion() (int, error) {
	serverInfo, err := r.GetServerInfo()
	if err != nil {
		return 0, err
	}
	// Some version of OpenShift (e.g. v3.11) have a plus in the minor Kubernetes version: v1.11+
	minorServerVersion, err := strconv.Atoi(strings.Trim(serverInfo.Minor, "+"))
	if err != nil {
		r.Log.Error(err, fmt.Sprintf("Failed to convert minor server version %s to integer", serverInfo.Minor))
		return 0, err
	}
	return minorServerVersion, nil
}

func GetTimeNow() string {
	return time.Now().Format(time.RFC3339)
}

func GetFromResourceMap(resourceList core.ResourceList, key string) string {
	quantity := resourceList[core.ResourceName(key)]
	return quantity.String()
}

// ParseTemplate allows parsing of assets as Go templates.
// Left and Right delims override default {{ and }} delimiters.
// Delims overriding can be helpful for assets where {{ and }} are already present, e.g. in Grafana dashboards.
func ParseTemplate(fileContent, filePath, leftDelim, rightDelim string, parameters interface{}) (string, error) {
	funcMap := sprig.TxtFuncMap()
	funcMap["resIndex"] = GetFromResourceMap
	funcMap["timeNow"] = GetTimeNow

	goTemplate, err := template.New(filePath).Delims(leftDelim, rightDelim).Funcs(funcMap).Parse(fileContent)
	if err != nil {
		err := fmt.Errorf("the template for file %s cannot be parsed. Error: %s", filePath, err.Error())
		return "", err
	}

	writer := strings.Builder{}

	if err := goTemplate.Execute(&writer, parameters); err != nil {
		err := fmt.Errorf("the template for file %s cannot be executed. Error: %s", filePath, err.Error())
		return "", err
	}

	return writer.String(), nil
}

func GetTagFromImage(image string) string {
	partsOfImage := strings.Split(image, ":")
	return partsOfImage[len(partsOfImage)-1]
}

func GetInstanceLabel(name, namespace string) string {
	label := fmt.Sprintf("%s-%s", name, namespace)
	return TruncLabel(label)
}

func TruncLabel(label string) string {
	if len(label) >= 63 {
		return strings.Trim(label[:63], "-")
	}
	return strings.Trim(label, "-")
}

const defaultPodsPendingTimeout = time.Minute * 5

func (r *ComponentReconciler) WaitForPodsReadiness(o K8sResource) error {
	start := time.Now()
	for {
		if time.Since(start) >= defaultPodsPendingTimeout {
			return errors.New(fmt.Sprintf("Timeout %s waiting for pods readiness", defaultPodsPendingTimeout.Round(time.Second).String()))
		}
		if err := r.GetResource(o); err != nil {
			return err
		}

		switch o.GetObjectKind().GroupVersionKind() {
		case appsv1.SchemeGroupVersion.WithKind("Deployment"):
			depl, ok := o.(*appsv1.Deployment)
			if !ok {
				return errors.New("Could not convert resource to deployment")
			}
			if depl.Status.Replicas == *depl.Spec.Replicas && depl.Status.ReadyReplicas == *depl.Spec.Replicas && depl.Status.AvailableReplicas == *depl.Spec.Replicas && depl.Status.UpdatedReplicas == *depl.Spec.Replicas {
				return nil
			}
		case appsv1.SchemeGroupVersion.WithKind("StatefulSet"):
			ss, ok := o.(*appsv1.StatefulSet)
			if !ok {
				return errors.New("Could not convert resource to statefulset")
			}
			if ss.Status.Replicas == *ss.Spec.Replicas && ss.Status.ReadyReplicas == *ss.Spec.Replicas {
				return nil
			}
		default:
			return errors.New("Could not get status of pods")
		}
		time.Sleep(time.Second * 5)
	}
}