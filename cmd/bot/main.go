package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alexey-y-a/code-review-notifier/internal/adapters/telegram"
	"github.com/alexey-y-a/code-review-notifier/internal/config"
	"github.com/alexey-y-a/code-review-notifier/internal/delivery/github"
	deliverytelegram "github.com/alexey-y-a/code-review-notifier/internal/delivery/telegram"
	"github.com/alexey-y-a/code-review-notifier/internal/repository"
	"github.com/alexey-y-a/code-review-notifier/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		zap.L().Fatal("Failed to load config", zap.Error(err))
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Password, cfg.Database.Name)
	repo, err := repository.NewPostgresUserRepository(dsn)
	if err != nil {
		zap.L().Fatal("Failed to initialize repository", zap.Error(err))
	}

	tgClient, err := telegram.NewTelegramBotClient(cfg.Telegram.Token)
	if err != nil {
		zap.L().Fatal("Failed to initialize Telegram client", zap.Error(err))
	}

	svc := service.NewBotService(repo, tgClient)

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		zap.L().Fatal("Failed to initialize Telegram bot", zap.Error(err))
	}
	tgHandler := deliverytelegram.NewHandler(svc, bot)
	go deliverytelegram.StartBot(bot, tgHandler)

	ghHandler := github.NewHandler(svc, cfg.GitHub.Secret)
	router := github.NewRouter(ghHandler)

	server := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: router,
	}
	zap.L().Info("Starting server", zap.String("port", cfg.Server.Port))
	if err := server.ListenAndServe(); err != nil {
		zap.L().Fatal("Failed to start server", zap.Error(err))
	}
}
