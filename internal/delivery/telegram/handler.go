package telegram

import (
	"github.com/alexey-y-a/code-review-notifier/internal/service"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strings"
)

type Handler struct {
	service service.Service
	bot     *tgbotapi.BotAPI
}

func NewHandler(service service.Service, bot *tgbotapi.BotAPI) *Handler {
	return &Handler{service: service, bot: bot}
}

func (h *Handler) HandleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	text := update.Message.Text

	if strings.HasPrefix(text, "/register ") {
		githubLogin := strings.TrimSpace(strings.TrimPrefix(text, "/register "))
		if githubLogin == "" {
			h.sendResponse(chatID, "Пожалуйста, укажите ваш GitHub логин: /register <ваш_логин>")
			return
		}
		if err := h.service.RegisterUser(chatID, githubLogin); err != nil {
			h.sendResponse(chatID, "Ошибка при регистрации: "+err.Error())
			return
		}
		h.sendResponse(chatID, "Вы успешно зарегистрированы с GitHub логином: "+githubLogin)
	} else {
		h.sendResponse(chatID, "Используйте команду /register <ваш_логин> для регистрации")
	}
}

func (h *Handler) sendResponse(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := h.bot.Send(msg); err != nil {
		zap.L().Error("Failed to send response", zap.Error(err))
	}
}
