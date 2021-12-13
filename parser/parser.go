package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/SakuraBurst/gitlab-bot/models"
)

const OPENED = "opened"

func Parser(repo string, token string) models.MergeRequests {
	request, err := http.NewRequest("GET", fmt.Sprintf("https://gitlab.innostage-group.ru/api/v4/projects/%s/merge_requests", url.QueryEscape(repo)), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("PRIVATE-TOKEN", token)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	mergeRequests := make([]models.MergeRequestListItem, 0)
	err = decoder.Decode(&mergeRequests)
	if err != nil {
		log.Fatal(err)
	}
	responseWaiters := make([]chan models.MergeRequestFileChanges, 0, len(mergeRequests))
	for _, v := range mergeRequests {
		if v.State == OPENED {
			responseWaiter := make(chan models.MergeRequestFileChanges)
			responseWaiters = append(responseWaiters, responseWaiter)
			go getMRDiffs(v.Iid, responseWaiter, repo, token)
		}

	}
	mergeRequestsWithChanges := models.MergeRequests{Length: len(responseWaiters), On: time.Now(), MergeRequests: make([]models.MergeRequestFileChanges, 0, len(mergeRequests))}
	for _, v := range responseWaiters {
		mergeRequestsWithChanges.MergeRequests = append(mergeRequestsWithChanges.MergeRequests, <-v)
	}
	return mergeRequestsWithChanges
}

func getMRDiffs(iid int, resChan chan models.MergeRequestFileChanges, repo, token string) {
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
		fmt.Println(request.URL)
		log.Fatal(test)
	}

	mrWithFileChanges := models.MergeRequestFileChanges{}
	err = decoder.Decode(&mrWithFileChanges)
	if err != nil {
		log.Fatal(err)
	}
	resChan <- mrWithFileChanges
}