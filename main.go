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
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var withDiffs = false
var project = ""
var token = ""

func init() {
	godotenv.Load()
	withDiffs = os.Getenv("VIEW_CHANGES") == "true"
	project = os.Getenv("PROJECT")
	token = os.Getenv("TOKEN")
}

func main() {
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
	logger.LoggerInit()
	log.WithFields(log.Fields{"withDiffs": withDiffs, "project": project, "token": token}).Info("Проект инициализирован")
	Wait()
	//

}

func Wait() {
	t := time.Now()
	if t.Hour() != 12 {
		waitFor := 0
		if t.Hour() > 12 {
			waitFor = 24 - t.Hour() + 12
		} else {
			waitFor = 12 - t.Hour()
		}
		log.WithField("current time", t).Infof("sleep for %d hour(s)", waitFor)
		time.Sleep(time.Hour * time.Duration(waitFor))
		Wait()
	} else {
		mrWithDiffs := parser.Parser(project, token, withDiffs)
		telegram.SendMessage(mrWithDiffs, withDiffs)
		log.WithField("текущее время", t).Info("sleep for 24 hours")
		time.Sleep(time.Hour * 24)
		Wait()
	}
}
