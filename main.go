package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/SakuraBurst/gitlab-bot/models"
)

const OPENED = "opened"

func main() {
	http.DefaultClient.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			// UNSAFE!
			// DON'T USE IN PRODUCTION!
			InsecureSkipVerify: true,
		},
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("https://gitlab.innostage-group.ru/api/v4/projects/%s/merge_requests", url.QueryEscape("gpe/ais-upu/ais-upu-frontend")), nil)
	// resp, err := http.Get("http://gitlab.innostage-group.ru/api/v4/projects/gpe%2Fais-upu%2Fais-upu-frontend/merge_requests")
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("PRIVATE-TOKEN", "ymrsGzzNEofRKhoX2f5G")

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
			go getMRDiffs(v.Iid, responseWaiter)
		}

	}
	mergeRequestsWithChanges := make([]models.MergeRequestFileChanges, 0, len(mergeRequests))
	for _, v := range responseWaiters {
		m := <-v
		mergeRequestsWithChanges = append(mergeRequestsWithChanges, m)
	}
	testRequest := map[string]string{
		"chat_id":    "@mrchicki",
		"text":       mergeRequestsWithChanges[0].Description,
		"parse_mode": "Markdown",
	}
	testBytes, err := json.Marshal(testRequest)
	if err != nil {
		log.Fatal(err)
	}
	reader := bytes.NewReader(testBytes)
	respon, err := http.Post("https://api.telegram.org/bot5021252898:AAFJr-XK1_pTKNEW3Ju7tvT-z1VOb75zycw/sendMessage", "application/json", reader)
	if err != nil {
		log.Fatal(err)
	}
	defer respon.Body.Close()
	decoder = json.NewDecoder(respon.Body)

	if respon.StatusCode != http.StatusOK {
		test := make(map[string]interface{})
		decoder.Decode(&test)
		fmt.Println(request.URL)
		log.Fatal(test)
	}
	fmt.Println(respon.StatusCode)

	// fmt.Println(mergeRequestsWithChanges)
}

func getMRDiffs(iid int, resChan chan models.MergeRequestFileChanges) {
	request, err := http.NewRequest("GET", fmt.Sprintf("http://gitlab.innostage-group.ru/api/v4/projects/%s/merge_requests/%d/changes", url.QueryEscape("gpe/ais-upu/ais-upu-frontend"), iid), nil)
	request.Header.Add("PRIVATE-TOKEN", "ymrsGzzNEofRKhoX2f5G")
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
