package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/SakuraBurst/gitlab-bot/models"
	"github.com/SakuraBurst/gitlab-bot/templates"
)

func SendMessage(mrWithDiffs models.MergeRequests) {
	buff := bytes.NewBuffer([]byte{})
	templates.TelegramMessageTemplate.Execute(buff, mrWithDiffs)

	testRequest := map[string]string{
		"chat_id":    "@mrchicki",
		"text":       buff.String(),
		"parse_mode": "html",
	}
	testBytes, err := json.Marshal(testRequest)
	if err != nil {
		log.Fatal(err)
	}
	reader := bytes.NewReader(testBytes)
	respon, err := http.Post("https://api.telegram.org/bot5021252898:AAFJr-XK1_pTKNEW3Ju7tvT-z1VOb75zycw/sendMessage", "application/json", reader)
	if err != nil {
		log.Fatal(err)
	}
	defer respon.Body.Close()
	decoder := json.NewDecoder(respon.Body)

	if respon.StatusCode != http.StatusOK {
		test := make(map[string]interface{})
		decoder.Decode(&test)
		log.Fatal(test)
	}
	fmt.Println(respon.StatusCode)
}
