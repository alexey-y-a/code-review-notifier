package service

import (
	"testing"

	"github.com/alexey-y-a/code-review-notifier/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) SaveUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepository) GetUserByGitHubLogin(login string) (*model.User, error) {
	args := m.Called(login)
	return args.Get(0).(*model.User), args.Error(1)
}

type MockTelegramClient struct {
	mock.Mock
}

func (m *MockTelegramClient) SendMessage(chatID int64, text string) error {
	args := m.Called(chatID, text)
	return args.Error(0)
}

func TestBotService(t *testing.T) {
	repo := &MockRepository{}
	tgClient := &MockTelegramClient{}
	svc := NewBotService(repo, tgClient)

	t.Run("RegisterUser", func(t *testing.T) {
		user := &model.User{TelegramID: 123, GitHubLogin: "alexey-y-a"}
		repo.On("SaveUser", user).Return(nil)

		err := svc.RegisterUser(123, "alexey-y-a")
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("HandleGitHubEvent", func(t *testing.T) {
		event := &model.PullRequestEvent{
			Action:     "assigned",
			Assignee:   "alexey-y-a",
			Title:      "Test PR",
			HTMLURL:    "https://github.com/test/repo/pull/1",
			Repository: "test/repo",
		}
		user := &model.User{TelegramID: 123, GitHubLogin: "alexey-y-a"}
		repo.On("GetUserByGitHubLogin", "alexey-y-a").Return(user, nil)
		tgClient.On("SendMessage", int64(123), mock.Anything).Return(nil)

		err := svc.HandleGitHubEvent(event)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
		tgClient.AssertExpectations(t)
	})

	t.Run("HandleGitHubEvent_NoUser", func(t *testing.T) {
		event := &model.PullRequestEvent{Assignee: "unknown"}
		repo.On("GetUserByGitHubLogin", "unknown").Return(nil, nil)

		err := svc.HandleGitHubEvent(event)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}
