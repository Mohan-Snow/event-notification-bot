package main

import (
	"context"
	"log"

	"event-notification-bot/internal/app"
	"event-notification-bot/internal/config"
	"event-notification-bot/internal/repo/postgres"
	"event-notification-bot/internal/service/bot"
	"event-notification-bot/internal/service/cat"
	"event-notification-bot/internal/service/chat"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	// Set up application configs
	appConfig, err := config.NewConfig(logger)
	if err != nil {
		logger.Error("Application configuring error", zap.Error(err))
	}

	// Establish data source connection
	db, err := postgres.New(ctx, appConfig.DBConfig)
	if err != nil {
		logger.Panic("Database initializing error", zap.Error(err))
	} else {
		logger.Info("Established connection to database postgres")
	}
	defer db.Close()

	// Establish telegram api connection
	tgbot, err := tgbotapi.NewBotAPI(appConfig.TelegramToken)
	if err != nil {
		logger.Panic("Establishing connection to telegram failed", zap.Error(err))
	}
	tgbot.Debug = appConfig.TelegramBotDebug
	logger.Info("Authorized on telegram account", zap.String("Bot Name", tgbot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = appConfig.TelegramUpdateTimeout

	updates := tgbot.GetUpdatesChan(u)

	chatService := chat.NewService(db, db, db)
	botSender := bot.NewSender(tgbot, chatService)
	catSender := cat.NewSender(chatService, botSender, logger)

	tgApp := app.New(tgbot, updates, catSender, chatService, logger)
	tgApp.Run(ctx)
}
