package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SakuraBurst/gitlab-bot/api/clients"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestBot_SendMessage_CreateURLError(t *testing.T) {
	err := telegramBotInvalidTokenMock.SendMessage("test")
	assert.NotNil(t, err)
	assert.EqualError(t, err, invalidURLErrorStringMock)
}

func TestBot_SendMessage_RequestError(t *testing.T) {
	clients.EnableMock()
	responseErr := errors.New("response error")
	clients.Mocks.AddMock(telegramValidURLMock, clients.Mock{
		Response: nil,
		Err:      responseErr,
	})
	err := telegramBotMock.SendMessage("test")
	assert.NotNil(t, err)
	assert.EqualError(t, err, responseErr.Error())
	clients.DisableMock()
}

func TestBot_SendMessage_TelegramError(t *testing.T) {
	clients.EnableMock()
	telegramError := models.TelegramError{
		Ok:          false,
		ErrorCode:   http.StatusUnauthorized,
		Description: http.StatusText(http.StatusUnauthorized),
	}
	telegramErrorBytes, err := json.Marshal(telegramError)
	require.Nil(t, err)
	telegramErrorBody := bytes.NewReader(telegramErrorBytes)
	clients.Mocks.AddMock(telegramValidURLMock, clients.Mock{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusUnauthorized),
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(telegramErrorBody),
		},
		Err: nil,
	})
	err = telegramBotMock.SendMessage("test")
	assert.NotNil(t, err)
	assert.EqualError(t, err, telegramError.Error())
	clients.DisableMock()
}

func TestBot_SendMessage_InvalidBody(t *testing.T) {
	clients.EnableMock()
	invalidBody, err := os.Open("123sdcfc90")
	require.NotNil(t, err)
	clients.Mocks.AddMock(telegramValidURLMock, clients.Mock{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			Body:       invalidBody,
		},
		Err: nil,
	})
	err = telegramBotMock.SendMessage("test")
	assert.Nil(t, err)
	clients.DisableMock()
}

func TestBot_SendMessage(t *testing.T) {
	clients.EnableMock()
	emptyBody := bytes.NewReader(nil)
	clients.Mocks.AddMock(telegramValidURLMock, clients.Mock{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(emptyBody),
		},
		Err: nil,
	})
	err := telegramBotMock.SendMessage("test")
	assert.Nil(t, err)
	clients.DisableMock()
}
