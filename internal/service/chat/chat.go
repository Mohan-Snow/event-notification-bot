package chat

import (
	"context"
	"errors"

	"event-notification-bot/internal/model"
	"event-notification-bot/internal/repo/postgres"
)

type Saver interface {
	SaveChat(ctx context.Context, id int64) error
}

type Provider interface {
	ListChats(ctx context.Context) ([]*model.Chat, error)
	FindChatById(ctx context.Context, id int64) (*model.Chat, error)
}

type Deleter interface {
	DeleteChat(ctx context.Context, id int64) error
}

type Service struct {
	chatSaver    Saver
	chatProvider Provider
	chatDeleter  Deleter
}

var (
	ErrChatAlreadyExists = errors.New("chat already exists")
)

func NewService(chatSaver Saver, chatProvider Provider, chatDeleter Deleter) *Service {
	return &Service{
		chatSaver:    chatSaver,
		chatProvider: chatProvider,
		chatDeleter:  chatDeleter,
	}
}

func (s *Service) AddChat(ctx context.Context, id int64) error {
	isExists, err := s.checkIfChatExists(ctx, id)
	if err != nil {
		return err
	}
	if isExists {
		return ErrChatAlreadyExists
	}

	return s.chatSaver.SaveChat(ctx, id)
}

func (s *Service) DeleteChat(ctx context.Context, id int64) error {
	return s.chatDeleter.DeleteChat(ctx, id)
}

func (s *Service) ListChat(ctx context.Context) ([]*model.Chat, error) {
	return s.chatProvider.ListChats(ctx)
}

func (s *Service) checkIfChatExists(ctx context.Context, id int64) (bool, error) {
	_, err := s.chatProvider.FindChatById(ctx, id)
	if err != nil {
		if errors.Is(err, postgres.ErrChatNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
