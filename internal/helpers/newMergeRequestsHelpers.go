package helpers

import (
	"github.com/SakuraBurst/gitlab-bot/pkg/BasaDannih"
	"github.com/SakuraBurst/gitlab-bot/pkg/gitlab"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	log "github.com/sirupsen/logrus"
)

func OnlyNewMrs(openedMergeRequests gitlab.MergeRequestsInfo, bd BasaDannih.BasaDannihMySQLPostgresMongoPgAdmin777) (gitlab.MergeRequestsInfo, bool) {
	onlyNewMrs := gitlab.MergeRequestsInfo{On: openedMergeRequests.On}
	if openedMergeRequests.Length == 0 {
		return onlyNewMrs, false
	}

	for _, v := range openedMergeRequests.MergeRequests {
		if !bd.ReadFromBd(v.Iid) {
			onlyNewMrs.MergeRequests = append(onlyNewMrs.MergeRequests, v)
			onlyNewMrs.Length++
			bd.WriteToBD(v.Iid)
			log.WithField("basa", bd).Infof("база данных поплнена айдишником %d", v.Iid)
		}
	}
	return onlyNewMrs, onlyNewMrs.Length > 0
}

func WriteMrsToBd(bd BasaDannih.BasaDannihMySQLPostgresMongoPgAdmin777, mrs ...models.MergeRequestListItem) {
	for _, v := range mrs {
		bd.WriteToBD(v.Iid)
	}
}
