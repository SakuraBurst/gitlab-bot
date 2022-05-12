package telegram

import (
	"bytes"
	"github.com/SakuraBurst/gitlab-bot/internal/templates"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	log "github.com/sirupsen/logrus"
)

func (t Bot) SendMergeRequestMessage(mergeRequests models.MergeRequestsInfo, newMr, withDiffs bool) {
	buff := bytes.NewBuffer(nil)
	if err := templates.GetRightTemplate(newMr, withDiffs).Execute(buff, mergeRequests); err != nil {
		log.Fatal(err)
	}
	t.SendMessage(buff.String())
}
