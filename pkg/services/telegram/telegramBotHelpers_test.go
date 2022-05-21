package telegram

import (
	"bytes"
	"encoding/json"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"testing"
)

var invalidArgumentErrorStringMock = "invalid argument"

func TestDecodeTelegramResponse_ErrorInvalidBody(t *testing.T) {
	invalidBody, err := os.Open("123sdcfc90")
	require.NotNil(t, err)
	response := &http.Response{
		Status:     http.StatusText(http.StatusUnauthorized),
		StatusCode: http.StatusUnauthorized,
		Body:       invalidBody,
	}
	err = decodeTelegramResponse(response)
	assert.NotNil(t, err)
	assert.EqualError(t, err, invalidArgumentErrorStringMock)
}

func TestDecodeTelegramResponse_ErrorBody(t *testing.T) {
	tgError := models.TelegramError{
		Ok:          false,
		ErrorCode:   http.StatusUnauthorized,
		Description: http.StatusText(http.StatusUnauthorized),
	}
	tgErrorBytes, err := json.Marshal(tgError)
	require.Nil(t, err)
	tgErrorBody := bytes.NewReader(tgErrorBytes)
	response := &http.Response{
		Status:     http.StatusText(http.StatusUnauthorized),
		StatusCode: http.StatusUnauthorized,
		Body:       io.NopCloser(tgErrorBody),
	}
	err = decodeTelegramResponse(response)
	assert.NotNil(t, err)
	assert.EqualError(t, err, tgError.Error())
}

func TestDecodeTelegramResponse(t *testing.T) {
	response := &http.Response{
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
	}
	err := decodeTelegramResponse(response)
	assert.Nil(t, err)
}
