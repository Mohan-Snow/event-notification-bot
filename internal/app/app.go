package app

import (
	"context"
	"errors"

	"event-notification-bot/internal/service/chat"
	"event-notification-bot/internal/service/scheduler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type ChatService interface {
	AddChat(ctx context.Context, id int64) error
	DeleteChat(ctx context.Context, id int64) error
}

type CatSender interface {
	SendPhoto(ctx context.Context) error
}

type App struct {
	bot         *tgbotapi.BotAPI
	updates     tgbotapi.UpdatesChannel
	logger      *zap.Logger
	chatService ChatService
	catSender   CatSender
}

func New(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, catSender CatSender, chatService ChatService, logger *zap.Logger) *App {
	return &App{
		bot:         bot,
		updates:     updates,
		logger:      logger,
		chatService: chatService,
		catSender:   catSender,
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

			// Добавляем chatID в список
			err := a.chatService.AddChat(ctx, update.Message.Chat.ID)
			if err != nil && !errors.Is(err, chat.ErrChatAlreadyExists) {
				a.logger.Error("Can't save chat", zap.Error(err))

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка")
				msg.ReplyToMessageID = update.Message.MessageID

				a.bot.Send(msg)
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я буду отправлять фотку котов каждые 10 сек")
			a.bot.Send(msg)
		}
	}
}
