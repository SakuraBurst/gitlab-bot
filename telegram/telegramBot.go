package telegram

type TelegramBot struct {
	token       string
	mainChannel string
}

func NewBot(token, mainChannel string) TelegramBot {
	return TelegramBot{token: token, mainChannel: mainChannel}
}
