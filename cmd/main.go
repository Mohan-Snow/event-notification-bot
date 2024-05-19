package main

import (
	"database/sql"
	"event-notification-bot/config"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	// Set up application configs
	appConfig, err := config.NewConfig(logger)

	// establish data source connection
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		appConfig.DbHost, appConfig.DbPort, appConfig.DbUsername, appConfig.DbPassword, appConfig.DbName)
	db, err := sql.Open("postgres", psqlInfo) // to config
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
	bot.Debug = true // to config
	logger.Info("Authorized on telegram account", zap.String("Bot Name", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

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
