package controllers

import (
	"os"
	"testing"
	"time"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	secv1 "github.com/openshift/api/security/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	fakediscovery "k8s.io/client-go/discovery/fake"
	ktesting "k8s.io/client-go/testing"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func TestSuccessfulReconcileResultZeroDoesNotHotLoop(t *testing.T) {
	t.Setenv("RECONCILIATION_INTERVAL", "0")
	result, err := successfulReconcileResult()
	if err != nil {
		t.Fatal(err)
	}
	if result.Requeue || result.RequeueAfter != 0 {
		t.Fatalf("interval 0 must return empty Result, got %#v", result)
	}
}

func TestSuccessfulReconcileResultPositiveIsRepairRequeue(t *testing.T) {
	t.Setenv("RECONCILIATION_INTERVAL", "3600")
	result, err := successfulReconcileResult()
	if err != nil {
		t.Fatal(err)
	}
	if result.RequeueAfter != time.Hour {
		t.Fatalf("got RequeueAfter=%s, want 1h", result.RequeueAfter)
	}
}

func TestSuccessfulReconcileResultDefaultOneHour(t *testing.T) {
	t.Setenv("RECONCILIATION_INTERVAL", "")
	result, err := successfulReconcileResult()
	if err != nil {
		t.Fatal(err)
	}
	if result.RequeueAfter != time.Hour {
		t.Fatalf("got RequeueAfter=%s, want 1h", result.RequeueAfter)
	}
}

func TestJaegerServiceWatchRequiresPrivilegedRights(t *testing.T) {
	prev := utils.PrivilegedRights
	t.Cleanup(func() { utils.PrivilegedRights = prev })
	utils.PrivilegedRights = false
	if jaegerServiceWatchAllowed() {
		t.Fatal("namespaced/non-privileged installs must not start the Jaeger Service informer")
	}
	utils.PrivilegedRights = true
	if !jaegerServiceWatchAllowed() {
		t.Fatal("privileged installs may register the mapped Jaeger Service watch")
	}
}

func TestPlatformMonitoringPredicate(t *testing.T) {
	p := platformMonitoringPredicate()
	oldObj := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "pm",
			Namespace:       "monitoring",
			Generation:      1,
			ResourceVersion: "1",
			Labels:          map[string]string{"a": "1"},
			Annotations:     map[string]string{"k": "v"},
		},
	}

	t.Run("status-only skipped", func(t *testing.T) {
		newObj := oldObj.DeepCopy()
		newObj.ResourceVersion = "2"
		newObj.Status.Conditions = []monv1.PlatformMonitoringCondition{{Type: "Successful"}}
		if p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj}) {
			t.Fatal("status-only update must not enqueue")
		}
	})
	t.Run("generation change enqueued", func(t *testing.T) {
		newObj := oldObj.DeepCopy()
		newObj.Generation = 2
		if !p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj}) {
			t.Fatal("generation change must enqueue")
		}
	})
	t.Run("label change enqueued", func(t *testing.T) {
		newObj := oldObj.DeepCopy()
		newObj.Labels = map[string]string{"a": "2"}
		if !p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj}) {
			t.Fatal("label change must enqueue")
		}
	})
	t.Run("annotation change enqueued", func(t *testing.T) {
		newObj := oldObj.DeepCopy()
		newObj.Annotations = map[string]string{"k": "other"}
		if !p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: newObj}) {
			t.Fatal("annotation change must enqueue")
		}
	})
}

