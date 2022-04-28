package telegram

func (t Bot) SendInitMessage(ver string) {
	t.SendMessage("Бот запущен, версия " + ver)
}
