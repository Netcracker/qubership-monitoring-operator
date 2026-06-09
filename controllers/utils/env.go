package utils

import "os"

var (
	_defaultEnvValues = map[string]string{
		"RECONCILIATION_INTERVAL": "60",
	}
)

func GetEnvWithDefaultValue(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return _defaultEnvValues[key]
	}
	return value
}
