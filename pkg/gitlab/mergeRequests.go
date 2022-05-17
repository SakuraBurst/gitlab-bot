package gitlab

import (
	"github.com/SakuraBurst/gitlab-bot/api/clients"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	log "github.com/sirupsen/logrus"
	"time"
)

type MergeRequestsInfo struct {
	Length        int
	On            time.Time
	MergeRequests []models.MergeRequest
}

func (g Gitlab) MergeRequests() (*MergeRequestsInfo, error) {
	log.WithFields(log.Fields{"repo": g.repo}).Info("парсер начал работу")
	url, headers, err := g.getMergeRequestURL()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resp, err := clients.Get(url.String(), headers)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// TODO: подумоть
			log.Error(err)
		}
	}()

	openedMergeRequests, err := decodeMergeRequestsInfo(resp)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if g.WithDiffs {
		return getMrsWithDiffs(g, openedMergeRequests), nil
	}
	log.WithFields(log.Fields{"Количество мрок со статусом opened": openedMergeRequests.Length}).Info("парсер закончил работу")
	return openedMergeRequests, nil
}

func getMrsWithDiffs(g Gitlab, mri *MergeRequestsInfo) *MergeRequestsInfo {
	responseWaiters := make(chan models.MergeRequest, mri.Length)
	for _, v := range mri.MergeRequests {
		go g.getMRDiffs(v.Iid, responseWaiters)
	}

	openedMergeRequestsWithDiffs := MergeRequestsInfo{
		Length:        mri.Length,
		On:            mri.On,
		MergeRequests: make([]models.MergeRequest, 0, mri.Length),
	}

	for i := 0; i < mri.Length; i++ {
		openedMergeRequestsWithDiffs.MergeRequests = append(openedMergeRequestsWithDiffs.MergeRequests, <-responseWaiters)
	}

	close(responseWaiters)
	return &openedMergeRequestsWithDiffs
}

func (g Gitlab) getMRDiffs(iid int, resChan chan models.MergeRequest) {
	log.WithFields(log.Fields{"iid": iid}).Info("получение отдельного открытого мр с доп даннымми")
	url, headers, err := g.getSingleMergeRequestWithChangesURL(iid)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := clients.Get(url.String(), headers)

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	mrWithFileChanges, err := decodeSingleMergeRequestItem(resp)
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{"iid": iid, "mrWithFileChanges": mrWithFileChanges}).Info("мр успешно получен")
	resChan <- *mrWithFileChanges
}
