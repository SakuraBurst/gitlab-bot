package telegram

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/SakuraBurst/gitlab-bot/models"
	"github.com/SakuraBurst/gitlab-bot/templates"
	log "github.com/sirupsen/logrus"
)

func SendMessage(mergeRequests models.MergeRequests, newMr, withDiffs bool, channel, token string) {

	buff := bytes.NewBuffer(nil)
	templates.GetRightTemplate(newMr, withDiffs).Execute(buff, mergeRequests)

	tgRequest := map[string]string{
		"chat_id":    channel,
		"text":       buff.String(),
		"parse_mode": "html",
	}

	log.WithFields(log.Fields{"tgRequest": tgRequest}).Info("начата отправка сообщения в телеграм")
	testBytes, err := json.Marshal(tgRequest)
	if err != nil {
		log.Fatal(err)
	}
	reader := bytes.NewReader(testBytes)
	respon, err := http.Post("https://api.telegram.org/bot"+token+"/sendMessage", "application/json", reader)
	if err != nil {
		log.Fatal(err)
	}
	defer respon.Body.Close()
	decoder := json.NewDecoder(respon.Body)

	if respon.StatusCode != http.StatusOK {
		test := make(map[string]interface{})
		decoder.Decode(&test)
		log.WithFields(log.Fields{"ошибка отправки в телеграм": test}).Fatal()
	}
	log.Info("Сообщение успешно отправлено")
}
