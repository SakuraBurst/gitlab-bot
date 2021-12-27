package telegram

import (
	"bytes"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (t Bot) sendMessage(text string) {
	tgRequest := map[string]string{
		"chat_id":    t.mainChannel,
		"text":       text,
		"parse_mode": "html",
	}

	log.WithFields(log.Fields{"tgRequest": tgRequest}).Info("начата отправка сообщения в телеграм")
	testBytes, err := json.Marshal(tgRequest)
	if err != nil {
		log.Fatal(err)
	}
	reader := bytes.NewReader(testBytes)
	response, err := http.Post("https://api.telegram.org/bot"+t.token+"/sendMessage", "application/json", reader)
	if err != nil {
		log.Fatal(err)
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
}
