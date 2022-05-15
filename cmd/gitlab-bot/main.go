package main

import (
	"github.com/SakuraBurst/gitlab-bot/internal/helpers"
	"github.com/SakuraBurst/gitlab-bot/internal/logger"
	"github.com/SakuraBurst/gitlab-bot/internal/worker"
	"github.com/SakuraBurst/gitlab-bot/pkg/basa_dannih"
	"github.com/SakuraBurst/gitlab-bot/pkg/gitlab"
	"github.com/SakuraBurst/gitlab-bot/pkg/telegram"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

const True = "true"

var envMap map[string]string
var bd = make(basa_dannih.BasaDannihMySQLPostgresMongoPgAdmin777)
var neededEnv = []string{"VIEW_CHANGES", "PROJECT", "GITLAB_TOKEN", "TELEGRAM_CHANEL", "TELEGRAM_BOT_TOKEN", "FATAL_REMINDER"}

func init() {
	absLoggerFilePath, _ := filepath.Abs("../gitlab-bot/internal/logger/logger.log")
	f, err := os.OpenFile(absLoggerFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	absDotEnvPath, _ := filepath.Abs("../gitlab-bot/.env")
	logger.Init(log.InfoLevel, f)
	err = godotenv.Load(absDotEnvPath)
	if err != nil {
		log.Fatal(err)
	}
	err = helpers.CheckForEnv(neededEnv)
	if err != nil {
		log.Fatal(err)
	}
	envMap = helpers.GetEnvMap()
	reminderBot := telegram.NewBot(envMap["TELEGRAM_BOT_TOKEN"], envMap["FATAL_REMINDER"])
	logger.AddHook(&logger.FatalNotifier{
		Bot:     reminderBot,
		LogFile: f,
	})
	log.WithFields(log.Fields{
		"with diffs":         envMap["VIEW_CHANGES"],
		"project":            envMap["PROJECT"],
		"gitlab token":       envMap["GITLAB_TOKEN"],
		"telegram chanel":    envMap["TELEGRAM_CHANEL"],
		"telegram bot token": envMap["TELEGRAM_BOT_TOKEN"],
	}).Info("Проект инициализирован")
}

func main() {
	if envMap["WITHOUT_NOTIFIER"] == True && envMap["WITHOUT_REMINDER"] == True {
		log.SetOutput(os.Stderr)
		log.Fatal("ну и чего ты ожидал? Без объявлялки и напоминалки это бот ничего не умеет делать")
	}

	git := gitlab.NewGitlabConn(envMap["VIEW_CHANGES"] == True, envMap["PROJECT"], envMap["GITLAB_TOKEN"], "https://gitlab.innostage-group.ru")
	tlBot := telegram.NewBot(envMap["TELEGRAM_BOT_TOKEN"], envMap["TELEGRAM_CHANEL"])

	mergeRequests, err := git.MergeRequests()
	if err != nil {
		log.Fatal(err)
	}
	helpers.WriteMrsToBd(bd, mergeRequests.MergeRequests...)

	stop := make(chan error)
	if envMap["WITHOUT_NOTIFIER"] != True {
		go worker.WaitFor24Hours(stop, git, tlBot)
	}
	if envMap["WITHOUT_REMINDER"] != True {
		go worker.WaitForMinute(stop, git, tlBot, bd)
	}

	err = <-stop
	log.Fatal(err)
}
