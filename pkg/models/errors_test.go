package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGitlabError_Error(t *testing.T) {
	glError := GitlabError{Message: "test"}
	assert.Error(t, glError)
	assert.Equal(t, glError.Error(), "test")
}

func TestTelegramError_Error(t *testing.T) {
	tlError := TelegramError{
		Ok:          false,
		ErrorCode:   401,
		Description: "test",
	}
	assert.Error(t, tlError)
	assert.Equal(t, tlError.Error(), "Code 401, Message: test")
}
