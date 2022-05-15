package helpers

import (
	"github.com/SakuraBurst/gitlab-bot/pkg/basa_dannih"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	log "github.com/sirupsen/logrus"
)

func OnlyNewMrs(openedMergeRequests []models.MergeRequest, bd basa_dannih.BDInterface) (unWrittenOpenedMergeRequests []models.MergeRequest, writtenMergeRequests int, ok bool) {
	if len(openedMergeRequests) == 0 {
		return nil, 0, false
	}
	unWrittenOpenedMergeRequests = make([]models.MergeRequest, 0, len(openedMergeRequests))
	for _, v := range openedMergeRequests {
		if !bd.ReadFromBd(v.Iid) {
			unWrittenOpenedMergeRequests = append(unWrittenOpenedMergeRequests, v)
			writtenMergeRequests++
			bd.WriteToBD(v.Iid)
			log.WithField("basa", bd).Infof("база данных поплнена айдишником %d", v.Iid)
		}
	}
	return unWrittenOpenedMergeRequests, writtenMergeRequests, writtenMergeRequests > 0
}

func WriteMrsToBd(bd basa_dannih.BasaDannihMySQLPostgresMongoPgAdmin777, mrs ...models.MergeRequest) {
	for _, v := range mrs {
		bd.WriteToBD(v.Iid)
	}
}
