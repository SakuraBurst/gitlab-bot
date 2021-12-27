package telegram

func (t Bot) SendInitMessage(ver string) {
	t.sendMessage("Бот запущен, версия " + ver)
}
