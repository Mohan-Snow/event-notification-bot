package main

import (
	"database/sql"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"event-notification-bot/internal/config"
)

func main() {
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
	connectionString := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		appConfig.DbHost, appConfig.DbPort, appConfig.DbUsername, appConfig.DbPassword, appConfig.DbName)
	db, err := sql.Open(appConfig.DataSource, connectionString)
	if err != nil {
		logger.Panic("Database initializing error", zap.Error(err))
	} else {
		logger.Info("Established connection to database postgres")
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		logger.Error("Database ping error", zap.Error(err))
	}

	// Establish telegram api connection
	bot, err := tgbotapi.NewBotAPI(appConfig.TelegramToken)
	if err != nil {
		logger.Panic("Establishing connection to telegram failed", zap.Error(err))
	}
	bot.Debug = appConfig.TelegramBotDebug
	logger.Info("Authorized on telegram account", zap.String("Bot Name", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = appConfig.TelegramUpdateTimeout

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			logger.Info("Incoming message", zap.String("User", update.Message.From.UserName),
				zap.String("Message", update.Message.Text))

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

