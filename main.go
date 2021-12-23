package main

import (
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/SakuraBurst/gitlab-bot/logger"
	"github.com/SakuraBurst/gitlab-bot/parser"
	"github.com/SakuraBurst/gitlab-bot/telegram"
	"github.com/SakuraBurst/gitlab-bot/worker"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var withDiffs = false
var withoutReminder = false
var withoutNotifier = false
var project = ""
var gitlabToken = ""
var telegramChanel = ""
var telegramBotToken = ""

func init() {
	godotenv.Load()
	withDiffs = os.Getenv("VIEW_CHANGES") == "true"
	withoutReminder = os.Getenv("WITHOUT_REMINDER") == "true"
	withoutNotifier = os.Getenv("WITHOUT_NOTIFIER") == "true"
	project = os.Getenv("PROJECT")
	gitlabToken = os.Getenv("GITLAB_TOKEN")
	telegramChanel = os.Getenv("TELEGRAM_CHANEL")
	telegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")

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
}

func main() {
	logger.LoggerInit()
	log.WithFields(log.Fields{
		"with diffs":         withDiffs,
		"project":            project,
		"gitlab token":       gitlabToken,
		"telegram chanel":    telegramChanel,
		"telegram bot token": telegramBotToken,
	}).Info("Проект инициализирован")

	if withoutNotifier && withoutReminder {
		log.SetOutput(os.Stderr)
		log.Fatal("ну и чего ты ожидал? Без объявлялки и напоминалки это бот ничего не умеет делать")
	}
	log.Info("начат первый тестовый прогон")
	mergeRequests, err := parser.Parser(project, gitlabToken, withDiffs)
	if err != nil {
		log.Fatal(err)
	}
	mergeRequests, _ = worker.OnlyNewMrs(mergeRequests)
	telegram.SendMessage(mergeRequests, false, withDiffs, telegramChanel, telegramBotToken)
	stop := make(chan bool)
	if !withoutNotifier {
		go worker.WaitFor24Hours(withDiffs, stop, project, gitlabToken, telegramChanel, telegramBotToken)
	}
	if !withoutReminder {
		go worker.WaitForMinute(withDiffs, stop, project, gitlabToken, telegramChanel, telegramBotToken)
	}

	if <-stop {
		log.Fatal("что-то пошло не так")
	}
	//

}
