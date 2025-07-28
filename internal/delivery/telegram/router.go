package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func StartBot(bot *tgbotapi.BotAPI, handler *Handler) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		go handler.HandleUpdate(update)
	}
}
