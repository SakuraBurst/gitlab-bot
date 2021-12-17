package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func LoggerInit() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{PrettyPrint: true})
	f, err := os.OpenFile("logger.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(f)
	}
	log.Info("Логер инициализирован")

}
