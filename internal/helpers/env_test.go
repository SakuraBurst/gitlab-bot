package helpers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGetEnvMap(t *testing.T) {
	testEnv := "TestGetEnvMap"
	err := os.Setenv(testEnv, testEnv)
	require.Nil(t, err)
	envMap := GetEnvMap()
	assert.NotNil(t, envMap)
	assert.Equal(t, envMap[testEnv], testEnv)
	assert.Empty(t, envMap["NEVER_USED_ENV"])
}

func TestCheckForEnvError(t *testing.T) {
	neededEnv := []string{"TestCheckForEnvError"}
	err := CheckForEnv(neededEnv)
	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.Equal(t, "Отсутствует TestCheckForEnvError", err.Error())
}

func TestCheckForEnvOK(t *testing.T) {
	testEnv := "TestCheckForEnvOK"
	err := os.Setenv(testEnv, testEnv)
	require.Nil(t, err)
	neededEnv := []string{testEnv}
	err = CheckForEnv(neededEnv)
	assert.Nil(t, err)
}