func TestOwnedChildPredicateEnqueuesGenerationlessBodyChange(t *testing.T) {
	p := ownedChildPredicate()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "creds", Namespace: "monitoring", ResourceVersion: "1"},
		Data:       map[string][]byte{"password": []byte("a")},
	}
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "grafana", Namespace: "monitoring", ResourceVersion: "1"},
		Spec:       corev1.ServiceSpec{Ports: []corev1.ServicePort{{Port: 80}}},
	}
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{Name: "grafana", Namespace: "monitoring", ResourceVersion: "1"},
		Rules:      []rbacv1.PolicyRule{{Verbs: []string{"get"}, Resources: []string{"secrets"}}},
	}
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{Name: "grafana", Namespace: "monitoring", ResourceVersion: "1"},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{corev1.ResourceStorage: resource.MustParse("1Gi")},
			},
		},
	}
	endpoints := &corev1.Endpoints{ //nolint:staticcheck // SA1019: v1 Endpoints is deprecated but still served.
		ObjectMeta: metav1.ObjectMeta{Name: "vm-kubelet", Namespace: "monitoring", ResourceVersion: "1"},
		Subsets: []corev1.EndpointSubset{{
			Addresses: []corev1.EndpointAddress{{IP: "10.0.0.1"}},
			Ports:     []corev1.EndpointPort{{Port: 10250}},
		}},
	}

	secretData := secret.DeepCopy()
	secretData.ResourceVersion = "2"
	secretData.Data["password"] = []byte("b")

	secretFields := secret.DeepCopy()
	secretFields.ResourceVersion = "2"
	secretFields.ManagedFields = []metav1.ManagedFieldsEntry{{Manager: "kubectl"}}

	servicePorts := service.DeepCopy()
	servicePorts.ResourceVersion = "2"
	servicePorts.Spec.Ports[0].Port = 3000

	serviceStatus := service.DeepCopy()
	serviceStatus.ResourceVersion = "2"
	serviceStatus.Status.LoadBalancer.Ingress = []corev1.LoadBalancerIngress{{IP: "203.0.113.8"}}

	roleRules := role.DeepCopy()
	roleRules.ResourceVersion = "2"
	roleRules.Rules[0].Verbs = []string{"get", "list"}

	pvcSpec := pvc.DeepCopy()
	pvcSpec.ResourceVersion = "2"
	pvcSpec.Spec.Resources.Requests[corev1.ResourceStorage] = resource.MustParse("2Gi")

	pvcStatus := pvc.DeepCopy()
	pvcStatus.ResourceVersion = "2"
	pvcStatus.Status.Phase = corev1.ClaimBound

	endpointSubsets := endpoints.DeepCopy()
	endpointSubsets.ResourceVersion = "2"
	endpointSubsets.Subsets[0].Addresses[0].IP = "10.0.0.2"

	tests := []struct {
		name    string
		old     client.Object
		updated client.Object
		want    bool
	}{
		{name: "secret data change", old: secret, updated: secretData, want: true},
		{name: "secret managed fields only", old: secret, updated: secretFields, want: false},
		{name: "service port change", old: service, updated: servicePorts, want: true},
		{name: "service load balancer status", old: service, updated: serviceStatus, want: false},
		{name: "role rules change", old: role, updated: roleRules, want: true},
		{name: "pvc storage request change", old: pvc, updated: pvcSpec, want: true},
		{name: "pvc phase change", old: pvc, updated: pvcStatus, want: false},
		{name: "endpoints address change", old: endpoints, updated: endpointSubsets, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.Update(event.UpdateEvent{ObjectOld: tt.old, ObjectNew: tt.updated})
			if got != tt.want {
				t.Fatalf("ownedChildPredicate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOwnedChildPredicateIgnoresStatus(t *testing.T) {
	p := ownedChildPredicate()
	oldObj := &grafv1.GrafanaDashboard{
		ObjectMeta: metav1.ObjectMeta{Name: "dash", Namespace: "monitoring", Generation: 1, ResourceVersion: "1", Labels: map[string]string{"app": "x"}},
	}
	statusOnly := oldObj.DeepCopy()
	statusOnly.ResourceVersion = "9"
	statusOnly.Status.NoMatchingInstances = true
	if p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: statusOnly}) {
		t.Fatal("GrafanaDashboard status-only update must not enqueue")
	}

	specChange := oldObj.DeepCopy()
	specChange.Generation = 2
	if !p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: specChange}) {
		t.Fatal("spec/generation change must enqueue")
	}

	labelChange := oldObj.DeepCopy()
	labelChange.Labels = map[string]string{"app": "y"}
	if !p.Update(event.UpdateEvent{ObjectOld: oldObj, ObjectNew: labelChange}) {
		t.Fatal("label change must enqueue")
	}
}

func TestJaegerServicePredicate(t *testing.T) {
	p := jaegerServicePredicate()
	matching := &corev1.Service{ObjectMeta: metav1.ObjectMeta{
		Name: "query", Namespace: "tracing", Labels: map[string]string{
			"app":                         "jaeger",
			"app.kubernetes.io/component": "query",
			"app.kubernetes.io/part-of":   "jaeger",
		},
	}}
	other := &corev1.Service{ObjectMeta: metav1.ObjectMeta{
		Name: "web", Namespace: "tracing", Labels: map[string]string{"app": "nginx"},
	}}
	if !p.Create(event.CreateEvent{Object: matching}) {
		t.Fatal("matching Jaeger Service create must enqueue")
	}
	if p.Create(event.CreateEvent{Object: other}) {
		t.Fatal("unrelated Service must not enqueue")
	}
}

