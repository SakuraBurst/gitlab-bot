package telegram

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var telegramBotMock = Bot{
	token:       "1",
	mainChannel: "1",
}

var telegramValidURLMock = "https://api.telegram.org/bot1/sendMessage"

var telegramBotInvalidTokenMock = Bot{
	token:       "_394e2904dfsfs234)Ue)R(UEWR#$%@%",
	mainChannel: "",
}

var invalidURLErrorStringMock = "parse \"https://api.telegram.org/bot_394e2904dfsfs234)Ue)R(UEWR#$%@%/sendMessage\": invalid URL escape \"%@%\""

func TestBot_CreateSendMessageURL_InvalidURL(t *testing.T) {
	sendMessageURL, headers, err := telegramBotInvalidTokenMock.CreateSendMessageURL()
	assert.Nil(t, sendMessageURL)
	assert.Nil(t, headers)
	assert.NotNil(t, err)
	assert.EqualError(t, err, invalidURLErrorStringMock)
}

func TestBot_CreateSendMessageURL(t *testing.T) {
	sendMessageURL, headers, err := telegramBotMock.CreateSendMessageURL()
	assert.NotNil(t, sendMessageURL)
	assert.NotNil(t, headers)
	assert.Nil(t, err)
	assert.Equal(t, sendMessageURL.String(), telegramValidURLMock)
	assert.Contains(t, headers, http.CanonicalHeaderKey("content-type"))
}
