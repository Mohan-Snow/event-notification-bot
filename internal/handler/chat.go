package handler

import (
	"context"
)

type ChatService interface {
	AddChat(ctx context.Context, id int64) error
	DeleteChat(ctx context.Context, id int64) error
}

type ChatHandler struct {
	chatService ChatService
}

func NewChatHandler(chatService ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

func (h *ChatHandler) StartChat(ctx context.Context, chatID int64) error {
	return h.chatService.AddChat(ctx, chatID)
}

func (h *ChatHandler) StopChat(ctx context.Context, chatID int64) error {
	return h.chatService.DeleteChat(ctx, chatID)
}
