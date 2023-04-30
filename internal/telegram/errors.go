package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	errorInvalidUrl   = errors.New("this is not valid link")
	errorUnauthorized = errors.New("user unauthorized")
	errorUnableToSave = errors.New("unable to save this link")
)

func (b *Bot) HandleErrors(chatId int64, err error) {
	msg := tgbotapi.NewMessage(chatId, b.messages.Default)
	switch err {
	case errorInvalidUrl:
		msg.Text = b.messages.InvalidLink
		b.bot.Send(msg)

	case errorUnauthorized:
		msg.Text = b.messages.Unauthorized
		b.bot.Send(msg)

	case errorUnableToSave:
		msg.Text = b.messages.UnableToSave
		b.bot.Send(msg)

	default:
		b.bot.Send(msg)
	}

}