func TestMapJaegerServicesToPlatformMonitorings(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := monv1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}

	interested := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Name: "pm", Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Integration: &monv1.Integration{
				Jaeger: &monv1.Jaeger{CreateGrafanaDataSource: true},
			},
		},
	}
	ignored := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Name: "other", Namespace: "monitoring"},
	}
	svc := &corev1.Service{ObjectMeta: metav1.ObjectMeta{
		Name: "query", Namespace: "tracing",
		Labels: map[string]string{
			"app":                         "jaeger",
			"app.kubernetes.io/component": "query",
			"app.kubernetes.io/part-of":   "jaeger",
		},
	}}

	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(interested, ignored, svc).Build()
	r := &PlatformMonitoringReconciler{Client: c, Log: log.Log}
	reqs := r.mapJaegerServicesToPlatformMonitorings(t.Context(), svc)
	if len(reqs) != 1 || reqs[0].Name != "pm" || reqs[0].Namespace != "monitoring" {
		t.Fatalf("got %#v, want the Jaeger-enabled PlatformMonitoring", reqs)
	}
}

func TestHasKind(t *testing.T) {
	dc := &fakediscovery.FakeDiscovery{Fake: &ktesting.Fake{}}
	dc.Resources = []*metav1.APIResourceList{
		{
			GroupVersion: grafv1.SchemeGroupVersion.String(),
			APIResources: []metav1.APIResource{{Kind: "Grafana"}, {Kind: "GrafanaDashboard"}},
		},
	}
	r := &PlatformMonitoringReconciler{DiscoveryClient: dc, Log: log.Log}
	if !r.hasKind(grafv1.SchemeGroupVersion.String(), "Grafana") {
		t.Fatal("expected Grafana GVK to be served")
	}
	if r.hasKind("route.openshift.io/v1", "Route") {
		t.Fatal("Route must be absent when discovery does not list it")
	}
}

func TestWatchBasedReconcileEnv(t *testing.T) {
	t.Setenv(watchBasedReconcileEnv, "")
	os.Unsetenv(watchBasedReconcileEnv)
	if !watchBasedReconcileEnabled() {
		t.Fatal("default must enable owned-child watches")
	}
	t.Setenv(watchBasedReconcileEnv, "false")
	if watchBasedReconcileEnabled() {
		t.Fatal("false must disable owned-child watches")
	}
	t.Setenv(watchBasedReconcileEnv, "0")
	if watchBasedReconcileEnabled() {
		t.Fatal("0 must disable owned-child watches")
	}
	t.Setenv(watchBasedReconcileEnv, "true")
	if !watchBasedReconcileEnabled() {
		t.Fatal("true must enable owned-child watches")
	}
}

func TestOwnedChildPredicateDeploymentDelete(t *testing.T) {
	p := ownedChildPredicate()
	d := &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: "grafana-operator", Namespace: "monitoring"}}
	if !p.Delete(event.DeleteEvent{Object: d}) {
		t.Fatal("confirmed Deployment delete must enqueue")
	}
	if p.Delete(event.DeleteEvent{Object: d, DeleteStateUnknown: true}) {
		t.Fatal("unknown delete must not enqueue")
	}
}

func TestClusterScopedWatchRequiresPrivilegedRights(t *testing.T) {
	prev := utils.PrivilegedRights
	t.Cleanup(func() { utils.PrivilegedRights = prev })
	utils.PrivilegedRights = false
	if clusterScopedWatchAllowed() {
		t.Fatal("namespaced installs must not start ClusterRole/SCC informers")
	}
	utils.PrivilegedRights = true
	if !clusterScopedWatchAllowed() {
		t.Fatal("privileged installs may register cluster-scoped watches")
	}
}

