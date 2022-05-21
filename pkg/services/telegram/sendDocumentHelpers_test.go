package telegram

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

var invalidSendFileURLErrorStringMock = "parse \"https://api.telegram.org/bot_394e2904dfsfs234)Ue)R(UEWR#$%@%/sendDocument\": invalid URL escape \"%@%\""

var telegramValidSendFileURLMock = "https://api.telegram.org/bot1/sendDocument"

var multipartFormValueMockBuilder = "--%s\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n\r\n--%s\r\nContent-Disposition: form-data; name=\"document\"; filename=\"\"\r\nContent-Type: application/octet-stream\r\n\r\n\r\n--%s--\r\n"

func TestBot_CreateSendFileURL_InvalidURL(t *testing.T) {
	sendMessageURL, err := telegramBotInvalidTokenMock.createSendFileURL()
	assert.Nil(t, sendMessageURL)
	assert.NotNil(t, err)
	assert.EqualError(t, err, invalidSendFileURLErrorStringMock)
}

func TestBot_CreateSendFileURL(t *testing.T) {
	sendMessageURL, err := telegramBotMock.createSendFileURL()
	assert.NotNil(t, sendMessageURL)
	assert.Nil(t, err)
	assert.Equal(t, telegramValidSendFileURLMock, sendMessageURL.String())
}

func TestBot_MakeBodyToSendFile_InvalidFile(t *testing.T) {
	body, headers, err := telegramBotInvalidTokenMock.makeBodyAndHeadersToSendFile(nil, "")
	assert.Nil(t, body)
	assert.Nil(t, headers)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "invalid argument")

}

func TestBot_MakeBodyToSendFile(t *testing.T) {
	file, err := os.Create("test.txt")
	require.Nil(t, err)
	_, err = file.Write([]byte("test"))
	require.Nil(t, err)
	body, headers, err := telegramBotInvalidTokenMock.makeBodyAndHeadersToSendFile(file, "")
	assert.NotNil(t, body)
	assert.NotNil(t, headers)
	assert.Nil(t, err)
	buffer := bytes.NewBuffer(nil)
	b, err := io.Copy(buffer, body)
	require.Nil(t, err)
	assert.Equal(t, int64(352), b)
	boundary := getBoundary(headers.Get("Content-type"))
	multipartFormValueMock := fmt.Sprintf(multipartFormValueMockBuilder, boundary, boundary, boundary)
	assert.Equal(t, multipartFormValueMock, buffer.String())
	assert.Contains(t, headers, http.CanonicalHeaderKey("Content-type"))
	err = file.Close()
	require.Nil(t, err)
	err = os.Remove("test.txt")
	require.Nil(t, err)
}

func TestPrepareSendMessageRequestHelper(t *testing.T) {
	sendDocumentHelper := prepareSendMessageRequestHelper(nil)
	assert.Nil(t, sendDocumentHelper.File)
	assert.NotNil(t, sendDocumentHelper.Writer)
	assert.NotNil(t, sendDocumentHelper.MainBuffer)
}

func getBoundary(contentType string) string {
	boundaryKey := "boundary="
	index := strings.Index(contentType, boundaryKey)
	log.Println(contentType)

	return contentType[index+len(boundaryKey):]
}
