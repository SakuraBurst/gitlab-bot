package telegram

import (
	"github.com/SakuraBurst/gitlab-bot/api/clients"
	log "github.com/sirupsen/logrus"
)

func (t Bot) SendMessage(text string) error {
	tgRequest := map[string]string{
		"chat_id":    t.mainChannel,
		"text":       text,
		"parse_mode": "html",
	}
	log.WithFields(log.Fields{"tgRequest": tgRequest}).Info("начата отправка сообщения в телеграм")
	sendMessageURL, headers, err := t.CreateSendMessageURL()
	if err != nil {
		return err
	}
	response, err := clients.Post(sendMessageURL.String(), tgRequest, headers)
	if err != nil {
		return err
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			// TODO: опять подумоть
			log.Error(err)
		}
	}()
	err = decodeTelegramResponse(response)
	if err != nil {
		return err
	}
	log.Info("Сообщение успешно отправлено")
	return nil
}
