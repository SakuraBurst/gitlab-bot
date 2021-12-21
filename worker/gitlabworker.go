package worker

import (
	"time"

	"github.com/SakuraBurst/gitlab-bot/models"
	"github.com/SakuraBurst/gitlab-bot/parser"
	"github.com/SakuraBurst/gitlab-bot/telegram"
	log "github.com/sirupsen/logrus"
)

var BASA_DANNIH_MY_SQL_POSTGRES_MONGO_PG_ADMIN_777 = make(map[int]bool)

func WaitFor24Hours(withDiffs bool, stop chan bool, project, gitlabToken, telegramChanel, telegramBotToken string) {
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
			mergeRequests, err := parser.Parser(project, gitlabToken, withDiffs)
			if err != nil {
				log.Error("gg")
				stop <- true
			}
			telegram.SendMessage(mergeRequests, false, withDiffs, telegramChanel, telegramBotToken)
			log.Info("sleep for 24 hours")
			time.Sleep(time.Hour * 24)
		}
	}

}

func WaitForMinute(withDiffs bool, stop chan bool, project, gitlabToken, telegramChanel, telegramBotToken string) {
	for {
		log.Info("sleep for 1 minute")
		time.Sleep(time.Minute)
		mergeRequests, err := parser.Parser(project, gitlabToken, withDiffs)
		if err != nil {
			log.Error("gg")
			stop <- true
		}
		mergeRequests, ok := OnlyNewMrs(mergeRequests)
		log.WithFields(log.Fields{"Количество новых мрок": mergeRequests.Length, "Статус": ok}).Info("Ежеминутный обход")
		if ok {
			telegram.SendMessage(mergeRequests, true, withDiffs, telegramChanel, telegramBotToken)

		}
	}

}

func OnlyNewMrs(allOpendeMergeRequests models.MergeRequests) (models.MergeRequests, bool) {
	onlyNewMrs := models.MergeRequests{On: allOpendeMergeRequests.On}
	if allOpendeMergeRequests.Length == 0 {
		return onlyNewMrs, false
	}

	for _, v := range allOpendeMergeRequests.MergeRequests {
		if !BASA_DANNIH_MY_SQL_POSTGRES_MONGO_PG_ADMIN_777[v.Iid] {
			onlyNewMrs.MergeRequests = append(onlyNewMrs.MergeRequests, v)
			onlyNewMrs.Length++
			BASA_DANNIH_MY_SQL_POSTGRES_MONGO_PG_ADMIN_777[v.Iid] = true
			log.WithField("basa", BASA_DANNIH_MY_SQL_POSTGRES_MONGO_PG_ADMIN_777).Infof("база данных поплнена айдишником %d", v.Iid)
		}
	}
	return onlyNewMrs, onlyNewMrs.Length > 0
}
