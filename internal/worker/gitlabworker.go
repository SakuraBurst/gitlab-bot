package worker

import (
	"github.com/SakuraBurst/gitlab-bot/internal/helpers"
	"github.com/SakuraBurst/gitlab-bot/pkg/basa_dannih"
	"github.com/SakuraBurst/gitlab-bot/pkg/gitlab"
	"github.com/SakuraBurst/gitlab-bot/pkg/telegram"
	"time"

	log "github.com/sirupsen/logrus"
)

const halfDay = 12

const fullDay = 24

var errorCounter = 0

func WaitFor24Hours(stop chan<- error, glConn gitlab.Gitlab, tlBot telegram.Bot) {
	for {
		t := time.Now()

		if t.Hour() != halfDay {
			waitFor := 0
			if t.Hour() > halfDay {
				waitFor = fullDay - t.Hour() + halfDay
			} else {
				waitFor = halfDay - t.Hour()
			}
			log.Infof("sleep for %d hour(s)", waitFor)
			time.Sleep(time.Hour * time.Duration(waitFor))
			continue
		} else {
			if t.Weekday() == 0 || t.Weekday() == 6 {
				time.Sleep(time.Hour * fullDay)
				continue
			}
			mergeRequests, err := glConn.MergeRequests()
			errorChecker(err, stop)
			err = tlBot.SendMergeRequestMessage(mergeRequests, false, glConn.WithDiffs)
			errorChecker(err, stop)
			log.Info("sleep for 24 hours")
			time.Sleep(time.Hour * fullDay)
		}
	}
}

func WaitForMinute(stop chan error, glConn gitlab.Gitlab, tlBot telegram.Bot, bd basa_dannih.BasaDannihMySQLPostgresMongoPgAdmin777) {
	for {
		mergeRequests, err := glConn.MergeRequests()
		errorChecker(err, stop)
		openedMrs, writtenMrs, ok := helpers.OnlyNewMrs(mergeRequests.MergeRequests, bd)
		mergeRequests.MergeRequests = openedMrs
		mergeRequests.Length = writtenMrs
		log.WithFields(log.Fields{"Количество новых мрок": mergeRequests.Length, "Статус": ok}).Info("Ежеминутный обход")
		if ok {
			err := tlBot.SendMergeRequestMessage(mergeRequests, true, glConn.WithDiffs)
			errorChecker(err, stop)
		}
		log.Info("sleep for 1 minute")
		time.Sleep(time.Minute)
	}

}

func errorChecker(err error, stopChanel chan<- error) {
	if err != nil {
		log.Error("gg")
		errorCounter++

	}
	if errorCounter == 100 {
		stopChanel <- err
	}
}
