package repository

import (
	"database/sql"
	"testing"

	"github.com/alexey-y-a/code-review-notifier/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestPostgresUserRepository(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE users (telegram_id INTEGER PRIMARY KEY, github_login TEXT NOT NULL)`)
	assert.NoError(t, err)

	repo := &PostgresUserRepository{db: db}

	t.Run("SaveUser", func(t *testing.T) {
		user := &model.User{TelegramID: 123, GitHubLogin: "alexey-y-a"}
		err := repo.SaveUser(user)
		assert.NoError(t, err)

		var savedLogin string
		err = db.QueryRow("SELECT github_login FROM users WHERE telegram_id = ?", 123).Scan(&savedLogin)
		assert.NoError(t, err)
		assert.Equal(t, "alexey-y-a", savedLogin)
	})

	t.Run("GetUserByGitHubLogin", func(t *testing.T) {
		user, err := repo.GetUserByGitHubLogin("alexey-y-a")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, int64(123), user.TelegramID)
	})
}
