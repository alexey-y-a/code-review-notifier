package repository

import (
	"database/sql"

	"github.com/alexey-y-a/code-review-notifier/internal/model"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type UserRepository interface {
	SaveUser(user *model.User) error
	GetUserByGitHubLogin(login string) (*model.User, error)
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(dsn string) (*PostgresUserRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		zap.L().Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}
	if err := db.Ping(); err != nil {
		zap.L().Error("Failed to ping database", zap.Error(err))
		return nil, err
	}
	return &PostgresUserRepository{db: db}, nil
}

func (r *PostgresUserRepository) SaveUser(user *model.User) error {
	query := `INSERT INTO users (telegram_id, github_login) VALUES ($1, $2) ON CONFLICT (telegram_id) DO UPDATE SET github_login = $2`
	_, err := r.db.Exec(query, user.TelegramID, user.GitHubLogin)
	if err != nil {
		zap.L().Error("Failed to save user", zap.Error(err))
	}
	return err
}

func (r *PostgresUserRepository) GetUserByGitHubLogin(login string) (*model.User, error) {
	query := `SELECT telegram_id, github_login FROM users WHERE github_login = $1`
	var user model.User
	err := r.db.QueryRow(query, login).Scan(&user.TelegramID, &user.GitHubLogin)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		zap.L().Error("Failed to get user", zap.Error(err))
		return nil, err
	}
	return &user, nil
}
