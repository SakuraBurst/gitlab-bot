package telegram

import (
	"bytes"
	log "github.com/sirupsen/logrus"

	"github.com/SakuraBurst/gitlab-bot/models"
	"github.com/SakuraBurst/gitlab-bot/templates"
)

func (t Bot) SendMergeRequestMessage(mergeRequests models.MergeRequests, newMr, withDiffs bool) {
	buff := bytes.NewBuffer(nil)
	if err := templates.GetRightTemplate(newMr, withDiffs).Execute(buff, mergeRequests); err != nil {
		log.Fatal(err)
	}
	t.sendMessage(buff.String())
}
