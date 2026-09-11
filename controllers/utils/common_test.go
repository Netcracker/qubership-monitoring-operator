package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	fakediscovery "k8s.io/client-go/discovery/fake"
	clienttesting "k8s.io/client-go/testing"
)

func TestResourceExistsIgnoresUnrelatedDiscoveryFailures(t *testing.T) {
	discoveryClient := &fakediscovery.FakeDiscovery{Fake: &clienttesting.Fake{
		Resources: []*metav1.APIResourceList{{
			GroupVersion: "security.openshift.io/v1",
			APIResources: []metav1.APIResource{{Kind: "SecurityContextConstraints"}},
		}},
	}}
	discoveryClient.PrependReactor("get", "group", func(clienttesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("unrelated API group discovery failed")
	})

	exists, err := ResourceExists(
		discoveryClient,
		"security.openshift.io/v1",
		"SecurityContextConstraints",
	)

	require.NoError(t, err)
	assert.True(t, exists)
}

func TestResourceExistsReturnsFalseWhenGroupVersionIsMissing(t *testing.T) {
	discoveryClient := &fakediscovery.FakeDiscovery{Fake: &clienttesting.Fake{}}

	exists, err := ResourceExists(discoveryClient, "security.openshift.io/v1", "SecurityContextConstraints")

	require.NoError(t, err)
	assert.False(t, exists)
}

func TestResourceExistsReturnsTargetDiscoveryError(t *testing.T) {
	expectedErr := errors.New("target API group discovery failed")
	discoveryClient := &fakediscovery.FakeDiscovery{Fake: &clienttesting.Fake{}}
	discoveryClient.PrependReactor("get", "resource", func(clienttesting.Action) (bool, runtime.Object, error) {
		return true, nil, expectedErr
	})

	exists, err := ResourceExists(discoveryClient, "security.openshift.io/v1", "SecurityContextConstraints")

	assert.False(t, exists)
	require.ErrorIs(t, err, expectedErr)
}
