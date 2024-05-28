package api

import (
	"strconv"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type TelegramApiConfig struct {
	TelegramToken         string `required:"true" envconfig:"TELEGRAM_TOKEN"`
	TelegramBotDebug      bool   `envconfig:"TELEGRAM_BOT_DEBUG" default:"false"`
	TelegramUpdateTimeout int    `envconfig:"TELEGRAM_UPDATE_TIMEOUT" default:"60"`
}

// NewTelegramApiConfig Loads environmental variables from .env file and populates TelegramApiConfig struct with found values
func NewTelegramApiConfig(logger *zap.Logger) (*TelegramApiConfig, error) {
	envPath := "internal/config/.env"
	if err := godotenv.Load(envPath); err != nil {
		logger.Error("No .env file found for path", zap.String("path", envPath), zap.Error(err))
	}
	cfg := TelegramApiConfig{}
	// Process method tries to search environmental variable first by prefix+envconfig
	// if value isn't found it uses specified envconfig
	if err := envconfig.Process("", &cfg); err != nil {
		logger.Error("Error occurred while processing configs", zap.Error(err))
		return nil, err
	}
	logger.Info("Telegram API configuration.", zap.String("Telegram token", cfg.TelegramToken),
		zap.String("Telegram bot debug", strconv.FormatBool(cfg.TelegramBotDebug)),
		zap.String("Telegram update timeout", strconv.Itoa(cfg.TelegramUpdateTimeout)))
	return &cfg, nil
}
