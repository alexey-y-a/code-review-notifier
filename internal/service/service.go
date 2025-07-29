package service

import (
	"fmt"

	"github.com/alexey-y-a/code-review-notifier/internal/adapters/telegram"
	"github.com/alexey-y-a/code-review-notifier/internal/model"
	"github.com/alexey-y-a/code-review-notifier/internal/repository"
	"go.uber.org/zap"
)

type Service interface {
	RegisterUser(telegramID int64, githubLogin string) error
	HandleGitHubEvent(event *model.PullRequestEvent) error
}

type BotService struct {
	repo     repository.UserRepository
	telegram telegram.TelegramClient
}

func NewBotService(repo repository.UserRepository, telegram telegram.TelegramClient) *BotService {
	return &BotService{repo: repo, telegram: telegram}
}

func (s *BotService) RegisterUser(telegramID int64, githubLogin string) error {
	user := &model.User{TelegramID: telegramID, GitHubLogin: githubLogin}
	if err := s.repo.SaveUser(user); err != nil {
		zap.L().Error("Failed to register user", zap.Error(err))
		return err
	}
	zap.L().Info("User registered", zap.Int64("telegram_id", telegramID), zap.String("github_login", githubLogin))
	return nil
}

func (s *BotService) HandleGitHubEvent(event *model.PullRequestEvent) error {
	user, err := s.repo.GetUserByGitHubLogin(event.Assignee)
	if err != nil {
		zap.L().Error("Failed to get user by GitHub login", zap.Error(err))
		return err
	}
	if user == nil {
		zap.L().Warn("No user found for GitHub login", zap.String("github_login", event.Assignee))
		return nil
	}

	message := fmt.Sprintf("На вас назначен pull request: %s — %s (%s)", event.Title, event.HTMLURL, event.Repository)
	if err := s.telegram.SendMessage(user.TelegramID, message); err != nil {
		zap.L().Error("Failed to send Telegram message", zap.Error(err))
		return err
	}
	zap.L().Info("Notification sent", zap.Int64("telegram_id", user.TelegramID), zap.String("message", message))
	return nil
}
