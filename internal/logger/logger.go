package logger

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func Init() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{PrettyPrint: true})
	f, err := os.OpenFile("logger.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(f)
	} else {
		fmt.Println(err)
	}
	log.Info("Логер инициализирован")

	log.SetLevel(log.WarnLevel)
}

func AddHook(hook log.Hook) {
	log.AddHook(hook)
}
