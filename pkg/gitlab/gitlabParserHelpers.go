package gitlab

import (
	"encoding/json"
	"fmt"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"time"
)

const OPENED = "opened"

func getMergeRequest(g Gitlab) (*http.Request, error) {
	mergeRequestsURL, err := url.Parse(fmt.Sprintf("%s/api/v4/projects/%s/merge_requests", g.Url, url.QueryEscape(g.repo)))
	if err != nil {
		return nil, err
	}
	query := url.Values{}
	query.Set("state", OPENED)
	mergeRequestsURL.RawQuery = query.Encode()
	request, err := http.NewRequest("GET", mergeRequestsURL.String(), nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	request.Header.Add("PRIVATE-TOKEN", g.token)
	log.Info(request.URL.String())
	return request, err
}

func decodeMergeRequestsInfo(body io.Reader) (models.MergeRequestsInfo, error) {
	decoder := json.NewDecoder(body)
	mergeRequests := make([]models.MergeRequestListItem, 0)
	err := decoder.Decode(&mergeRequests)
	if err != nil {
		log.Error(err)
		return models.MergeRequestsInfo{}, err
	}
	log.WithFields(log.Fields{"Количество мрок": len(mergeRequests)}).Info("парсер получил список мрок")
	return models.MergeRequestsInfo{
		Length:        len(mergeRequests),
		On:            time.Now(),
		MergeRequests: mergeRequests,
	}, err
}

func decodeSingleMergeRequestItem(body io.Reader) (models.MergeRequestListItem, error) {
	decoder := json.NewDecoder(body)

	// TODO: замапать гитлаб эррор
	//if resp.StatusCode != http.StatusOK {
	//	test := make(map[string]interface{})
	//	err = decoder.Decode(&test)
	//	if err != nil {
	//		log.Panic(err)
	//	}
	//	log.WithFields(log.Fields{"url": request.URL}).Fatal(test)
	//}

	mrListItem := models.MergeRequestListItem{}
	err := decoder.Decode(&mrListItem)
	if err != nil {
		log.Fatal(err)
	}
	return mrListItem, err
}

func getSingleMergeRequestWithChanges(g Gitlab, iid int) (*http.Request, error) {
	mergeRequestURL, err := url.Parse(fmt.Sprintf("%s/api/v4/projects/%s/merge_requests/%d/changes", g.Url, url.QueryEscape(g.repo), iid))
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("GET", mergeRequestURL.String(), nil)

	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("PRIVATE-TOKEN", g.token)
	return request, err
}
