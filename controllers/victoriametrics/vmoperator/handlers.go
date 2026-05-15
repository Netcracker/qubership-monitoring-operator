package vmoperator

import (
	"context"
	"fmt"
	"reflect"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"

	"github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	secv1 "github.com/openshift/api/security/v1"
	errs "github.com/pkg/errors"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *VmOperatorReconciler) handleRole(cr *monv1.PlatformMonitoring) error {
	m, err := vmOperatorRole(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Role manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmOperator.Image)

	e := &rbacv1.Role{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.SetLabels(m.GetLabels())
	e.SetName(m.GetName())
	e.Rules = m.Rules

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmOperatorReconciler) handleServiceAccount(cr *monv1.PlatformMonitoring) error {
	m, err := vmOperatorServiceAccount(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceAccount manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmOperator.Image)

	e := &corev1.ServiceAccount{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.SetLabels(m.GetLabels())

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmOperatorReconciler) handleRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := vmOperatorRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating RoleBinding manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmOperator.Image)

	e := &rbacv1.RoleBinding{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.SetLabels(m.GetLabels())

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmOperatorReconciler) handleClusterRole(cr *monv1.PlatformMonitoring) error {
	m, err := vmOperatorClusterRole(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRole manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmOperator.Image)

	e := &rbacv1.ClusterRole{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.SetLabels(m.GetLabels())
	e.SetName(m.GetName())
	e.Rules = m.Rules

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmOperatorReconciler) handleClusterRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := vmOperatorClusterRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRoleBinding manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmOperator.Image)

	e := &rbacv1.ClusterRoleBinding{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.SetLabels(m.GetLabels())

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmOperatorReconciler) handleDeployment(cr *monv1.PlatformMonitoring) error {
	m, err := vmOperatorDeployment(r, cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Deployment manifest")
		return err
	}
	e := &appsv1.Deployment{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.SetLabels(m.GetLabels())
	e.Spec.Selector = m.Spec.Selector
	e.Spec.Template.SetLabels(m.Spec.Template.GetLabels())
	e.Spec.Template.Spec.SecurityContext = m.Spec.Template.Spec.SecurityContext
	e.Spec.Template.Spec.Affinity = m.Spec.Template.Spec.Affinity
	e.Spec.Template.Spec.Volumes = m.Spec.Template.Spec.Volumes
	e.Spec.Template.Spec.Containers = m.Spec.Template.Spec.Containers
	e.Spec.Template.Spec.ServiceAccountName = m.Spec.Template.Spec.ServiceAccountName
	e.Spec.Template.Spec.PriorityClassName = m.Spec.Template.Spec.PriorityClassName
	e.Spec.Replicas = m.Spec.Replicas

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmOperatorReconciler) handleService(cr *monv1.PlatformMonitoring) error {
	m, err := vmOperatorService(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Service manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmOperator.Image)

	e := &corev1.Service{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.SetLabels(m.GetLabels())
	e.Spec.Ports = m.Spec.Ports
	e.Spec.Selector = m.Spec.Selector

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmOperatorReconciler) handleKubeletService(cr *monv1.PlatformMonitoring) error {
	m, err := vmKubeletService(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Service manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmOperator.Image)

	e := &corev1.Service{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.SetLabels(m.GetLabels())
	e.Spec.Ports = m.Spec.Ports
	e.Spec.Selector = m.Spec.Selector

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmOperatorReconciler) handleKubeletServiceEndpoints(cr *monv1.PlatformMonitoring) error {
	nodes, err := r.KubeClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		r.Log.Error(err, "Failed to retrieve nodes to get addresses")
		return errs.Wrap(err, "Failed to list nodes to get addresses")
	}
	if r.HasApi(discoveryv1.SchemeGroupVersion, "EndpointSlice") {
		eps, err := vmKubeletServiceEndpointSlice(cr)
		if err != nil {
			r.Log.Error(err, "Failed creating EndpointSlice manifest")
			return err
		}
		return r.upsertNodeEndpointSlice(cr, eps, nodes)
	}
	eps, err := vmKubeletServiceEndpoints(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Endpoints manifest")
		return err
	}
	return r.upsertNodeEndpoints(cr, eps, nodes)
}

func (r *VmOperatorReconciler) handleKubeSchedulerService(cr *monv1.PlatformMonitoring) error {
	m, err := vmKubeSchedulerService(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Service manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmOperator.Image)

	e := &corev1.Service{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.SetLabels(m.GetLabels())
	e.Spec.Ports = m.Spec.Ports
	e.Spec.Selector = m.Spec.Selector

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmOperatorReconciler) handleKubeSchedulerServiceEndpoints(cr *monv1.PlatformMonitoring) error {
	allNodes, err := r.KubeClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		r.Log.Error(err, "Failed to retrieve nodes to get addresses")
		return errs.Wrap(err, "Failed to list nodes to get addresses")
	}
	cpNodes := filterNodes(allNodes, isControlPlaneNode)
	if r.HasApi(discoveryv1.SchemeGroupVersion, "EndpointSlice") {
		eps, err := vmKubeSchedulerServiceEndpointSlice(cr)
		if err != nil {
			r.Log.Error(err, "Failed creating EndpointSlice manifest")
			return err
		}
		return r.upsertNodeEndpointSlice(cr, eps, cpNodes)
	}
	eps, err := vmKubeSchedulerServiceEndpoints(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Endpoints manifest")
		return err
	}
	return r.upsertNodeEndpoints(cr, eps, cpNodes)
}
func (r *VmOperatorReconciler) handleKubeControllerManagerService(cr *monv1.PlatformMonitoring) error {
	m, err := vmKubeControllerManagerService(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Service manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmOperator.Image)

	e := &corev1.Service{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.SetLabels(m.GetLabels())
	e.Spec.Ports = m.Spec.Ports
	e.Spec.Selector = m.Spec.Selector

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmOperatorReconciler) handleKubeControllerManagerServiceEndpoints(cr *monv1.PlatformMonitoring) error {
	allNodes, err := r.KubeClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		r.Log.Error(err, "Failed to retrieve nodes to get addresses")
		return errs.Wrap(err, "Failed to list nodes to get addresses")
	}
	cpNodes := filterNodes(allNodes, isControlPlaneNode)
	if r.HasApi(discoveryv1.SchemeGroupVersion, "EndpointSlice") {
		eps, err := vmKubeControllerManagerServiceEndpointSlice(cr)
		if err != nil {
			r.Log.Error(err, "Failed creating EndpointSlice manifest")
			return err
		}
		return r.upsertNodeEndpointSlice(cr, eps, cpNodes)
	}
	eps, err := vmKubeControllerManagerServiceEndpoints(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Endpoints manifest")
		return err
	}
	return r.upsertNodeEndpoints(cr, eps, cpNodes)
}

func filterNodes(nodes *corev1.NodeList, predicate func(corev1.Node) bool) *corev1.NodeList {
	out := &corev1.NodeList{}
	for _, n := range nodes.Items {
		if predicate(n) {
			out.Items = append(out.Items, n)
		}
	}
	return out
}
func isControlPlaneNode(n corev1.Node) bool {
	if _, ok := n.Labels["node-role.kubernetes.io/control-plane"]; ok {
		return true
	}
	if _, ok := n.Labels["node-role.kubernetes.io/master"]; ok {
		return true
	}
	for _, t := range n.Spec.Taints {
		if t.Key == "node-role.kubernetes.io/control-plane" ||
			t.Key == "node-role.kubernetes.io/master" {
			return true
		}
	}
	return false
}

func nodeAddress(node corev1.Node) (string, map[corev1.NodeAddressType][]string, error) {
	m := map[corev1.NodeAddressType][]string{}
	for _, a := range node.Status.Addresses {
		m[a.Type] = append(m[a.Type], a.Address)
	}

	if addresses, ok := m[corev1.NodeInternalIP]; ok {
		return addresses[0], m, nil
	}
	if addresses, ok := m[corev1.NodeExternalIP]; ok {
		return addresses[0], m, nil
	}
	return "", m, fmt.Errorf("host address unknown")
}

// upsertNodeEndpoints creates or updates a legacy v1 Endpoints object populated from the given nodes.
// Kept for clusters that don't serve discovery.k8s.io/v1 EndpointSlice; new clusters take the upsertNodeEndpointSlice path.
func (r *VmOperatorReconciler) upsertNodeEndpoints(cr *monv1.PlatformMonitoring, eps *corev1.Endpoints, nodes *corev1.NodeList) error { //nolint:staticcheck // SA1019: see manifest.go.
	addresses, ers := getNodeAddresses(nodes)
	for _, err := range ers {
		r.Log.Error(err, "")
	}
	eps.Subsets[0].Addresses = addresses

	eps.Labels["name"] = utils.TruncLabel(eps.GetName())
	eps.Labels["app.kubernetes.io/name"] = utils.TruncLabel(eps.GetName())
	eps.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(eps.GetName(), eps.GetNamespace())
	eps.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmOperator.Image)

	e := &corev1.Endpoints{ObjectMeta: eps.ObjectMeta} //nolint:staticcheck // SA1019: see above.
	if err := r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return r.CreateResource(cr, eps)
		}
		return err
	}
	e.SetLabels(eps.GetLabels())
	e.Subsets = eps.Subsets
	return r.UpdateResource(e)
}

