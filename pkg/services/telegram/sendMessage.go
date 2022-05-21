package telegram

import (
	"github.com/SakuraBurst/gitlab-bot/api/clients"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	log "github.com/sirupsen/logrus"
)

func (t Bot) SendMessage(text string) error {
	tgRequest := models.TelegramMessageRequest{
		ChatID:    t.mainChannel,
		Text:      text,
		ParseMode: "html",
	}
	log.WithFields(log.Fields{"tgRequest": tgRequest}).Info("начата отправка сообщения в телеграм")
	sendMessageURL, headers, err := t.createSendMessageURL()
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
