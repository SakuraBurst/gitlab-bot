package telegram

import (
	"encoding/json"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func decodeTelegramResponse(response *http.Response) error {
	decoder := json.NewDecoder(response.Body)
	if response.StatusCode != http.StatusOK {
		log.Println(response.StatusCode, response.Status)
		var tgError models.TelegramError
		err := decoder.Decode(&tgError)
		if err != nil {
			return err
		}
		return tgError
	}
	return nil
}
