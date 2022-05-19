package telegram

type Bot struct {
	token       string
	mainChannel string
}

const telegramApi = "https://api.telegram.org"

func NewBot(token, mainChannel string) Bot {
	return Bot{token: token, mainChannel: mainChannel}
}
