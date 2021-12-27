package telegram

func (t TelegramBot) SendInitMessage(ver string) {
	t.sendMessage("Бот запущен, версия " + ver)
}
