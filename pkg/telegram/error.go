package telegram

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	errInvalidURL   = errors.New("url is invalid")
	errUnauthorized = errors.New("url is not authorized")
	errUnableToSave = errors.New("unable to save")
)

// Обработка ошибок
func (b *Bot) HandleError(chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, b.messages.Default)
	switch err {
	case errInvalidURL:
		msg.Text = b.messages.InvalidURL
		b.bot.Send(msg)
	case errUnauthorized:
		msg.Text = b.messages.Unauthorized
		b.bot.Send(msg)
	case errUnableToSave:
		msg.Text = b.messages.UnableToSave
		b.bot.Send(msg)
	default:
		b.bot.Send(msg)

	}
}
