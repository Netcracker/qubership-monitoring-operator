package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPRouteSpecWithDefaultsFilters(t *testing.T) {
	spec := map[string]interface{}{
		"rules": []interface{}{map[string]interface{}{
			"filters": []interface{}{
				map[string]interface{}{"requestMirror": map[string]interface{}{
					"backendRef": map[string]interface{}{"name": "mirror"},
					"fraction":   map[string]interface{}{"numerator": 5},
				}},
				map[string]interface{}{"cors": map[string]interface{}{}},
			},
		}},
	}

	got, err := httpRouteSpecWithDefaults(spec)
	require.NoError(t, err)

	rules := got.(map[string]interface{})["rules"].([]interface{})
	filters := rules[0].(map[string]interface{})["filters"].([]interface{})
	mirror := filters[0].(map[string]interface{})["requestMirror"].(map[string]interface{})
	backendRef := mirror["backendRef"].(map[string]interface{})
	assert.Equal(t, "", backendRef["group"])
	assert.Equal(t, "Service", backendRef["kind"])
	assert.Equal(t, float64(100), mirror["fraction"].(map[string]interface{})["denominator"])
	assert.Equal(t, float64(5), filters[1].(map[string]interface{})["cors"].(map[string]interface{})["maxAge"])
}

func TestHTTPRouteSpecWithDefaultsEdgeCases(t *testing.T) {
	got, err := httpRouteSpecWithDefaults(nil)
	require.NoError(t, err)
	assert.Nil(t, got)

	_, err = httpRouteSpecWithDefaults(make(chan int))
	require.Error(t, err, "a value that cannot be encoded as JSON")

	_, err = httpRouteSpecWithDefaults([]interface{}{"not", "an", "object"})
	require.Error(t, err, "a spec that is not a JSON object")

	assert.False(t, httpRouteSpecMatches(map[string]interface{}{}, make(chan int)))
}
