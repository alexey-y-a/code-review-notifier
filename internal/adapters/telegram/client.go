package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type TelegramClient interface {
	SendMessage(chatID int64, text string) error
}

type TelegramBotClient struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramBotClient(token string) (*TelegramBotClient, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		zap.L().Error("Failed to initialize Telegram bot", zap.Error(err))
		return nil, err
	}
	return &TelegramBotClient{bot: bot}, nil
}

func (c *TelegramBotClient) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := c.bot.Send(msg)
	if err != nil {
		zap.L().Error("Failed to send message", zap.Error(err))
	}
	return err
}
