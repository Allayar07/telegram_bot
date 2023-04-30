package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"telegramBot/internal/repository"
)

func (b *Bot) InitAuthClient(message *tgbotapi.Message) error {
	authLink, err := b.generateAuthURL(message.Chat.ID)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(message.Chat.ID,
		fmt.Sprintf(startReply, authLink))
	_, err = b.bot.Send(msg)

	return err
}

func (b *Bot) getAccessToken(chatId int64) (string, error) {
	return b.tokenRepository.Get(chatId, repository.AccessToken)
}

func (b *Bot) generateAuthURL(chatID int64) (string, error) {
	RedirectURL := b.generateRedirectURL(chatID)

	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), RedirectURL)
	if err != nil {
		return "", err
	}

	if err := b.tokenRepository.Save(chatID, requestToken, repository.RequestToken); err != nil {
		return "", err
	}

	return b.pocketClient.GetAuthorizationURL(requestToken, RedirectURL)
}

func (b *Bot) generateRedirectURL(chatId int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.redirectUrl, chatId)
}
