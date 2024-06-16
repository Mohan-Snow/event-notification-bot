package cat

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"event-notification-bot/internal/model"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type ChatService interface {
	ListChat(ctx context.Context) ([]*model.Chat, error)
}

type BotSender interface {
	SendMessageToAllChats(ctx context.Context, message string) error
}

type Service struct {
	chatService ChatService
	botSender   BotSender
	logger      *zap.Logger
}

func NewSender(chatService ChatService, botSender BotSender, logger *zap.Logger) *Service {
	return &Service{
		chatService: chatService,
		botSender:   botSender,
		logger:      logger,
	}
}

const (
	catAPIURL = "https://api.thecatapi.com/v1/images/search"
)

// SendPhoto Функция для отправки случайного фото кота
func (s *Service) SendPhoto(ctx context.Context) error {
	files, err := os.ReadDir("cats") // Папка, в которой хранятся фото котов
	if err != nil || len(files) == 0 {
		s.logger.Info("Local cat images not found or directory is empty, fetching from the internet")

		return s.sendPhotoFromInternet(ctx)
	}

	// Выбираем случайное фото
	randomIndex := rand.Intn(len(files))
	photo := files[randomIndex]

	filePath := "cats/" + photo.Name()

	return s.botSender.SendMessageToAllChats(ctx, filePath)
}

// sendPhotoFromInternet Функция для отправки случайного фото кота из интернета
func (s *Service) sendPhotoFromInternet(ctx context.Context) error {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(catAPIURL)
	res := fasthttp.AcquireResponse()

	client := &fasthttp.Client{}
	if err := client.Do(req, res); err != nil {
		s.logger.Info("Error fetching image from the internet")
	}
	fasthttp.ReleaseRequest(req)

	var result []struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(res.Body(), &result); err != nil {
		return err
	}
	fasthttp.ReleaseResponse(res)

	if len(result) == 0 {
		return fmt.Errorf("No image URL found in response")
	}

	imageURL := result[0].URL

	return s.botSender.SendMessageToAllChats(ctx, imageURL)
}
