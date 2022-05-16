package workers

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func FullDayWorker(stop chan<- error, work Work, wakeUpHour int) {
	for {
		t := time.Now()
		if t.Hour() != wakeUpHour {
			waitFor := 0
			if t.Hour() > wakeUpHour {
				waitFor = fullDay - t.Hour() + wakeUpHour
			} else {
				waitFor = wakeUpHour - t.Hour()
			}
			log.Infof("sleep for %d hour(s)", waitFor)
			time.Sleep(time.Hour * time.Duration(waitFor))
			continue
		} else {
			err := work()
			errorChecker(err, stop)
			log.Info("sleep for 24 hours")
			time.Sleep(time.Hour * fullDay)
		}
	}
}
