package workers

import log "github.com/sirupsen/logrus"

const fullDay = 24

const maxPossibleErrors = 100

var errorCounter = 0

type Work func() error

func errorChecker(err error, stopChanel chan<- error) {
	if err != nil {
		log.Error(err)
		errorCounter++

	} else {
		errorCounter = 0
	}
	if errorCounter == maxPossibleErrors {
		stopChanel <- err
	}
}
