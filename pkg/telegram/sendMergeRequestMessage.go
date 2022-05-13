package telegram

import (
	"bytes"
	"github.com/SakuraBurst/gitlab-bot/internal/templates"
	"github.com/SakuraBurst/gitlab-bot/pkg/gitlab"
	log "github.com/sirupsen/logrus"
)

func (t Bot) SendMergeRequestMessage(mergeRequests gitlab.MergeRequestsInfo, newMr, withDiffs bool) error {
	buff := bytes.NewBuffer(nil)
	if err := templates.GetRightTemplate(newMr, withDiffs).Execute(buff, mergeRequests); err != nil {
		log.Fatal(err)
	}
	return t.SendMessage(buff.String())
}
