package config

import (
	"strconv"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type AppConfig struct {
	DataSource			  string `required:"true" envconfig:"DATA_SOURCE" default:"postgres"`
	DbPort				  string `required:"true" envconfig:"DATA_SOURCE_PORT" default:"5432"`
	DbHost				  string `required:"true" envconfig:"DATA_SOURCE_HOST" default:"127.0.0.1"`
	DbName				  string `required:"true" envconfig:"DATA_SOURCE_NAME" default:"postgres"`
	DbUsername			  string `required:"true" envconfig:"DATA_SOURCE_USERNAME" default:"postgres"`
	DbPassword			  string `required:"true" envconfig:"DATA_SOURCE_PASSWORD" default:"postgres"`
	TelegramToken		  string `required:"true" envconfig:"TELEGRAM_TOKEN"`
	TelegramBotDebug	  bool	 `envconfig:"TELEGRAM_BOT_DEBUG" default:"false"`
	TelegramUpdateTimeout int	 `envconfig:"TELEGRAM_UPDATE_TIMEOUT" default:"60"`
}

// NewConfig Loads environmental variables from .env file and populates AppConfig struct with found values
func NewConfig(logger *zap.Logger) (*AppConfig, error) {
	cfg := AppConfig{}
	// Process method tries to search environmental variable first by prefix+envconfig
	// if value isn't found it uses specified envconfig
	if err := envconfig.Process("", &cfg); err != nil {
		logger.Error("Error occurred while processing configs", zap.Error(err))
		return nil, err
	}
	logger.Info("Application configured", zap.String("Project data source", cfg.DataSource),
		zap.String("Database port", cfg.DbPort), zap.String("Database host", cfg.DbHost),
		zap.String("Database name", cfg.DbName), zap.String("Database username", cfg.DbUsername),
		zap.String("Database password", cfg.DbPassword), zap.String("Telegram token", cfg.TelegramToken),
		zap.String("Telegram bot debug", strconv.FormatBool(cfg.TelegramBotDebug)),
		zap.String("Telegram update timeout", strconv.Itoa(cfg.TelegramUpdateTimeout)))
	return &cfg, nil
}

