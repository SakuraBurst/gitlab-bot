package gitlab

import (
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	log "github.com/sirupsen/logrus"
	"time"
)

type Gitlab struct {
	url       string
	repo      string
	token     string
	WithDiffs bool
}

type MergeRequestsInfo struct {
	Length        int
	On            time.Time
	MergeRequests []models.MergeRequest
}

type MergeRequestTransfer struct {
	mergeRequest *models.MergeRequest
	error        error
}

func NewGitlabConn(withDiffs bool, repo, token, url string) Gitlab {
	return Gitlab{repo: repo, token: token, WithDiffs: withDiffs, url: url}
}

func (g Gitlab) MergeRequests() (*MergeRequestsInfo, error) {
	log.WithFields(log.Fields{"repo": g.repo}).Info("парсер начал работу")
	openedMergeRequests, err := g.getAllOpenedMergeRequests()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if g.WithDiffs {
		return getAllOpenedMrsWithDiffs(g, openedMergeRequests)
	}
	log.WithFields(log.Fields{"Количество мрок со статусом opened": openedMergeRequests.Length}).Info("парсер закончил работу")
	return openedMergeRequests, nil
}

func getAllOpenedMrsWithDiffs(g Gitlab, mri *MergeRequestsInfo) (*MergeRequestsInfo, error) {
	responseWaiters := make(chan MergeRequestTransfer, mri.Length)
	closed := make(chan bool)
	var isClosed = func() bool {
		select {
		case <-closed:
			return true
		default:
			return false
		}
	}

	for _, v := range mri.MergeRequests {
		go g.getMRWithDiffs(v.Iid, responseWaiters, isClosed)
	}

	openedMergeRequestsWithDiffs := MergeRequestsInfo{
		Length:        mri.Length,
		On:            mri.On,
		MergeRequests: make([]models.MergeRequest, 0, mri.Length),
	}

	for i := 0; i < mri.Length; i++ {
		if !isClosed() {
			result := <-responseWaiters
			if result.error != nil {
				close(closed)
				close(responseWaiters)
				return nil, result.error
			}
			openedMergeRequestsWithDiffs.MergeRequests = append(openedMergeRequestsWithDiffs.MergeRequests, *result.mergeRequest)
		}
	}

	close(responseWaiters)
	return &openedMergeRequestsWithDiffs, nil
}
