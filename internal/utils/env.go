package utils

import (
	"os"
)

func GetEnvStr(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
