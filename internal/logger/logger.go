package logger

import (
	"github.com/SakuraBurst/gitlab-bot/pkg/services/telegram"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
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

func GetLogLevel(isProduction bool) log.Level {
	if isProduction {
		return log.ErrorLevel
	}
	return log.InfoLevel
}

type FatalNotifier struct {
	Bot     telegram.Bot
	LogFile *os.File
}

func (f *FatalNotifier) Levels() []log.Level {
	return []log.Level{log.FatalLevel}
}

func (f *FatalNotifier) Fire(entry *log.Entry) error {
	err := f.Bot.SendMessage(entry.Message)
	if err != nil {
		panic(err)
	}
	absLoggerFilePath, err := filepath.Abs("logger.json")
	if err != nil {
		panic(err)
	}
	file, err := os.Open(absLoggerFilePath)
	defer func() {
		if err := file.Close(); err != nil {
			log.Info(err)
		}
	}()
	if err != nil {
		panic(err)
	}
	err = f.Bot.SendDocument(file, filepath.Base(file.Name()))
	if err != nil {
		panic(err)
	}
	return nil
}
