package main

import (
	"database/sql"
	"event-notification-bot/config"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	// Set up application configs
	appConfig, err := config.NewConfig()

	// establish data source connection
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		appConfig.DbHost, appConfig.DbPort, appConfig.DbUsername, appConfig.DbPassword, appConfig.DbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("Database initializing error")
		log.Fatal(err)
	} else {
		log.Println("Established connection to database postgres")
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Println("Database ping error")
		log.Fatal(err)
	}

	// Establish telegram api connection
	bot, err := tgbotapi.NewBotAPI(appConfig.TelegramToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}