func (r *VmOperatorReconciler) upsertNodeEndpointSlice(cr *monv1.PlatformMonitoring, eps *discoveryv1.EndpointSlice, nodes *corev1.NodeList) error {
	endpointEntries, ers := getNodeEndpoints(nodes)
	for _, err := range ers {
		r.Log.Error(err, "")
	}
	eps.Endpoints = endpointEntries

	eps.Labels["name"] = utils.TruncLabel(eps.GetName())
	eps.Labels["app.kubernetes.io/name"] = utils.TruncLabel(eps.GetName())
	eps.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(eps.GetName(), eps.GetNamespace())
	eps.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmOperator.Image)

	e := &discoveryv1.EndpointSlice{ObjectMeta: eps.ObjectMeta}
	if err := r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return r.CreateResource(cr, eps)
		}
		return err
	}
	e.SetLabels(eps.GetLabels())
	e.AddressType = eps.AddressType
	e.Ports = eps.Ports
	e.Endpoints = eps.Endpoints
	return r.UpdateResource(e)
}

func getNodeAddresses(nodes *corev1.NodeList) ([]corev1.EndpointAddress, []error) { //nolint:staticcheck // SA1019: legacy code path kept for clusters without discovery.k8s.io/v1 EndpointSlice support.
	addresses := make([]corev1.EndpointAddress, 0) //nolint:staticcheck // SA1019: see above.
	ers := make([]error, 0)

	for _, n := range nodes.Items {
		address, _, err := nodeAddress(n)
		if err != nil {
			ers = append(ers, errs.Wrapf(err, "failed to determine hostname for node (%s)", n.Name))
			continue
		}
		addresses = append(addresses, corev1.EndpointAddress{ //nolint:staticcheck // SA1019: see above.
			IP: address,
			TargetRef: &corev1.ObjectReference{
				Kind:       "Node",
				Name:       n.Name,
				UID:        n.UID,
				APIVersion: n.APIVersion,
			},
		})
	}

	return addresses, ers
}

