package postgres

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type DataSourceConfig struct {
	DataSource string `required:"true" envconfig:"DATA_SOURCE" default:"postgres"`
	DbPort     string `required:"true" envconfig:"DATA_SOURCE_PORT" default:"5432"`
	DbHost     string `required:"true" envconfig:"DATA_SOURCE_HOST" default:"127.0.0.1"`
	DbName     string `required:"true" envconfig:"DATA_SOURCE_NAME" default:"postgres"`
	DbUsername string `required:"true" envconfig:"DATA_SOURCE_USERNAME" default:"postgres"`
	DbPassword string `required:"true" envconfig:"DATA_SOURCE_PASSWORD" default:"postgres"`
}

// NewDataSourceConfig Loads environmental variables from .env file and populates DataSourceConfig struct with found values
func NewDataSourceConfig(logger *zap.Logger) (*DataSourceConfig, error) {
	envPath := "internal/config/.env"
	if err := godotenv.Load(envPath); err != nil {
		logger.Error("No .env file found for path", zap.String("path", envPath), zap.Error(err))
	}
	cfg := DataSourceConfig{}
	// Process method tries to search environmental variable first by prefix+envconfig
	// if value isn't found it uses specified envconfig
	if err := envconfig.Process("", &cfg); err != nil {
		logger.Error("Error occurred while processing configs", zap.Error(err))
		return nil, err
	}
	logger.Info("Datasource connection configuration.", zap.String("Project data source", cfg.DataSource),
		zap.String("Database port", cfg.DbPort), zap.String("Database host", cfg.DbHost),
		zap.String("Database name", cfg.DbName), zap.String("Database username", cfg.DbUsername),
		zap.String("Database password", cfg.DbPassword))
	return &cfg, nil
}
