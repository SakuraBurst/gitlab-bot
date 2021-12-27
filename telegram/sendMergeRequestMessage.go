package telegram

import (
	"bytes"

	"github.com/SakuraBurst/gitlab-bot/models"
	"github.com/SakuraBurst/gitlab-bot/templates"
)

func (t TelegramBot) SendMergeRequestMessage(mergeRequests models.MergeRequests, newMr, withDiffs bool) {

	buff := bytes.NewBuffer(nil)
	templates.GetRightTemplate(newMr, withDiffs).Execute(buff, mergeRequests)
	t.sendMessage(buff.String())

}
