package parser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/SakuraBurst/gitlab-bot/models"
	log "github.com/sirupsen/logrus"
)

const OPENED = "opened"

func Parser(repo string, token string, withDiffs bool) (models.MergeRequests, error) {
	log.WithFields(log.Fields{"repo": repo}).Info("парсер начал работу")
	request, err := http.NewRequest("GET", fmt.Sprintf("https://gitlab.innostage-group.ru/api/v4/projects/%s/merge_requests", url.QueryEscape(repo)), nil)
	if err != nil {
		log.Error(err)
		return models.MergeRequests{}, err
	}
	request.Header.Add("PRIVATE-TOKEN", token)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Error(err)
		return models.MergeRequests{}, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	mergeRequests := make([]models.MergeRequestListItem, 0)
	err = decoder.Decode(&mergeRequests)
	if err != nil {
		log.Error(err)
		return models.MergeRequests{}, err
	}
	responseWaiters := make([]chan models.MergeRequestFileChanges, 0, len(mergeRequests))
	log.WithFields(log.Fields{"Количество мрок": len(mergeRequests)}).Info("парсер получил список мрок")
	openedMergeRequests := models.MergeRequests{Length: 0, On: time.Now(), MergeRequests: make([]models.MergeRequestFileChanges, 0, len(mergeRequests))}
	for _, v := range mergeRequests {
		if v.State == OPENED {
			openedMergeRequests.Length++
			if withDiffs {
				responseWaiter := make(chan models.MergeRequestFileChanges)
				responseWaiters = append(responseWaiters, responseWaiter)
				go getMRDiffs(v.Iid, responseWaiter, repo, token)
			} else {
				openedMergeRequests.MergeRequests = append(openedMergeRequests.MergeRequests, models.MergeRequestFileChanges{MergeRequestListItem: v, Changes: []models.FileChanges{}})

			}

		}

	}
	if withDiffs {
		for _, v := range responseWaiters {
			openedMergeRequests.MergeRequests = append(openedMergeRequests.MergeRequests, <-v)
		}
	}
	log.WithFields(log.Fields{"Количество мрок со статусом opened": openedMergeRequests.Length}).Info("парсер закончил работу")
	return openedMergeRequests, nil
}

func getMRDiffs(iid int, resChan chan models.MergeRequestFileChanges, repo, token string) {
	log.WithFields(log.Fields{"iid": iid}).Info("получение отдельного открытого мр с доп даннымми")
	request, err := http.NewRequest("GET", fmt.Sprintf("http://gitlab.innostage-group.ru/api/v4/projects/%s/merge_requests/%d/changes", url.QueryEscape(repo), iid), nil)
	request.Header.Add("PRIVATE-TOKEN", token)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		test := make(map[string]interface{})
		decoder.Decode(&test)
		log.WithFields(log.Fields{"url": request.URL}).Fatal(test)
	}

	mrWithFileChanges := models.MergeRequestFileChanges{}
	err = decoder.Decode(&mrWithFileChanges)
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{"iid": iid, "mrWithFileChanges": mrWithFileChanges}).Info("мр успешно получен")
	resChan <- mrWithFileChanges
}
