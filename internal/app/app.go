package app

import (
	"context"
	"errors"

	"event-notification-bot/internal/service/chat"
	"event-notification-bot/internal/service/scheduler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type ChatHandler interface {
	StartChat(ctx context.Context, chatID int64) error
	StopChat(ctx context.Context, chatID int64) error
}

type CatSender interface {
	SendPhoto(ctx context.Context) error
}

type App struct {
	bot         *tgbotapi.BotAPI
	updates     tgbotapi.UpdatesChannel
	logger      *zap.Logger
	catSender   CatSender
	chatHandler ChatHandler
}

func New(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, chatHandler ChatHandler, catSender CatSender, logger *zap.Logger) *App {
	return &App{
		bot:         bot,
		updates:     updates,
		logger:      logger,
		catSender:   catSender,
		chatHandler: chatHandler,
	}
}

func (a *App) Run(ctx context.Context) {
	// Устанавливаем шедулер для отправки фото каждые 10 сек
	err := scheduler.Start(func() {
		err := a.catSender.SendPhoto(ctx)
		if err != nil {
			a.logger.Error("Can't send photo", zap.Error(err))
		}
	})
	if err != nil {
		a.logger.Error("Can't start scheduler", zap.Error(err))
	}

	for update := range a.updates {
		if update.Message != nil { // If we got a message
			a.logger.Info("Incoming message", zap.String("User", update.Message.From.UserName),
				zap.String("Message", update.Message.Text))

			responseMessage, err := a.handleCommands(ctx, update.Message.Text, update.Message.Chat.ID)
			if err != nil {
				a.logger.Error("Error while handle commands", zap.Error(err))
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseMessage)
			if _, err := a.bot.Send(msg); err != nil {
				a.logger.Error("Error while sending message", zap.Error(err))
			}
		}
	}
}

func (a *App) handleCommands(ctx context.Context, command string, chatID int64) (string, error) {
	defaultMessage := "Я буду отправлять фотку котов каждые 10 сек"
	byeMessage := "Я больше не буду отправлять фотки котов. Пока"
	errMessage := "Произошла ошибка"

	switch command {
	case "/start":
		err := a.chatHandler.StartChat(ctx, chatID)
		if err != nil {
			if errors.Is(err, chat.ErrChatAlreadyExists) {
				return defaultMessage, nil
			}

			return errMessage, err
		}

		return defaultMessage, nil
	case "/stop":
		err := a.chatHandler.StopChat(ctx, chatID)
		if err != nil {
			return errMessage, err
		}

		return byeMessage, nil
	default:
		return defaultMessage, nil
	}
}
