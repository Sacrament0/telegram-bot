package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken     string
	PocketConsumerKey string
	AuthServerURL     string
	TelegramBotURL    string `mapstructure:"bot_url"`
	DBParth           string `mapstructure:"db_file"`

	Messages Messages
}

type Messages struct {
	Errors
	Responses
}

type Errors struct {
	Default      string `mapstructure:"default"`
	InvalidURL   string `mapstructure:"invalid_url"`
	Unauthorized string `mapstructure:"unauthorized"`
	UnableToSave string `mapstructure:"unable_to_save"`
}

type Responses struct {
	Start               string `mapstructure:"start"`
	AlreadyUnauthorized string `mapstructure:"already_authorized"`
	SavedSuccessfully   string `mapstructure:"saved_saccessfully"`
	UnknownCommand      string `mapstructure:"unknown_command"`
}

// Инициализация конфига
func Init() (*Config, error) {
	// Путь к директории конфига (имя папки)
	viper.AddConfigPath("configs")
	// Имя конфига в папке
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.responses", &cfg.Messages.Responses); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.errors", &cfg.Messages.Errors); err != nil {
		return nil, err
	}

	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Парсит переменные окружения
func parseEnv(cfg *Config) error {
	os.Setenv("TOKEN", "5141221734:AAHSxTYUHYNfbXascUBKn9pTi4v0WvToSxU")
	os.Setenv("CONSUMER_KEY", "101425-5b0a22249f035a3b41b52d0")
	os.Setenv("AUTH_SERVER_URL", "http://localhost/")

	if err := viper.BindEnv("token"); err != nil {
		return err
	}

	if err := viper.BindEnv("consumer_key"); err != nil {
		return err
	}

	if err := viper.BindEnv("auth_server_url"); err != nil {
		return err
	}

	cfg.TelegramToken = viper.GetString("token")
	cfg.PocketConsumerKey = viper.GetString("consumer_key")
	cfg.TelegramBotURL = viper.GetString("auth_server_url")

	return nil
}

