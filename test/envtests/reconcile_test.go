package bdd_tests

import (
	"context"
	"embed"

	v1beta1 "github.com/Netcracker/qubership-monitoring-operator/api/v1beta1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/alertmanager"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/grafana"
	grafana_operator "github.com/Netcracker/qubership-monitoring-operator/controllers/grafana-operator"
	kubernetes_monitors "github.com/Netcracker/qubership-monitoring-operator/controllers/kubernetes-monitors"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/kubestatemetrics"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/nodeexporter"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/prometheus"
	prometheus_operator "github.com/Netcracker/qubership-monitoring-operator/controllers/prometheus-operator"
	prometheus_rules "github.com/Netcracker/qubership-monitoring-operator/controllers/prometheus-rules"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	grafv1 "github.com/grafana-operator/grafana-operator/v4/api/integreatly/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

//go:embed  assets/*.yaml
var assets embed.FS

var _ = Describe("Reconcile", func() {
	It("Create PlatformMonitoring CR", func() {
		//Create test CR
		cr = v1beta1.PlatformMonitoring{}
		err = yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, "assets/platformmonitoring.yaml"), 100).Decode(&cr)
		cr.SetUID(uuid.NewUUID())
		Expect(err).NotTo(HaveOccurred())
		Expect(cr).NotTo(BeNil())
	})
	It("Prometheus-operator", func() {
		// Run Prometheus-operator
		poReconciler := prometheus_operator.NewPrometheusOperatorReconciler(k8sClient, scheme.Scheme)
		err = poReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		// Get Role
		role := rbacv1.Role{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.PrometheusOperatorComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: role.Namespace,
			Name:      role.Name,
		}, &role)
		Expect(err).NotTo(HaveOccurred())
		Expect(role).NotTo(BeNil())
		logf.Log.Info("Getting Role successful")

		// Get ServiceAccount
		sa := corev1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.PrometheusOperatorComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: sa.Namespace,
			Name:      sa.Name,
		}, &sa)
		Expect(err).NotTo(HaveOccurred())
		Expect(sa).NotTo(BeNil())
		logf.Log.Info("Getting ServiceAccount successful")

		// Get RoleBinding
		roleBinding := rbacv1.RoleBinding{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.PrometheusOperatorComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: roleBinding.Namespace,
			Name:      roleBinding.Name,
		}, &roleBinding)
		Expect(err).NotTo(HaveOccurred())
		Expect(roleBinding).NotTo(BeNil())
		logf.Log.Info("Getting RoleBinding successful")

		// Get Service
		service := corev1.Service{ObjectMeta: metav1.ObjectMeta{
			Name:      utils.PrometheusOperatorComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: service.Namespace,
			Name:      service.Name,
		}, &service)
		Expect(err).NotTo(HaveOccurred())
		Expect(service).NotTo(BeNil())
		logf.Log.Info("Getting Service successful")

		// Get Deployment
		deployment := appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{
			Name:      utils.PrometheusOperatorComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: deployment.Namespace,
			Name:      deployment.Name,
		}, &deployment)
		Expect(err).NotTo(HaveOccurred())
		Expect(deployment).NotTo(BeNil())
		Expect(deployment.Spec.Template.Spec.SecurityContext.FSGroup).Should(Equal(cr.Spec.Prometheus.Operator.SecurityContext.FSGroup))
		Expect(deployment.Spec.Template.Spec.SecurityContext.RunAsUser).Should(Equal(cr.Spec.Prometheus.Operator.SecurityContext.RunAsUser))
		logf.Log.Info("Getting Deployment successful")
	})
	It("KubernetesMonitors", func() {
		//Run KubernetesMonitors reconcile
		kubernetesMonitorsReconciler := kubernetes_monitors.NewKubernetesMonitorsReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = kubernetesMonitorsReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		// Get ApiServerServiceMonitor
		sm := promv1.ServiceMonitor{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + "kube-apiserver-service-monitor",
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: sm.Namespace,
			Name:      sm.Name,
		}, &sm)
		Expect(err).NotTo(HaveOccurred())
		Expect(sm).NotTo(BeNil())
		logf.Log.Info("Getting ApiServerServiceMonitor successful")

		// Get ControllerManagerServiceMonitor
		sm = promv1.ServiceMonitor{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + "kube-controller-manager-service-monitor",
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: sm.Namespace,
			Name:      sm.Name,
		}, &sm)
		Expect(err).NotTo(HaveOccurred())
		Expect(sm).NotTo(BeNil())
		logf.Log.Info("Getting ControllerManagerServiceMonitor successful")

		// Get SchedulerServiceMonitor
		sm = promv1.ServiceMonitor{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + "kube-scheduler-service-monitor",
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: sm.Namespace,
			Name:      sm.Name,
		}, &sm)
		Expect(err).NotTo(HaveOccurred())
		Expect(sm).NotTo(BeNil())
		logf.Log.Info("Getting SchedulerServiceMonitor successful")

		// Get KubeletServiceMonitor
		sm = promv1.ServiceMonitor{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + "kubelet-service-monitor",
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: sm.Namespace,
			Name:      sm.Name,
		}, &sm)
		Expect(err).NotTo(HaveOccurred())
		Expect(sm).NotTo(BeNil())
		logf.Log.Info("Getting KubeletServiceMonitor successful")
	})
	It("Prometheus", func() {
		//Run Prometheus reconcile
		pReconciler := prometheus.NewPrometheusReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = pReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		// Get ServiceAccount
		sa := corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cr.GetNamespace() + "-" + utils.PrometheusComponentName,
				Namespace: cr.GetNamespace(),
			},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: sa.Namespace,
			Name:      sa.Name,
		}, &sa)
		Expect(err).NotTo(HaveOccurred())
		Expect(sa).NotTo(BeNil())
		logf.Log.Info("Getting ServiceAccount successful")

		// Get Prometheus
		prom := promv1.Prometheus{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "k8s",
				Namespace: cr.GetNamespace(),
			},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: prom.Namespace,
			Name:      prom.Name,
		}, &prom)
		Expect(err).NotTo(HaveOccurred())
		Expect(prom).NotTo(BeNil())
		Expect(prom.Spec.ServiceAccountName).Should(Equal(cr.GetNamespace() + "-" + utils.PrometheusComponentName))
		Expect(prom.Spec.SecurityContext.FSGroup).Should(Equal(cr.Spec.Prometheus.SecurityContext.FSGroup))
		Expect(prom.Spec.SecurityContext.RunAsUser).Should(Equal(cr.Spec.Prometheus.SecurityContext.RunAsUser))
		logf.Log.Info("Getting Prometheus successful")

		// Get Ingress
		ingress := netv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cr.GetNamespace() + "-" + utils.PrometheusComponentName,
				Namespace: cr.GetNamespace(),
			},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: ingress.Namespace,
			Name:      ingress.Name,
		}, &ingress)
		Expect(err).NotTo(HaveOccurred())
		Expect(ingress).NotTo(BeNil())
		logf.Log.Info("Getting Ingress successful")

		// Get PodMonitor
		podMonitor := promv1.PodMonitor{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cr.GetNamespace() + "-" + "prometheus-pod-monitor",
				Namespace: cr.GetNamespace(),
			},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: podMonitor.Namespace,
			Name:      podMonitor.Name,
		}, &podMonitor)
		Expect(err).NotTo(HaveOccurred())
		Expect(podMonitor).NotTo(BeNil())
		logf.Log.Info("Getting PodMonitor successful")
	})
	It("Alertmanager", func() {
		// Run AlertManager reconcile
		aReconciler := alertmanager.NewAlertManagerReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = aReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		// Get Secret
		secret := corev1.Secret{ObjectMeta: metav1.ObjectMeta{
			Name:      "alertmanager-k8s",
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: secret.Namespace,
			Name:      secret.Name,
		}, &secret)
		Expect(err).NotTo(HaveOccurred())
		Expect(secret).NotTo(BeNil())
		logf.Log.Info("Getting Secret successful")

		// Get ServiceAccount
		sa := corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cr.GetNamespace() + "-" + utils.AlertManagerComponentName,
				Namespace: cr.GetNamespace(),
			},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: sa.Namespace,
			Name:      sa.Name,
		}, &sa)
		Expect(err).NotTo(HaveOccurred())
		Expect(sa).NotTo(BeNil())
		logf.Log.Info("Getting ServiceAccount successful")

		// Get Service
		service := corev1.Service{ObjectMeta: metav1.ObjectMeta{
			Name:      utils.AlertManagerComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: service.Namespace,
			Name:      service.Name,
		}, &service)
		Expect(err).NotTo(HaveOccurred())
		Expect(service).NotTo(BeNil())
		logf.Log.Info("Getting Service successful")

		// Get AlertManager
		alertmanager := promv1.Alertmanager{ObjectMeta: metav1.ObjectMeta{
			Name:      "k8s",
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: alertmanager.Namespace,
			Name:      alertmanager.Name,
		}, &alertmanager)
		Expect(err).NotTo(HaveOccurred())
		Expect(alertmanager).NotTo(BeNil())
		Expect(alertmanager.Spec.SecurityContext.FSGroup).Should(Equal(cr.Spec.AlertManager.SecurityContext.FSGroup))
		Expect(alertmanager.Spec.SecurityContext.RunAsUser).Should(Equal(cr.Spec.AlertManager.SecurityContext.RunAsUser))
		logf.Log.Info("Getting AlertManager successful")

		// Get Ingress
		ingress := netv1.Ingress{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.AlertManagerComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: ingress.Namespace,
			Name:      ingress.Name,
		}, &ingress)
		Expect(err).NotTo(HaveOccurred())
		Expect(ingress).NotTo(BeNil())
		logf.Log.Info("Getting Ingress successful")

		// Get PodMonitor
		podMonitor := promv1.PodMonitor{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + "alertmanager-pod-monitor",
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: podMonitor.Namespace,
			Name:      podMonitor.Name,
		}, &podMonitor)
		Expect(err).NotTo(HaveOccurred())
		Expect(podMonitor).NotTo(BeNil())
		logf.Log.Info("Getting PodMonitor successful")
	})
	It("KubeStateMetrics", func() {
		// Run KubeStateMetrics reconcile
		ksmReconciler := kubestatemetrics.NewKubeStateMetricsReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = ksmReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		// Get ServiceAccount
		sa := corev1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.KubestatemetricsComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: sa.Namespace,
			Name:      sa.Name,
		}, &sa)
		Expect(err).NotTo(HaveOccurred())
		Expect(sa).NotTo(BeNil())
		logf.Log.Info("Getting ServiceAccount successful")

		// Get Service
		service := corev1.Service{ObjectMeta: metav1.ObjectMeta{
			Name:      utils.KubestatemetricsComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: service.Namespace,
			Name:      service.Name,
		}, &service)
		Expect(err).NotTo(HaveOccurred())
		Expect(service).NotTo(BeNil())
		logf.Log.Info("Getting Service successful")

		// Get Deployment
		deployment := appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{
			Name:      utils.KubestatemetricsComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: deployment.Namespace,
			Name:      deployment.Name,
		}, &deployment)
		Expect(err).NotTo(HaveOccurred())
		Expect(deployment).NotTo(BeNil())
		Expect(deployment.Spec.Template.Spec.SecurityContext.FSGroup).Should(Equal(cr.Spec.KubeStateMetrics.SecurityContext.FSGroup))
		Expect(deployment.Spec.Template.Spec.SecurityContext.RunAsUser).Should(Equal(cr.Spec.KubeStateMetrics.SecurityContext.RunAsUser))
		logf.Log.Info("Getting Deployment successful")

		// Get ServiceMonitor
		sm := promv1.ServiceMonitor{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.KubestatemetricsComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: sm.Namespace,
			Name:      sm.Name,
		}, &sm)
		Expect(err).NotTo(HaveOccurred())
		Expect(sm).NotTo(BeNil())
		logf.Log.Info("Getting ServiceMonitor successful")
	})
	It("NodeExporter", func() {
		//Run NodeExporter reconcile
		neReconciler := nodeexporter.NewNodeExporterReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = neReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		// Get ServiceAccount
		sa := corev1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.NodeExporterComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: sa.Namespace,
			Name:      sa.Name,
		}, &sa)
		Expect(err).NotTo(HaveOccurred())
		Expect(sa).NotTo(BeNil())
		logf.Log.Info("Getting ServiceAccount successful")

		// Get Service
		service := corev1.Service{ObjectMeta: metav1.ObjectMeta{
			Name:      utils.NodeExporterComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: service.Namespace,
			Name:      service.Name,
		}, &service)
		Expect(err).NotTo(HaveOccurred())
		Expect(service).NotTo(BeNil())
		logf.Log.Info("Getting Service successful")

		// Get DaemonSet
		daemonSet := appsv1.DaemonSet{ObjectMeta: metav1.ObjectMeta{
			Name:      utils.NodeExporterComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: daemonSet.Namespace,
			Name:      daemonSet.Name,
		}, &daemonSet)
		Expect(err).NotTo(HaveOccurred())
		Expect(daemonSet).NotTo(BeNil())
		Expect(daemonSet.Spec.Template.Spec.SecurityContext.RunAsUser).Should(Equal(cr.Spec.NodeExporter.SecurityContext.RunAsUser))
		Expect(daemonSet.Spec.Template.Spec.SecurityContext.FSGroup).Should(Equal(cr.Spec.NodeExporter.SecurityContext.FSGroup))
		logf.Log.Info("Getting DaemonSet successful")

		// Get ServiceMonitor
		sm := promv1.ServiceMonitor{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.NodeExporterComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: sm.Namespace,
			Name:      sm.Name,
		}, &sm)
		Expect(err).NotTo(HaveOccurred())
		Expect(sm).NotTo(BeNil())
		logf.Log.Info("Getting ServiceMonitor successful")
	})
	It("Grafana-operator", func() {
		// Run Grafana-operator reconcile
		goReconciler := grafana_operator.NewGrafanaOperatorReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = goReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		// Get GrafanaDashboard
		dashboardList := []string{
			"alerts-overview",
			"etcd-dashboard",
			"go-processes",
			"jvm-processes",
			"home-dashboard",
			"kubernetes-cluster-overview",
			"kubernetes-kubelet",
			"kubernetes-distribution-by-labels",
			"kubernetes-namespace-resources",
			"kubernetes-nodes-resources",
			"kubernetes-pod-resources",
			"kubernetes-pods-distribution-by-node",
			"kubernetes-top-resources",
			"prometheus-cardinality-explorer",
			"prometheus-self-monitor",
			"alertmanager-overview",
			"grafana-overview",
			"tls-status",
			"ha-services",
		}
		for _, name := range dashboardList {
			dashboard := grafv1.GrafanaDashboard{ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: cr.GetNamespace(),
			},
			}
			err = k8sClient.Get(context.TODO(), client.ObjectKey{
				Namespace: dashboard.Namespace,
				Name:      dashboard.Name,
			}, &dashboard)
			Expect(err).NotTo(HaveOccurred())
			Expect(dashboard).NotTo(BeNil())
			logf.Log.Info("Getting GrafanaDashboard successful", "name", dashboard.Name)
		}

		// Get ServiceAccount
		sa := corev1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.GrafanaOperatorComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: sa.Namespace,
			Name:      sa.Name,
		}, &sa)
		Expect(err).NotTo(HaveOccurred())
		Expect(sa).NotTo(BeNil())
		logf.Log.Info("Getting ServiceAccount successful")

		// Get Role
		role := rbacv1.Role{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.GrafanaOperatorComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: role.Namespace,
			Name:      role.Name,
		}, &role)
		Expect(err).NotTo(HaveOccurred())
		Expect(role).NotTo(BeNil())
		logf.Log.Info("Getting Role successful")

		// Get RoleBinding
		roleBinding := rbacv1.RoleBinding{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.GrafanaOperatorComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: roleBinding.Namespace,
			Name:      roleBinding.Name,
		}, &roleBinding)
		Expect(err).NotTo(HaveOccurred())
		Expect(roleBinding).NotTo(BeNil())
		logf.Log.Info("Getting RoleBinding successful")

		// Get Deployment
		deployment := appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{
			Name:      utils.GrafanaOperatorComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: deployment.Namespace,
			Name:      deployment.Name,
		}, &deployment)
		Expect(err).NotTo(HaveOccurred())
		Expect(deployment).NotTo(BeNil())
		Expect(deployment.Spec.Template.Spec.SecurityContext.FSGroup).Should(Equal(cr.Spec.Grafana.Operator.SecurityContext.FSGroup))
		Expect(deployment.Spec.Template.Spec.SecurityContext.RunAsUser).Should(Equal(cr.Spec.Grafana.Operator.SecurityContext.RunAsUser))
		logf.Log.Info("Getting Deployment successful")
	})
	It("Grafana", func() {
		// Run Grafana reconcile
		gReconciler := grafana.NewGrafanaReconciler(k8sClient, scheme.Scheme, discoveryClient, cfg)
		err = gReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		// Get Grafana
		grafana := grafv1.Grafana{ObjectMeta: metav1.ObjectMeta{
			Name:      "grafana",
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: grafana.Namespace,
			Name:      grafana.Name,
		}, &grafana)
		Expect(err).NotTo(HaveOccurred())
		Expect(grafana).NotTo(BeNil())
		Expect(grafana.Spec.Deployment.SecurityContext.RunAsUser).Should(Equal(cr.Spec.Grafana.SecurityContext.RunAsUser))
		Expect(grafana.Spec.Deployment.SecurityContext.FSGroup).Should(Equal(cr.Spec.Grafana.SecurityContext.FSGroup))
		logf.Log.Info("Getting Grafana successful")

		// Get GrafanaDatasource
		grafanaDataSource := grafv1.GrafanaDataSource{ObjectMeta: metav1.ObjectMeta{
			Name:      "platform-monitoring-prometheus",
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: grafanaDataSource.Namespace,
			Name:      grafanaDataSource.Name,
		}, &grafanaDataSource)
		Expect(err).NotTo(HaveOccurred())
		Expect(grafanaDataSource).NotTo(BeNil())
		logf.Log.Info("Getting GrafanaDatasource successful")

		// Get Ingress
		ingress := netv1.Ingress{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + utils.GrafanaComponentName,
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: ingress.Namespace,
			Name:      ingress.Name,
		}, &ingress)
		Expect(err).NotTo(HaveOccurred())
		Expect(ingress).NotTo(BeNil())
		logf.Log.Info("Getting Ingress successful")

		// Get PodMonitor
		podMonitor := promv1.PodMonitor{ObjectMeta: metav1.ObjectMeta{
			Name:      cr.GetNamespace() + "-" + "grafana-pod-monitor",
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: podMonitor.Namespace,
			Name:      podMonitor.Name,
		}, &podMonitor)
		Expect(err).NotTo(HaveOccurred())
		Expect(podMonitor).NotTo(BeNil())
		logf.Log.Info("Getting PodMonitor successful")
	})
	It("PrometheusRule", func() {
		// Run PrometheusRule reconcile
		prReconciler := prometheus_rules.NewPrometheusRulesReconciler(k8sClient, scheme.Scheme)
		err = prReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		// Get PrometheusRule
		rule := promv1.PrometheusRule{ObjectMeta: metav1.ObjectMeta{
			Name:      "prometheus-rules",
			Namespace: cr.GetNamespace(),
		},
		}
		err = k8sClient.Get(context.TODO(), client.ObjectKey{
			Namespace: rule.Namespace,
			Name:      rule.Name,
		}, &rule)
		Expect(err).NotTo(HaveOccurred())
		Expect(rule).NotTo(BeNil())
		logf.Log.Info("Getting PrometheusRule successful")
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

		// Run AlertManager reconcile
		aReconciler := alertmanager.NewAlertManagerReconciler(k8sClient, scheme.Scheme, discoveryClient)
		err = aReconciler.Run(&cr)
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

		// Run Grafana reconcile
		gReconciler := grafana.NewGrafanaReconciler(k8sClient, scheme.Scheme, discoveryClient, cfg)
		err = gReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())

		// Run PrometheusRule reconcile
		prReconciler := prometheus_rules.NewPrometheusRulesReconciler(k8sClient, scheme.Scheme)
		err = prReconciler.Run(&cr)
		Expect(err).NotTo(HaveOccurred())
	})
})
