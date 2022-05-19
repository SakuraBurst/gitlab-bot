package telegram

import (
	"encoding/json"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	"net/http"
)

func decodeTelegramResponse(response *http.Response) error {
	decoder := json.NewDecoder(response.Body)
	if response.StatusCode != http.StatusOK {
		var tgError models.TelegramError
		err := decoder.Decode(&tgError)
		if err != nil {
			return err
		}
		return tgError
	}
	return nil
}
