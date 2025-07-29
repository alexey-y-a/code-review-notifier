package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func StartBot(bot *tgbotapi.BotAPI, handler *Handler) {
	u := tgbotapi.NewUpdate(0)
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		handler.HandleUpdate(update)
	}
}
