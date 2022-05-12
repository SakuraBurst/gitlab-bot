package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/SakuraBurst/gitlab-bot/internal/helpers"
	"github.com/SakuraBurst/gitlab-bot/internal/logger"
	"github.com/SakuraBurst/gitlab-bot/internal/worker"
	"github.com/SakuraBurst/gitlab-bot/pkg/gitlab"
	"github.com/SakuraBurst/gitlab-bot/pkg/models"
	"github.com/SakuraBurst/gitlab-bot/pkg/telegram"
	"github.com/joho/godotenv"
	"net"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

const True = "true"

var envMap map[string]string
var bd = make(models.BasaDannihMySQLPostgresMongoPgAdmin777)

func init() {
	var err error
	envMap, err = godotenv.Read("../../.env")
	if err != nil || len(envMap) == 0 {
		envMap = helpers.GetOsEnvMap()
	}

	http.DefaultClient.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			// UNSAFE!
			// DON'T USE IN PRODUCTION!
			InsecureSkipVerify: true,
		},
	}
	logger.Init()
	log.WithFields(log.Fields{
		"with diffs":         envMap["VIEW_CHANGES"],
		"project":            envMap["PROJECT"],
		"gitlab token":       envMap["GITLAB_TOKEN"],
		"telegram chanel":    envMap["TELEGRAM_CHANEL"],
		"telegram bot token": envMap["TELEGRAM_BOT_TOKEN"],
	}).Info("Проект инициализирован")
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
	testBytes, err := json.Marshal(tgRequest)
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(testBytes)
	response, err := http.Post("https://api.telegram.org/bot"+f.Token+"/sendMessage", "application/json", reader)
	if err != nil {
		panic(err)
	}
	response.Body.Close()
	return nil
}

func main() {

	if envMap["WITHOUT_NOTIFIER"] == True && envMap["WITHOUT_REMINDER"] == True {
		log.SetOutput(os.Stderr)
		log.Fatal("ну и чего ты ожидал? Без объявлялки и напоминалки это бот ничего не умеет делать")
	}

	git := gitlab.NewGitlabConn(envMap["VIEW_CHANGES"] == True, envMap["PROJECT"], envMap["GITLAB_TOKEN"], "https://gitlab.innostage-group.ru")
	tlBot := telegram.NewBot(envMap["TELEGRAM_BOT_TOKEN"], envMap["TELEGRAM_CHANEL"])

	logger.AddHook(&FatalReminderChannel{
		envMap["FATAL_REMINDER"], envMap["TELEGRAM_BOT_TOKEN"],
	})

	//mergeRequests, err := git.MergeRequests()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//helpers.WriteMrsToBd(bd, mergeRequests.MergeRequests...)

	stop := make(chan bool)
	if envMap["WITHOUT_NOTIFIER"] != True {
		go worker.WaitFor24Hours(stop, git, tlBot)
	}
	if envMap["WITHOUT_REMINDER"] != True {
		go worker.WaitForMinute(stop, git, tlBot, bd)
	}

	if <-stop {
		log.Fatal("что-то пошло не так")
	}
	//

}
