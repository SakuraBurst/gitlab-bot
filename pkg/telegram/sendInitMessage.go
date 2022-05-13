package telegram

func (t Bot) SendInitMessage(ver string) error {
	return t.SendMessage("Бот запущен, версия " + ver)
}
