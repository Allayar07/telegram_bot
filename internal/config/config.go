package config

import (
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	TelegramToken string
	ConsumerKey   string
	AuthServerUrl string

	TelegramBotURL string `mapstructure:"bot_url"`
	DbPath         string `mapstructure:"bolt_db"`

	Messages Messages
}

type Messages struct {
	Responses
	Errors
}

type Responses struct {
	Start          string `mapstructure:"start"`
	AlreadyAuth    string `mapstructure:"already_authorized"`
	UnknownCommand string `mapstructure:"unknown_command"`
	SuccessLink    string `mapstructure:"success_save"`
}

type Errors struct {
	Default      string `mapstructure:"default"`
	InvalidLink  string `mapstructure:"invalid_url"`
	Unauthorized string `mapstructure:"unauthorized"`
	UnableToSave string `mapstructure:"unable_to_save"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("message.responses", &cfg.Messages.Responses); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("message.error", &cfg.Messages.Errors); err != nil {
		return nil, err
	}

	if err := ParseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ParseEnv(cfg *Config) error {
	os.Setenv("TOKEN", "5735493357:AAEy0iE3QYewXa57YS87ad9pJ3_ZptdJm4M")
	os.Setenv("CONSUMER_KEY", "104781-e549e898c3bde1b7e0bb77b")
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
	cfg.ConsumerKey = viper.GetString("consumer_key")
	cfg.AuthServerUrl = viper.GetString("auth_server_url")

	return nil
}