func getNodeEndpoints(nodes *corev1.NodeList) ([]discoveryv1.Endpoint, []error) {
	endpoints := make([]discoveryv1.Endpoint, 0)
	ers := make([]error, 0)

	ready := true
	for _, n := range nodes.Items {
		address, _, err := nodeAddress(n)
		if err != nil {
			ers = append(ers, errs.Wrapf(err, "failed to determine hostname for node (%s)", n.Name))
			continue
		}
		endpoints = append(endpoints, discoveryv1.Endpoint{
			Addresses: []string{address},
			Conditions: discoveryv1.EndpointConditions{
				Ready: &ready,
			},
			TargetRef: &corev1.ObjectReference{
				Kind:       "Node",
				Name:       n.Name,
				UID:        n.UID,
				APIVersion: n.APIVersion,
			},
		})
	}

	return endpoints, ers
}

func (r *VmOperatorReconciler) handleServiceMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := vmOperatorServiceMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceMonitor manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmOperator.Image)

	e := &promv1.ServiceMonitor{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.SetLabels(m.GetLabels())
	e.Spec.JobLabel = m.Spec.JobLabel
	e.Spec.Endpoints = m.Spec.Endpoints
	e.Spec.NamespaceSelector = m.Spec.NamespaceSelector

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmOperatorReconciler) handleSecurityContextConstraints(cr *monv1.PlatformMonitoring) error {
	m, err := vmOperatorSecurityContextConstraints()
	if err != nil {
		r.Log.Error(err, "Failed creating SecurityContextConstraints manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmOperator.Image)

	e := &secv1.SecurityContextConstraints{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	//Set parameters
	e.SetLabels(m.GetLabels())

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmOperatorReconciler) deleteServiceAccount(cr *monv1.PlatformMonitoring) error {
	m, err := vmOperatorServiceAccount(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceAccount manifest")
		return err
	}
	e := &corev1.ServiceAccount{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err = r.DeleteResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmOperatorReconciler) deleteRole(cr *monv1.PlatformMonitoring) error {
	m, err := vmOperatorRole(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Role manifest")
		return err
	}
	e := &rbacv1.Role{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err = r.DeleteResource(e); err != nil {
		return err
	}

	return nil
}

func (r *VmOperatorReconciler) deleteRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := vmOperatorRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating RoleBinding manifest")
		return err
	}
	e := &rbacv1.RoleBinding{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err = r.DeleteResource(e); err != nil {
		return err
	}

	return nil
}

func (r *VmOperatorReconciler) deleteClusterRole(cr *monv1.PlatformMonitoring) error {
	m, err := vmOperatorClusterRole(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRole manifest")
		return err
	}
	e := &rbacv1.ClusterRole{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err = r.DeleteResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmOperatorReconciler) deleteClusterRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := vmOperatorClusterRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRoleBinding manifest")
		return err
	}
	e := &rbacv1.ClusterRoleBinding{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err = r.DeleteResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmOperatorReconciler) deleteVmOperatorDeployment(cr *monv1.PlatformMonitoring) error {
	m, err := vmOperatorDeployment(r, cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Deployment manifest")
		return err
	}
	e := &appsv1.Deployment{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err = r.DeleteResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmOperatorReconciler) deleteService(cr *monv1.PlatformMonitoring) error {
	m, err := vmOperatorService(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Service manifest")
		return err
	}
	e := &corev1.Service{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err = r.DeleteResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmOperatorReconciler) deleteServiceMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := vmOperatorServiceMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceMonitor manifest")
		return err
	}
	e := &promv1.ServiceMonitor{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err = r.DeleteResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmOperatorReconciler) deleteVmOperatorConfigMap(cr *monv1.PlatformMonitoring) error {
	e := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "57410f0d.victoriametrics.com",
		Namespace: cr.GetNamespace()}}

	if err := r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err := r.DeleteResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmOperatorReconciler) deleteAllCRDObjects(cr *monv1.PlatformMonitoring) error {

	objectList := []client.ObjectList{
		&v1beta1.VMAgentList{},
		&v1beta1.VMAlertList{},
		&v1beta1.VMAlertmanagerList{},
		&v1beta1.VMAlertmanagerConfigList{},
		&v1beta1.VMAuthList{},
		&v1beta1.VMClusterList{},
		&v1beta1.VMNodeScrapeList{},
		&v1beta1.VMPodScrapeList{},
		&v1beta1.VMProbeList{},
		&v1beta1.VMRuleList{},
		&v1beta1.VMServiceScrapeList{},
		&v1beta1.VMSingleList{},
		&v1beta1.VMStaticScrapeList{},
		&v1beta1.VMUserList{},
	}

	var foundObjectList []client.ObjectList
	for _, object := range objectList {
		if err := r.Client.List(context.Background(), object, client.InNamespace(cr.GetNamespace()), client.MatchingLabels{"app.kubernetes.io/component": "monitoring"}); err != nil {
			return err
		}
		if reflect.ValueOf(object).Elem().FieldByName("Items").Len() > 0 {
			foundObjectList = append(foundObjectList, object)
		}
	}

	r.Log.Info("CRD object list for deleting", "length - ", len(foundObjectList))

	for _, object := range foundObjectList {
		objType := getListObjectType(object)
		proto := reflect.New(objType).Interface()

		if err := r.Client.DeleteAllOf(context.Background(), proto.(client.Object), client.InNamespace(cr.GetNamespace()), client.MatchingLabels{"app.kubernetes.io/component": "monitoring"}); err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		}
	}

	return nil
}

func getListObjectType(list client.ObjectList) reflect.Type {
	objType := reflect.ValueOf(list).Elem().FieldByName("Items").Type().Elem()
	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}

	return objType
}
