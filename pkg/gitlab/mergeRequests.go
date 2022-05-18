package gitlab

import (
	"github.com/SakuraBurst/gitlab-bot/api/clients"
	log "github.com/sirupsen/logrus"
)

func (g Gitlab) getAllOpenedMergeRequests() (*MergeRequestsInfo, error) {
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

	return decodeMergeRequestsInfo(resp)
}

func (g Gitlab) getMRWithDiffs(iid int, resChan chan MergeRequestTransfer, isChannelClosed func() bool) {
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
		return
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
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// TODO: подумоть
			log.Error(err)
		}
	}()
	mrWithFileChanges, err := decodeSingleMergeRequestItem(resp)
	if err != nil {
		if isChannelClosed() {
			return
		}
		resChan <- MergeRequestTransfer{
			mergeRequest: nil,
			error:        err,
		}
		return
	}
	log.WithFields(log.Fields{"iid": iid, "mrWithFileChanges": mrWithFileChanges}).Info("мр успешно получен")
	if isChannelClosed() {
		return
	}
	resChan <- MergeRequestTransfer{
		mergeRequest: mrWithFileChanges,
		error:        nil,
	}
}
