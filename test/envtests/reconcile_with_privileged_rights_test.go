package bdd_tests

import (
	"context"
	"flag"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	grafana_operator "github.com/Netcracker/qubership-monitoring-operator/controllers/grafana-operator"
	kubernetes_monitors "github.com/Netcracker/qubership-monitoring-operator/controllers/kubernetes-monitors"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/kubestatemetrics"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/nodeexporter"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/prometheus"
	prometheus_operator "github.com/Netcracker/qubership-monitoring-operator/controllers/prometheus-operator"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var _ = Describe("Reconcile with privilegedRights", func() {
	It("Create PlatformMonitoring CR", func() {
		//Create test CR
		cr = monv1.PlatformMonitoring{}
		err = yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, "assets/platformmonitoring.yaml"), 100).Decode(&cr)
		cr.SetUID(uuid.NewUUID())
		Expect(err).NotTo(HaveOccurred())
		Expect(cr).NotTo(BeNil())
		flag.BoolVar(&utils.PrivilegedRights, "privilegedRights", true, "Indicates is extended privileges should be used for the monitoring components")
		flag.Parse()
	})
	It("Prometheus-operator", func() {
		// Run Prometheus-operator
		poReconciler := prometheus_operator.NewPrometheusOperatorReconciler(k8sClient, scheme.Scheme)
		err = poReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		// Get ClusterRole
		clusterRole := rbacv1.ClusterRole{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.PrometheusOperatorComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: clusterRole.Namespace,
			Name:      clusterRole.Name,
		}, &clusterRole)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusterRole).NotTo(BeNil())
		logf.Log.Info("Getting ClusterRole successful")

		// Get ClusterRoleBinding
		clusterRoleBinding := rbacv1.ClusterRoleBinding{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.PrometheusOperatorComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: clusterRoleBinding.Namespace,
			Name:      clusterRoleBinding.Name,
		}, &clusterRoleBinding)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusterRoleBinding).NotTo(BeNil())
		logf.Log.Info("Getting ClusterRoleBinding successful")
	})
	It("Prometheus", func() {
		//Run Prometheus reconcile
		pReconciler := prometheus.NewPrometheusReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = pReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		// Get ClusterRole
		clusterRole := rbacv1.ClusterRole{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.PrometheusComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: clusterRole.Namespace,
			Name:      clusterRole.Name,
		}, &clusterRole)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusterRole).NotTo(BeNil())
		logf.Log.Info("Getting ClusterRole successful")

		// Get ClusterRoleBinding
		clusterRoleBinding := rbacv1.ClusterRoleBinding{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.PrometheusComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: clusterRoleBinding.Namespace,
			Name:      clusterRoleBinding.Name,
		}, &clusterRoleBinding)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusterRoleBinding).NotTo(BeNil())
		logf.Log.Info("Getting ClusterRoleBinding successful")
	})
	It("KubeStateMetrics", func() {
		// Run KubeStateMetrics reconcile
		ksmReconciler := kubestatemetrics.NewKubeStateMetricsReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = ksmReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		// Get ClusterRole
		clusterRole := rbacv1.ClusterRole{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.KubestatemetricsComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: clusterRole.Namespace,
			Name:      clusterRole.Name,
		}, &clusterRole)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusterRole).NotTo(BeNil())
		logf.Log.Info("Getting ClusterRole successful")

		// Get ClusterRoleBinding
		clusterRoleBinding := rbacv1.ClusterRoleBinding{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.KubestatemetricsComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: clusterRoleBinding.Namespace,
			Name:      clusterRoleBinding.Name,
		}, &clusterRoleBinding)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusterRoleBinding).NotTo(BeNil())
		logf.Log.Info("Getting ClusterRoleBinding successful")
	})
	It("NodeExporter", func() {
		//Run NodeExporter reconcile
		neReconciler := nodeexporter.NewNodeExporterReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = neReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		// Get ClusterRole
		clusterRole := rbacv1.ClusterRole{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.NodeExporterComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: clusterRole.Namespace,
			Name:      clusterRole.Name,
		}, &clusterRole)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusterRole).NotTo(BeNil())
		logf.Log.Info("Getting ClusterRole successful")

		// Get ClusterRoleBinding
		clusterRoleBinding := rbacv1.ClusterRoleBinding{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.NodeExporterComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: clusterRoleBinding.Namespace,
			Name:      clusterRoleBinding.Name,
		}, &clusterRoleBinding)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusterRoleBinding).NotTo(BeNil())
		logf.Log.Info("Getting ClusterRoleBinding successful")
	})
	It("Grafana-operator", func() {
		// Run Grafana-operator reconcile
		goReconciler := grafana_operator.NewGrafanaOperatorReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = goReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		// Get ClusterRole
		clusterRole := rbacv1.ClusterRole{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.GrafanaOperatorComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: clusterRole.Namespace,
			Name:      clusterRole.Name,
		}, &clusterRole)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusterRole).NotTo(BeNil())
		logf.Log.Info("Getting ClusterRole successful")

		// Get ClusterRoleBinding
		clusterRoleBinding := rbacv1.ClusterRoleBinding{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.GrafanaOperatorComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: clusterRoleBinding.Namespace,
			Name:      clusterRoleBinding.Name,
		}, &clusterRoleBinding)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusterRoleBinding).NotTo(BeNil())
		logf.Log.Info("Getting ClusterRoleBinding successful")
	})
	It("Cleanup", func() {
		flag := false
		cr.Spec.KubernetesMonitors = nil
		cr.Spec.Prometheus.Install = &flag
		cr.Spec.AlertManager.Install = &flag
		cr.Spec.KubeStateMetrics.Install = &flag
		cr.Spec.NodeExporter.Install = &flag
		cr.Spec.GrafanaDashboards.Install = &flag
		cr.Spec.Grafana.Install = &flag
		cr.Spec.PrometheusRules.Install = &flag

		//Run KubernetesMonitors reconcile
		kubernetesMonitorsReconciler := kubernetes_monitors.NewKubernetesMonitorsReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = kubernetesMonitorsReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		//Run Prometheus reconcile
		pReconciler := prometheus.NewPrometheusReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = pReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		// Run KubeStateMetrics reconcile
		ksmReconciler := kubestatemetrics.NewKubeStateMetricsReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = ksmReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		//Run NodeExporter reconcile
		neReconciler := nodeexporter.NewNodeExporterReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = neReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		// Run Grafana-operator reconcile
		goReconciler := grafana_operator.NewGrafanaOperatorReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = goReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())
	})
})
