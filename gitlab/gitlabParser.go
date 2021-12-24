package gitlab

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

func (g Gitlab) Parser() (models.MergeRequests, error) {
	log.WithFields(log.Fields{"repo": g.repo}).Info("парсер начал работу")
	request, err := http.NewRequest("GET", fmt.Sprintf("https://gitlab.innostage-group.ru/api/v4/projects/%s/merge_requests", url.QueryEscape(g.repo)), nil)
	if err != nil {
		log.Error(err)
		return models.MergeRequests{}, err
	}
	request.Header.Add("PRIVATE-TOKEN", g.token)

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

	log.WithFields(log.Fields{"Количество мрок": len(mergeRequests)}).Info("парсер получил список мрок")
	openedMergeRequests := models.MergeRequests{Length: 0, On: time.Now(), MergeRequests: make([]models.MergeRequestFileChanges, 0, len(mergeRequests))}
	for _, v := range mergeRequests {
		if v.State == OPENED {
			openedMergeRequests.Length++
			openedMergeRequests.MergeRequests = append(openedMergeRequests.MergeRequests, models.MergeRequestFileChanges{MergeRequestListItem: v, Changes: []models.FileChanges{}})
		}

	}
	if g.WithDiffs {
		responseWaiters := make(chan models.MergeRequestFileChanges, openedMergeRequests.Length)
		for _, v := range openedMergeRequests.MergeRequests {
			go g.getMRDiffs(v.Iid, responseWaiters)
		}

		openedMergeRequests = models.MergeRequests{Length: openedMergeRequests.Length, On: openedMergeRequests.On, MergeRequests: make([]models.MergeRequestFileChanges, 0, openedMergeRequests.Length)}

		for i := 0; i < openedMergeRequests.Length; i++ {
			openedMergeRequests.MergeRequests = append(openedMergeRequests.MergeRequests, <-responseWaiters)
		}

		close(responseWaiters)
	}
	log.WithFields(log.Fields{"Количество мрок со статусом opened": openedMergeRequests.Length}).Info("парсер закончил работу")
	return openedMergeRequests, nil
}

func (g Gitlab) getMRDiffs(iid int, resChan chan models.MergeRequestFileChanges) {
	log.WithFields(log.Fields{"iid": iid}).Info("получение отдельного открытого мр с доп даннымми")
	request, err := http.NewRequest("GET", fmt.Sprintf("http://gitlab.innostage-group.ru/api/v4/projects/%s/merge_requests/%d/changes", url.QueryEscape(g.repo), iid), nil)
	request.Header.Add("PRIVATE-TOKEN", g.token)
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
