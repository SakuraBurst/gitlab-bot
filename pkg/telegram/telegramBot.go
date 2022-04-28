package telegram

type Bot struct {
	token       string
	mainChannel string
}

func NewBot(token, mainChannel string) Bot {
	return Bot{token: token, mainChannel: mainChannel}
}
