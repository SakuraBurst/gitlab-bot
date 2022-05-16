package workers

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func OneMinuteWorker(stop chan error, work Work) {
	for {
		err := work()
		errorChecker(err, stop)
		log.Info("sleep for 1 minute")
		time.Sleep(time.Minute)
	}
}
