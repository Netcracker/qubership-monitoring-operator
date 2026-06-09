package utils

import (
	"github.com/go-logr/logr"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// K8sResource abstract k8s resource which can be reconciled
type K8sResource interface {
	runtime.Object
	k8smetav1.Object

	GetObjectMeta() k8smetav1.Object
}

type ComponentReconciler struct {
	Client client.Client
	Scheme *runtime.Scheme
	Dc     discovery.DiscoveryInterface
	Log    logr.Logger
}
