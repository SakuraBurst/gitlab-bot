package logger

import (
	"github.com/SakuraBurst/gitlab-bot/pkg/telegram"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

type FixedJSONFormatter struct {
	OriginalLogger log.JSONFormatter
}

func (f *FixedJSONFormatter) Format(entry *log.Entry) ([]byte, error) {
	bytes, err := f.OriginalLogger.Format(entry)
	if err != nil {
		return bytes, err
	}
	bytes[len(bytes)-1] = ','
	bytes = append(bytes, '\n')
	return bytes, err
}

func Init(level log.Level, output io.Writer) {
	log.SetReportCaller(true)
	log.SetFormatter(&FixedJSONFormatter{log.JSONFormatter{PrettyPrint: true}})
	log.SetOutput(output)
	log.Info("Логер инициализирован")
	log.SetLevel(level)
}

func AddHook(hook log.Hook) {
	log.AddHook(hook)
}

type FatalNotifier struct {
	Bot     telegram.Bot
	LogFile *os.File
}

func (f *FatalNotifier) Levels() []log.Level {
	return []log.Level{log.FatalLevel}
}

func (f *FatalNotifier) Fire(entry *log.Entry) error {
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	err := f.Bot.SendMessage(entry.Message)
	if err != nil {
		panic(err)
	}
	return nil
}
