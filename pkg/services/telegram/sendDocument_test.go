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

var sendFileURLMock = "https://api.telegram.org/bot1/sendDocument"

func TestBot_SendFile_CreateURLError(t *testing.T) {
	err := telegramBotInvalidTokenMock.SendDocument(nil, "")
	assert.NotNil(t, err)
	assert.EqualError(t, err, invalidSendFileURLErrorStringMock)
}

func TestBot_SendFile_MakeBodyError(t *testing.T) {
	err := telegramBotMock.SendDocument(nil, "")
	assert.NotNil(t, err)
	assert.EqualError(t, err, "invalid argument")
}

func TestBot_SendFile_RequestError(t *testing.T) {
	clients.EnableMock()
	file, err := os.Create("test.txt")
	require.Nil(t, err)
	mockError := errors.New("mockError")
	clients.Mocks.AddMock(sendFileURLMock, clients.Mock{
		Response: nil,
		Err:      mockError,
	})
	err = telegramBotMock.SendDocument(file, "")

	assert.NotNil(t, err)
	assert.EqualError(t, err, mockError.Error())
	err = file.Close()
	require.Nil(t, err)
	err = os.Remove("test.txt")
	require.Nil(t, err)
	clients.DisableMock()
}

func TestBot_SendFile_TelegramError(t *testing.T) {
	clients.EnableMock()
	file, err := os.Create("test.txt")
	require.Nil(t, err)
	telegramError := models.TelegramError{
		Ok:          false,
		ErrorCode:   http.StatusUnauthorized,
		Description: http.StatusText(http.StatusUnauthorized),
	}
	telegramErrorBody, err := json.Marshal(telegramError)
	require.Nil(t, err)
	body := bytes.NewReader(telegramErrorBody)
	clients.Mocks.AddMock(sendFileURLMock, clients.Mock{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusUnauthorized),
			StatusCode: http.StatusUnauthorized,
			Body:       io.NopCloser(body),
		},
		Err: nil,
	})
	err = telegramBotMock.SendDocument(file, "")
	assert.NotNil(t, err)
	assert.EqualError(t, err, telegramError.Error())
	err = file.Close()
	require.Nil(t, err)
	err = os.Remove("test.txt")
	require.Nil(t, err)
	clients.DisableMock()
}

func TestBot_SendFile_InvalidBody(t *testing.T) {
	clients.EnableMock()
	file, err := os.Create("test.txt")
	require.Nil(t, err)
	invalidBody, err := os.Open("123sdcfc90")
	require.NotNil(t, err)
	clients.Mocks.AddMock(sendFileURLMock, clients.Mock{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			Body:       invalidBody,
		},
		Err: nil,
	})
	err = telegramBotMock.SendDocument(file, "")
	assert.Nil(t, err)
	err = file.Close()
	require.Nil(t, err)
	err = os.Remove("test.txt")
	require.Nil(t, err)
	clients.DisableMock()
}

func TestBot_SendFile(t *testing.T) {
	clients.EnableMock()
	file, err := os.Create("test.txt")
	require.Nil(t, err)
	nilReader := bytes.NewReader(nil)
	clients.Mocks.AddMock(sendFileURLMock, clients.Mock{
		Response: &http.Response{
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(nilReader),
		},
		Err: nil,
	})
	err = telegramBotMock.SendDocument(file, "")
	assert.Nil(t, err)
	err = file.Close()
	require.Nil(t, err)
	err = os.Remove("test.txt")
	require.Nil(t, err)
	clients.DisableMock()
}
