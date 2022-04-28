package helpers

import (
	"os"
	"strings"
)

func GetOsEnvMap() map[string]string {
	envMap := make(map[string]string)
	for _, v := range os.Environ() {
		env := strings.Split(v, "=")
		envMap[env[0]] = env[1]
	}
	return envMap
}
