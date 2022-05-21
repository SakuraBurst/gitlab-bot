package main

import (
	"github.com/SakuraBurst/gitlab-bot/internal/helpers"
	"github.com/SakuraBurst/gitlab-bot/internal/logger"
	"github.com/SakuraBurst/gitlab-bot/internal/templates"
	"github.com/SakuraBurst/gitlab-bot/internal/workers"
	"github.com/SakuraBurst/gitlab-bot/pkg/basa_dannih"
	"github.com/SakuraBurst/gitlab-bot/pkg/services/gitlab"
	"github.com/SakuraBurst/gitlab-bot/pkg/services/telegram"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

const True = "true"

var envMap map[string]string
var bd = make(basa_dannih.BasaDannihMySQLPostgresMongoPgAdmin777)
var neededEnv = []string{"PROJECT", "GITLAB_TOKEN", "TELEGRAM_CHANEL", "TELEGRAM_BOT_TOKEN", "FATAL_REMINDER"}

func init() {
	absLoggerFilePath, _ := filepath.Abs("../gitlab-bot/internal/logger/logger.json")
	f, err := os.OpenFile(absLoggerFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	absDotEnvPath, err := filepath.Abs("../gitlab-bot/.env")
	if err != nil {
		log.Fatal(err)
	}
	err = godotenv.Load(absDotEnvPath)
	if err != nil {
		log.Fatal(err)
	}
	err = helpers.CheckForEnv(neededEnv)
	if err != nil {
		log.Fatal(err)
	}
	envMap = helpers.GetEnvMap()
	logger.Init(logger.GetLogLevel(envMap["IS_PRODUCTION"] == True), f)
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

	var FullDayWork = func() error {
		t := time.Now()
		// отдыхаем
		if t.Weekday() == 0 || t.Weekday() == 6 {
			return nil
		}
		mergeRequests, err := git.MergeRequests()
		if err != nil {
			return err
		}
		log.WithFields(log.Fields{"Количество новых мрок": mergeRequests.Length}).Info("Ежедневный обход")
		openedMrTemplate := templates.GetRightTemplate(false, git.WithDiffs)
		message := templates.CreateStringFromTemplate(openedMrTemplate, mergeRequests)
		err = tlBot.SendMessage(message)
		return err
	}

	var OneMinuteWork = func() error {
		mergeRequests, err := git.MergeRequests()
		if err != nil {
			return err
		}
		openedMrs, writtenMrs, ok := helpers.OnlyNewMrs(mergeRequests.MergeRequests, bd)
		mergeRequests.MergeRequests = openedMrs
		mergeRequests.Length = writtenMrs
		log.WithFields(log.Fields{"Количество новых мрок": mergeRequests.Length, "Статус": ok}).Info("Ежеминутный обход")
		if ok {
			newMrTemplate := templates.GetRightTemplate(true, git.WithDiffs)
			message := templates.CreateStringFromTemplate(newMrTemplate, mergeRequests)
			err := tlBot.SendMessage(message)
			return err
		}
		return nil
	}

	stop := make(chan error)
	if envMap["WITHOUT_NOTIFIER"] != True {
		go workers.FullDayWorker(stop, FullDayWork, 12)
	}
	if envMap["WITHOUT_REMINDER"] != True {
		go workers.OneMinuteWorker(stop, OneMinuteWork)
	}

	err = <-stop
	log.Fatal(err)
}