func TestMapClusterScopedToPlatformMonitoring(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := monv1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	pm := &monv1.PlatformMonitoring{ObjectMeta: metav1.ObjectMeta{Name: "pm", Namespace: "monitoring"}}
	other := &monv1.PlatformMonitoring{ObjectMeta: metav1.ObjectMeta{Name: "other", Namespace: "other"}}
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(pm, other).Build()
	r := &PlatformMonitoringReconciler{Client: c, Log: log.Log}

	role := &rbacv1.ClusterRole{ObjectMeta: metav1.ObjectMeta{
		Name:   "monitoring-node-exporter",
		Labels: map[string]string{utils.InstallationNamespaceLabelKey: "monitoring"},
	}}
	reqs := r.mapClusterScopedToPlatformMonitoring(t.Context(), role)
	if len(reqs) != 1 || reqs[0].Name != "pm" || reqs[0].Namespace != "monitoring" {
		t.Fatalf("got %#v, want the PlatformMonitoring in the installation namespace", reqs)
	}

	unlabeled := &rbacv1.ClusterRole{ObjectMeta: metav1.ObjectMeta{Name: "system:monitoring"}}
	if reqs := r.mapClusterScopedToPlatformMonitoring(t.Context(), unlabeled); len(reqs) != 0 {
		t.Fatalf("unlabeled ClusterRole must not enqueue, got %#v", reqs)
	}
}

func TestMapSecurityContextConstraintsToPlatformMonitorings(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := monv1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	pm := &monv1.PlatformMonitoring{ObjectMeta: metav1.ObjectMeta{Name: "pm", Namespace: "monitoring"}}
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(pm).Build()
	r := &PlatformMonitoringReconciler{Client: c, Log: log.Log}

	t.Run("installation-namespace label", func(t *testing.T) {
		scc := &secv1.SecurityContextConstraints{ObjectMeta: metav1.ObjectMeta{
			Name:   "monitoring-node-exporter",
			Labels: map[string]string{utils.InstallationNamespaceLabelKey: "monitoring"},
		}}
		reqs := r.mapSecurityContextConstraintsToPlatformMonitorings(t.Context(), scc)
		if len(reqs) != 1 || reqs[0].Name != "pm" {
			t.Fatalf("got %#v", reqs)
		}
	})
	t.Run("node-exporter owner annotations", func(t *testing.T) {
		scc := &secv1.SecurityContextConstraints{ObjectMeta: metav1.ObjectMeta{
			Name: "monitoring-node-exporter",
			Annotations: map[string]string{
				nodeExporterSCCOwnerNameAnnotation:      "pm",
				nodeExporterSCCOwnerNamespaceAnnotation: "monitoring",
			},
		}}
		reqs := r.mapSecurityContextConstraintsToPlatformMonitorings(t.Context(), scc)
		if len(reqs) != 1 || reqs[0].Namespace != "monitoring" || reqs[0].Name != "pm" {
			t.Fatalf("got %#v", reqs)
		}
	})
	t.Run("fixed VM operator name", func(t *testing.T) {
		scc := &secv1.SecurityContextConstraints{ObjectMeta: metav1.ObjectMeta{Name: utils.VmOperatorComponentName}}
		reqs := r.mapSecurityContextConstraintsToPlatformMonitorings(t.Context(), scc)
		if len(reqs) != 1 || reqs[0].Name != "pm" {
			t.Fatalf("got %#v", reqs)
		}
	})
	t.Run("unrelated SCC", func(t *testing.T) {
		scc := &secv1.SecurityContextConstraints{ObjectMeta: metav1.ObjectMeta{Name: "privileged"}}
		if reqs := r.mapSecurityContextConstraintsToPlatformMonitorings(t.Context(), scc); len(reqs) != 0 {
			t.Fatalf("unrelated SCC must not enqueue, got %#v", reqs)
		}
	})
}

func TestOwnIfServedSkipsMissingGVK(t *testing.T) {
	dc := &fakediscovery.FakeDiscovery{Fake: &ktesting.Fake{}}
	dc.Resources = []*metav1.APIResourceList{
		{
			GroupVersion: grafv1.SchemeGroupVersion.String(),
			APIResources: []metav1.APIResource{{Kind: "Grafana"}},
		},
	}
	r := &PlatformMonitoringReconciler{DiscoveryClient: dc, Log: log.Log}
	if r.hasKind(grafv1.SchemeGroupVersion.String(), "GrafanaDatasource") {
		t.Fatal("GrafanaDatasource must be reported absent")
	}
	if !r.hasKind(grafv1.SchemeGroupVersion.String(), "Grafana") {
		t.Fatal("Grafana must be reported served")
	}
}

var _ client.Object = &monv1.PlatformMonitoring{}
