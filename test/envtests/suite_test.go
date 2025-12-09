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

package bdd_tests

import (
	"context"
	"path/filepath"
	"testing"

	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	grafv1 "github.com/grafana-operator/grafana-operator/v4/api/integreatly/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	v1beta1ext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var discoveryClient *discovery.DiscoveryClient
var err error
var cr v1alpha1.PlatformMonitoring

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {

	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "charts", "qubership-monitoring-operator", "crds"),
			filepath.Join("..", "..", "charts", "qubership-monitoring-operator", "charts", "grafana-operator", "crds"),
			filepath.Join("..", "..", "charts", "qubership-monitoring-operator", "charts", "prometheus-adapter-operator", "crds"),
			filepath.Join("..", "..", "charts", "qubership-monitoring-operator", "charts", "prometheus-operator", "crds"),
			filepath.Join("..", "..", "charts", "qubership-monitoring-operator", "charts", "victoriametrics-operator", "crds")},
	}
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	discoveryClient, err = discovery.NewDiscoveryClientForConfig(cfg)
	Expect(err).NotTo(HaveOccurred())

	err = v1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = promv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = grafv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = v1beta1ext.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	//Create test namespace
	namespace := corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "monitoring"}}
	err := k8sClient.Create(context.TODO(), &namespace)
	Expect(err).NotTo(HaveOccurred())

	By("successfully bootstrap test environment")
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
