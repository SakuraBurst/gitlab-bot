package worker

import (
	"github.com/SakuraBurst/gitlab-bot/internal/helpers"
	"github.com/SakuraBurst/gitlab-bot/pkg/gitlab"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	"github.com/SakuraBurst/gitlab-bot/pkg/telegram"
	"time"

	log "github.com/sirupsen/logrus"
)

func WaitFor24Hours(stop chan bool, glConn gitlab.Gitlab, tlBot telegram.Bot) {
	errorCounter := 0
	for {
		t := time.Now()

		if t.Hour() != 12 {
			waitFor := 0
			if t.Hour() > 12 {
				waitFor = 24 - t.Hour() + 12
			} else {
				waitFor = 12 - t.Hour()
			}
			log.Infof("sleep for %d hour(s)", waitFor)
			time.Sleep(time.Hour * time.Duration(waitFor))
			continue
		} else {
			if t.Weekday() == 0 || t.Weekday() == 6 {
				time.Sleep(time.Hour * 24)
				continue
			}
			mergeRequests, err := glConn.MergeRequests()
			if err != nil {
				log.Error("gg")
				errorCounter++
			}
			if errorCounter == 100 {
				stop <- true
			}
			tlBot.SendMergeRequestMessage(mergeRequests, false, glConn.WithDiffs)
			log.Info("sleep for 24 hours")
			time.Sleep(time.Hour * 24)
		}
	}

}

func WaitForMinute(stop chan bool, glConn gitlab.Gitlab, tlBot telegram.Bot, bd models.BasaDannihMySQLPostgresMongoPgAdmin777) {
	errorCounter := 0
	for {
		mergeRequests, err := glConn.MergeRequests()
		if err != nil {
			log.Error("gg")
			errorCounter++
		}
		if errorCounter == 100 {
			stop <- true
		}
		mergeRequests, ok := helpers.OnlyNewMrs(mergeRequests, bd)
		log.WithFields(log.Fields{"Количество новых мрок": mergeRequests.Length, "Статус": ok}).Info("Ежеминутный обход")
		if ok {
			tlBot.SendMergeRequestMessage(mergeRequests, true, glConn.WithDiffs)
		}
		log.Info("sleep for 1 minute")
		time.Sleep(time.Minute)
	}

}
