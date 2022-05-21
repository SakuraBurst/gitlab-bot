package telegram

import (
	"github.com/SakuraBurst/gitlab-bot/api/clients"
	log "github.com/sirupsen/logrus"
	"os"
)

func (t Bot) SendDocument(file *os.File, fileName string) error {
	sendFileURL, err := t.createSendFileURL()
	if err != nil {
		return err
	}
	body, headers, err := t.makeBodyAndHeadersToSendFile(file, fileName)
	if err != nil {
		return err
	}
	resp, err := clients.PostStream(sendFileURL.String(), body, headers)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// TODO: опять подумоть
			log.Error(err)
		}
	}()
	err = decodeTelegramResponse(resp)
	if err != nil {
		return err
	}
	return nil
}
