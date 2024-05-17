package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type AppConfig struct {
	DbPort        string `required:"true" envconfig:"DATA_SOURCE_PORT" default:"5432"`
	DbHost        string `required:"true" envconfig:"DATA_SOURCE_HOST" default:"127.0.0.1"`
	DbName        string `required:"true" envconfig:"DATA_SOURCE_NAME" default:"postgres"`
	DbUsername    string `required:"true" envconfig:"DATA_SOURCE_USERNAME" default:"postgres"`
	DbPassword    string `required:"true" envconfig:"DATA_SOURCE_PASSWORD" default:"postgres"`
	TelegramToken string `required:"true" envconfig:"TELEGRAM_TOKEN"`
}

// NewConfig Loads environmental variables from .env file and populates AppConfig struct with found values
func NewConfig() (*AppConfig, error) {
	if err := godotenv.Load("config/.env"); err != nil {
		log.Print("No .env file found")
	}
	cfg := AppConfig{}
	// Process method tries to search environmental variable first by prefix+envconfig
	// if value isn't found it uses specified envconfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	log.Printf("%+v\n", cfg)
	return &cfg, nil
}
