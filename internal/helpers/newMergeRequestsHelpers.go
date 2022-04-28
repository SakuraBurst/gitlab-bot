package helpers

import (
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	log "github.com/sirupsen/logrus"
)

func OnlyNewMrs(openedMergeRequests models.MergeRequests, bd models.BasaDannihMySQLPostgresMongoPgAdmin777) (models.MergeRequests, bool) {
	onlyNewMrs := models.MergeRequests{On: openedMergeRequests.On}
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

func WriteMrsToBd(bd models.BasaDannihMySQLPostgresMongoPgAdmin777, mrs ...models.MergeRequestFileChanges) {
	for _, v := range mrs {
		bd.WriteToBD(v.Iid)
	}
}
