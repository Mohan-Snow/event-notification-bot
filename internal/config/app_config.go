package config

import (
	"strconv"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type AppConfig struct {
	DBConfig              DBConfig
	TelegramToken         string `required:"false" envconfig:"TELEGRAM_TOKEN"`
	TelegramBotDebug      bool   `envconfig:"TELEGRAM_BOT_DEBUG" default:"true"`
	TelegramUpdateTimeout int    `envconfig:"TELEGRAM_UPDATE_TIMEOUT" default:"60"`
}

type DBConfig struct {
	DataSource string `required:"true" envconfig:"DATA_SOURCE" default:"postgres"`
	Port       string `required:"true" envconfig:"DATA_SOURCE_PORT" default:"5432"`
	Host       string `required:"true" envconfig:"DATA_SOURCE_HOST" default:"127.0.0.1"`
	Name       string `required:"true" envconfig:"DATA_SOURCE_NAME" default:"postgres"`
	Username   string `required:"true" envconfig:"DATA_SOURCE_USERNAME" default:"postgres"`
	Password   string `required:"true" envconfig:"DATA_SOURCE_PASSWORD"`
}

// NewConfig Loads environmental variables from .env file and populates AppConfig struct with found values
func NewConfig(logger *zap.Logger) (*AppConfig, error) {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file", zap.Error(err))
		return nil, err
	}

	cfg := AppConfig{}
	// Process method tries to search environmental variable first by prefix+envconfig
	// if value isn't found it uses specified envconfig
	if err := envconfig.Process("", &cfg); err != nil {
		logger.Error("Error occurred while processing configs", zap.Error(err))
		return nil, err
	}
	logger.Info("Application configured", zap.String("Project data source", cfg.DBConfig.DataSource),
		zap.String("Database port", cfg.DBConfig.Port), zap.String("Database host", cfg.DBConfig.Host),
		zap.String("Database name", cfg.DBConfig.Name), zap.String("Database username", cfg.DBConfig.Username),
		zap.String("Database password", cfg.DBConfig.Password), zap.String("Telegram token", cfg.TelegramToken),
		zap.String("Telegram bot debug", strconv.FormatBool(cfg.TelegramBotDebug)),
		zap.String("Telegram update timeout", strconv.Itoa(cfg.TelegramUpdateTimeout)))
	return &cfg, nil
}
