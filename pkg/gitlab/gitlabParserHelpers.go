package gitlab

import (
	"encoding/json"
	"fmt"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

const OPENED = "opened"

func getMergeRequestURL(g Gitlab) (*url.URL, http.Header, error) {
	mergeRequestsURL, err := url.Parse(fmt.Sprintf("%s/api/v4/projects/%s/merge_requests", g.Url, url.QueryEscape(g.repo)))
	if err != nil {
		return nil, nil, err
	}
	query := url.Values{}
	query.Set("state", OPENED)
	mergeRequestsURL.RawQuery = query.Encode()
	headers := make(http.Header)
	headers.Add("PRIVATE-TOKEN", g.token)
	log.Info(mergeRequestsURL.String())
	return mergeRequestsURL, headers, err
}

func decodeMergeRequestsInfo(request *http.Response) (MergeRequestsInfo, error) {
	decoder := json.NewDecoder(request.Body)
	if request.StatusCode != http.StatusOK {
		var gitlabError models.GitlabError
		err := decoder.Decode(&gitlabError)
		if err != nil {
			return MergeRequestsInfo{}, err
		}
		return MergeRequestsInfo{}, gitlabError
	}
	mergeRequests := make([]models.MergeRequest, 0)
	err := decoder.Decode(&mergeRequests)
	if err != nil {
		log.Error(err)
		return MergeRequestsInfo{}, err
	}
	log.WithFields(log.Fields{"Количество мрок": len(mergeRequests)}).Info("парсер получил список мрок")
	return MergeRequestsInfo{
		Length:        len(mergeRequests),
		On:            time.Now(),
		MergeRequests: mergeRequests,
	}, err
}

func decodeSingleMergeRequestItem(request *http.Response) (models.MergeRequest, error) {
	decoder := json.NewDecoder(request.Body)

	if request.StatusCode != http.StatusOK {
		var gitlabError models.GitlabError
		err := decoder.Decode(&gitlabError)
		if err != nil {
			return models.MergeRequest{}, err
		}
		return models.MergeRequest{}, gitlabError
	}

	mrListItem := models.MergeRequest{}
	err := decoder.Decode(&mrListItem)
	if err != nil {
		log.Fatal(err)
	}
	return mrListItem, err
}

func getSingleMergeRequestWithChangesURL(g Gitlab, iid int) (*url.URL, http.Header, error) {
	mergeRequestURL, err := url.Parse(fmt.Sprintf("%s/api/v4/projects/%s/merge_requests/%d/changes", g.Url, url.QueryEscape(g.repo), iid))
	if err != nil {
		return nil, nil, err
	}
	headers := make(http.Header)
	headers.Add("PRIVATE-TOKEN", g.token)
	return mergeRequestURL, headers, err
}
