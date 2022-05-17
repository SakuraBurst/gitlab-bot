package gitlab

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

const OPENED = "opened"

func (g Gitlab) getMergeRequestURL() (*url.URL, http.Header, error) {
	mergeRequestsURL, err := url.Parse(fmt.Sprintf("%s/api/v4/projects/%s/merge_requests", g.url, url.QueryEscape(g.repo)))
	if err != nil {
		return nil, nil, err
	}
	query := url.Values{}
	query.Set("state", OPENED)
	query.Set("with_merge_status_recheck", "true")
	mergeRequestsURL.RawQuery = query.Encode()
	headers := make(http.Header)
	headers.Add("PRIVATE-TOKEN", g.token)
	log.Info(mergeRequestsURL.String())
	return mergeRequestsURL, headers, err
}

func decodeMergeRequestsInfo(request *http.Response) (*MergeRequestsInfo, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	decoder := json.NewDecoder(request.Body)
	if request.StatusCode != http.StatusOK {
		var gitlabError models.GitlabError
		err := decoder.Decode(&gitlabError)
		if err != nil {
			return nil, err
		}
		return nil, gitlabError
	}
	mergeRequests := make([]models.MergeRequest, 0)
	err := decoder.Decode(&mergeRequests)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.WithFields(log.Fields{"Количество мрок": len(mergeRequests)}).Info("парсер получил список мрок")
	return &MergeRequestsInfo{
		Length:        len(mergeRequests),
		On:            time.Now(),
		MergeRequests: mergeRequests,
	}, err
}

func decodeSingleMergeRequestItem(request *http.Response) (*models.MergeRequest, error) {
	decoder := json.NewDecoder(request.Body)

	if request.StatusCode != http.StatusOK {
		var gitlabError models.GitlabError
		err := decoder.Decode(&gitlabError)
		if err != nil {
			return nil, err
		}
		return nil, gitlabError
	}

	mrListItem := &models.MergeRequest{}
	err := decoder.Decode(mrListItem)
	if err != nil {
		return nil, err
	}
	return mrListItem, err
}

func (g Gitlab) getSingleMergeRequestWithChangesURL(iid int) (*url.URL, http.Header, error) {
	mergeRequestURL, err := url.Parse(fmt.Sprintf("%s/api/v4/projects/%s/merge_requests/%d/changes", g.url, url.QueryEscape(g.repo), iid))
	if err != nil {
		return nil, nil, err
	}
	headers := make(http.Header)
	headers.Add("PRIVATE-TOKEN", g.token)
	return mergeRequestURL, headers, err
}
