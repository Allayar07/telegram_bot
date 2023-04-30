package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
	"telegramBot/internal/config"
	"telegramBot/internal/repository"
	"telegramBot/internal/repository/botdb"
	"telegramBot/internal/server"
	"telegramBot/internal/telegram"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient(cfg.ConsumerKey)
	if err != nil {
		log.Fatal(err)
	}

	db, err := repository.InitBoltDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	tokenRepository := botdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, cfg.AuthServerUrl, cfg.Messages)

	authServer := server.NewAuthServer(pocketClient, tokenRepository, cfg.TelegramBotURL)

	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)

		}
	}()

	if err := authServer.Run(); err != nil {
		log.Fatal(err)
	}
}
