package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnstructuredValuesEqual(t *testing.T) {
	tests := []struct {
		name string
		a, b interface{}
		want bool
	}{
		{name: "identical values", a: map[string]interface{}{"port": 80}, b: map[string]interface{}{"port": 80}, want: true},
		{name: "decoded and in-memory numbers", a: map[string]interface{}{"port": float64(80)}, b: map[string]interface{}{"port": int32(80)}, want: true},
		{name: "different values", a: map[string]interface{}{"port": 80}, b: map[string]interface{}{"port": 81}},
		{name: "first value cannot be encoded", a: make(chan int), b: map[string]interface{}{}},
		{name: "second value cannot be encoded", a: map[string]interface{}{}, b: make(chan int)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, UnstructuredValuesEqual(tt.a, tt.b))
		})
	}
}
