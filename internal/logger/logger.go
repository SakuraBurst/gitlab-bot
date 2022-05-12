package logger

import (
	"github.com/SakuraBurst/gitlab-bot/api/clients"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
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

type FatalReminderChannel struct {
	Chat  string
	Token string
}

func (f *FatalReminderChannel) Levels() []log.Level {
	return []log.Level{log.FatalLevel}
}

func (f *FatalReminderChannel) Fire(entry *log.Entry) error {
	tgRequest := map[string]string{
		"chat_id":    f.Chat,
		"text":       entry.Message,
		"parse_mode": "html",
	}
	headers := make(http.Header)
	headers.Set("Content-Type", "application/json")
	response, err := clients.Post("https://api.telegram.org/bot"+f.Token+"/sendMessage", tgRequest, headers)
	if err != nil {
		panic(err)
	}
	response.Body.Close()
	return nil
}
