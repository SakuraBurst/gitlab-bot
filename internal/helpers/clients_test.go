package helpers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestValidHeadersNil(t *testing.T) {
	headers := ValidHeaders(nil)
	assert.NotNil(t, headers)
}

func TestValidHeadersNotNIl(t *testing.T) {
	originalHeaders := http.Header{}
	headers := ValidHeaders(originalHeaders)
	assert.Equal(t, originalHeaders, headers)
}
