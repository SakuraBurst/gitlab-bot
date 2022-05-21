package telegram

import (
	"fmt"
	"net/http"
	"net/url"
)

func (t Bot) createSendMessageURL() (*url.URL, http.Header, error) {
	rawUrl := fmt.Sprintf("%s/bot%s/sendMessage", telegramApi, t.token)
	sendMessageURL, err := url.Parse(rawUrl)
	if err != nil {
		return nil, nil, err
	}
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	return sendMessageURL, headers, nil
}
