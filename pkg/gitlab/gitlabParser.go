package gitlab

import (
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (g Gitlab) MergeRequests() (models.MergeRequestsInfo, error) {
	log.WithFields(log.Fields{"repo": g.repo}).Info("парсер начал работу")
	request, err := getMergeRequest(g)
	if err != nil {
		log.Error(err)
		return models.MergeRequestsInfo{}, err
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Error(err)
		return models.MergeRequestsInfo{}, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	openedMergeRequests, err := decodeMergeRequestsInfo(resp.Body)
	if err != nil {
		log.Error(err)
		return models.MergeRequestsInfo{}, err
	}

	if g.WithDiffs {
		return getMrsWithDiffs(g, openedMergeRequests), nil
	}
	log.WithFields(log.Fields{"Количество мрок со статусом opened": openedMergeRequests.Length}).Info("парсер закончил работу")
	return openedMergeRequests, nil
}

func getMrsWithDiffs(g Gitlab, mri models.MergeRequestsInfo) models.MergeRequestsInfo {
	responseWaiters := make(chan models.MergeRequestListItem, mri.Length)
	for _, v := range mri.MergeRequests {
		go g.getMRDiffs(v.Iid, responseWaiters)
	}

	openedMergeRequestsWithDiffs := models.MergeRequestsInfo{
		Length:        mri.Length,
		On:            mri.On,
		MergeRequests: make([]models.MergeRequestListItem, 0, mri.Length),
	}

	for i := 0; i < mri.Length; i++ {
		openedMergeRequestsWithDiffs.MergeRequests = append(openedMergeRequestsWithDiffs.MergeRequests, <-responseWaiters)
	}

	close(responseWaiters)
	return openedMergeRequestsWithDiffs
}

func (g Gitlab) getMRDiffs(iid int, resChan chan models.MergeRequestListItem) {
	log.WithFields(log.Fields{"iid": iid}).Info("получение отдельного открытого мр с доп даннымми")
	request, err := getSingleMergeRequestWithChanges(g, iid)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	mrWithFileChanges, err := decodeSingleMergeRequestItem(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{"iid": iid, "mrWithFileChanges": mrWithFileChanges}).Info("мр успешно получен")
	resChan <- mrWithFileChanges
}
