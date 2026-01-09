/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	_ "net/http/pprof"
	"os"

	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/cache"

	qubershiporg1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	apis "github.com/grafana-operator/grafana-operator/v4/api"
	grafv1 "github.com/grafana-operator/grafana-operator/v4/api/integreatly/v1alpha1"
	secv1 "github.com/openshift/api/security/v1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	extensionsobj "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/discovery"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = utils.Logger("setup")
	logger   = utils.Logger("cmd")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(qubershiporg1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr, probeAddr, pprofAddr string
	var enableLeaderElection, pprofEnabled bool

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&utils.PrivilegedRights, "privilegedRights", false, "Indicates is extended privileges should be used for the monitoring components")
	flag.BoolVar(&pprofEnabled, "pprof-enable", false, "Enable pprof.")
	flag.StringVar(&pprofAddr, "pprof-address", ":9180", "The pprof address.")
	flag.Parse()

	ctrl.SetLogger(utils.Logger(""))

	namespace, found := os.LookupEnv("WATCH_NAMESPACE")
	if !found {
		namespace = "monitoring"
	}

	if !pprofEnabled {
		pprofAddr = ""
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		Metrics:                metricsserver.Options{BindAddress: metricsAddr},
		HealthProbeBindAddress: probeAddr,
		PprofBindAddress:       pprofAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "b0cb59fe.netcracker.com",
		NewCache: func(config *rest.Config, opts cache.Options) (cache.Cache, error) {
			opts.DefaultNamespaces = map[string]cache.Config{namespace: {}}
			return cache.New(config, opts)
		},
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}
	logger.Info("Registering Components.")

	// Setup Scheme for all resources
	if err = apis.AddToScheme(mgr.GetScheme()); err != nil {
		logger.Error(err, "")
		os.Exit(1)
	}

	// Add extention scheme
	if err = extensionsobj.AddToScheme(mgr.GetScheme()); err != nil {
		logger.Error(err, "")
		os.Exit(1)
	}

	// Add prometheus scheme
	if err = promv1.AddToScheme(mgr.GetScheme()); err != nil {
		logger.Error(err, "")
		os.Exit(1)
	}

	// Add victoriametrics scheme
	if err = vmetricsv1b1.AddToScheme(mgr.GetScheme()); err != nil {
		logger.Error(err, "")
		os.Exit(1)
	}

	// Add scc scheme
	if err = secv1.Install(mgr.GetScheme()); err != nil {
		logger.Error(err, "")
		os.Exit(1)
	}

	// Add grafana schema
	if err = grafv1.AddToScheme(mgr.GetScheme()); err != nil {
		logger.Error(err, "")
		os.Exit(1)
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(mgr.GetConfig())
	if err != nil {
		logger.Error(err, "Get discoveryClient failed")
	}
	if err = (&controllers.PlatformMonitoringReconciler{
		Client:          mgr.GetClient(),
		Scheme:          mgr.GetScheme(),
		Log:             utils.Logger("controller-platformmonitoring"),
		Config:          mgr.GetConfig(),
		DiscoveryClient: discoveryClient,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "PlatformMonitoring")
		os.Exit(1)
	}

	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err = mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
