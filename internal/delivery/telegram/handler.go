package telegram

import (
	"github.com/alexey-y-a/code-review-notifier/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Handler struct {
	service service.Service
	bot     *tgbotapi.BotAPI
}

func NewHandler(s service.Service, bot *tgbotapi.BotAPI) *Handler {
	return &Handler{service: s, bot: bot}
}

func (h *Handler) HandleUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		if update.Message.Text == "/start" {
			if err := h.service.RegisterUser(update.Message.Chat.ID, "test_user"); err != nil {
				zap.L().Error("Failed to register user", zap.Error(err))
			}
		}
	}
}
