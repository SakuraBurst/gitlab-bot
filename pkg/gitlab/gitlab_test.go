package gitlab

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGitlabConn(t *testing.T) {
	testString := "test"
	glConn := NewGitlabConn(true, testString, testString, testString)
	assert.NotNil(t, glConn)
	assert.True(t, glConn.WithDiffs)
	assert.Equal(t, glConn.repo, testString)
	assert.Equal(t, glConn.token, testString)
	assert.Equal(t, glConn.url, testString)

}
