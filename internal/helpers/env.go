package helpers

import (
	"errors"
	"os"
	"strings"
)

func GetEnvMap() map[string]string {
	envMap := make(map[string]string)
	for _, v := range os.Environ() {
		env := strings.Split(v, "=")
		envMap[env[0]] = env[1]
	}
	return envMap
}

func CheckForEnv(env []string) error {
	for _, s := range env {
		if _, ok := os.LookupEnv(s); !ok {
			return errors.New("Отсутствует " + s)
		}
	}
	return nil
}
