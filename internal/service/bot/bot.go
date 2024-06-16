package bot

import (
	"context"

	"event-notification-bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ChatService interface {
	ListChat(ctx context.Context) ([]*model.Chat, error)
}

type Sender struct {
	bot         *tgbotapi.BotAPI
	chatService ChatService
}

func NewSender(bot *tgbotapi.BotAPI, chatService ChatService) *Sender {
	return &Sender{
		bot:         bot,
		chatService: chatService,
	}
}

func (s *Sender) SendMessageToAllChats(ctx context.Context, message string) error {
	chatList, err := s.chatService.ListChat(ctx)
	if err != nil {
		return err
	}

	for _, chat := range chatList {
		msg := tgbotapi.NewPhoto(chat.ChatID, tgbotapi.FileURL(message))
		_, err := s.bot.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}
