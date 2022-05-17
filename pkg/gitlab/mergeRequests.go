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

type MergeRequestTransfer struct {
	mergeRequest *models.MergeRequest
	error        error
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
		return getMrsWithDiffs(g, openedMergeRequests)
	}
	log.WithFields(log.Fields{"Количество мрок со статусом opened": openedMergeRequests.Length}).Info("парсер закончил работу")
	return openedMergeRequests, nil
}

func getMrsWithDiffs(g Gitlab, mri *MergeRequestsInfo) (*MergeRequestsInfo, error) {
	responseWaiters := make(chan MergeRequestTransfer, mri.Length)
	closed := make(chan bool)
	var isClosed = func() bool {
		switch {
		case <-closed:
			return true
		default:
			return false
		}
	}
	for _, v := range mri.MergeRequests {
		go g.getMRDiffs(v.Iid, responseWaiters, isClosed)
	}

	openedMergeRequestsWithDiffs := MergeRequestsInfo{
		Length:        mri.Length,
		On:            mri.On,
		MergeRequests: make([]models.MergeRequest, 0, mri.Length),
	}

	for i := 0; i < mri.Length; i++ {
		result := <-responseWaiters
		if result.error != nil {
			close(closed)
			close(responseWaiters)
			return nil, result.error
		}
		openedMergeRequestsWithDiffs.MergeRequests = append(openedMergeRequestsWithDiffs.MergeRequests, *result.mergeRequest)
	}

	close(responseWaiters)
	return &openedMergeRequestsWithDiffs, nil
}

func (g Gitlab) getMRDiffs(iid int, resChan chan MergeRequestTransfer, isChannelClosed func() bool) {
	log.WithFields(log.Fields{"iid": iid}).Info("получение отдельного открытого мр с доп даннымми")
	url, headers, err := g.getSingleMergeRequestWithChangesURL(iid)
	if err != nil {
		if isChannelClosed() {
			return
		}
		resChan <- MergeRequestTransfer{
			mergeRequest: nil,
			error:        err,
		}
	}
	resp, err := clients.Get(url.String(), headers)

	if err != nil {
		if isChannelClosed() {
			return
		}
		resChan <- MergeRequestTransfer{
			mergeRequest: nil,
			error:        err,
		}
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// TODO: подумоть
			log.Error(err)
		}
	}()
	mrWithFileChanges, err := decodeSingleMergeRequestItem(resp)
	if err != nil {
		if err != nil {
			if isChannelClosed() {
				return
			}
			resChan <- MergeRequestTransfer{
				mergeRequest: nil,
				error:        err,
			}
		}
	}
	log.WithFields(log.Fields{"iid": iid, "mrWithFileChanges": mrWithFileChanges}).Info("мр успешно получен")
	resChan <- MergeRequestTransfer{
		mergeRequest: mrWithFileChanges,
		error:        nil,
	}
}
