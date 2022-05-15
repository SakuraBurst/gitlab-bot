package clients

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestValidHeadersNil(t *testing.T) {
	headers := validHeaders(nil)
	assert.NotNil(t, headers)
}

func TestValidHeadersNotNIl(t *testing.T) {
	originalHeaders := http.Header{}
	headers := validHeaders(originalHeaders)
	assert.Equal(t, originalHeaders, headers)
}
