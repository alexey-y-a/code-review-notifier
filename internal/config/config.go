package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Server struct {
		Port string
	}
	Telegram struct {
		Token string
	}
	GitHub struct {
		Secret string
	}
	Database struct {
		Host     string
		Port     int
		Username string
		Password string
		Name     string
	}
	LogLevel string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		zap.L().Error("Failed to read config file", zap.Error(err))
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		zap.L().Error("Failed to unmarshal config", zap.Error(err))
		return nil, err
	}

	return &cfg, nil
}
