package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGitlabError_Error(t *testing.T) {
	glError := GitlabError{Message: "test"}
	assert.Error(t, glError)
	assert.Equal(t, "test", glError.Error())
}

func TestTelegramError_Error(t *testing.T) {
	tlError := TelegramError{
		Ok:          false,
		ErrorCode:   401,
		Description: "test",
	}
	assert.Error(t, tlError)
	assert.EqualError(t, tlError, "Code 401, Message: test")
}
