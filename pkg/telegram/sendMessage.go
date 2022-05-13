package telegram

import (
	"encoding/json"
	"github.com/SakuraBurst/gitlab-bot/api/clients"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (t Bot) SendMessage(text string) error {
	tgRequest := map[string]string{
		"chat_id":    t.mainChannel,
		"text":       text,
		"parse_mode": "html",
	}
	log.WithFields(log.Fields{"tgRequest": tgRequest}).Info("начата отправка сообщения в телеграм")
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	response, err := clients.Post("https://api.telegram.org/bot"+t.token+"/sendMessage", tgRequest, headers)
	if err != nil {
		return err
	}
	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	decoder := json.NewDecoder(response.Body)

	if response.StatusCode != http.StatusOK {
		test := make(map[string]interface{})
		err = decoder.Decode(&test)
		if err != nil {
			log.Panic(err)
		}
		log.WithFields(log.Fields{"ошибка отправки в телеграм": test}).Fatal()
	}
	log.Info("Сообщение успешно отправлено")
	return nil
}
