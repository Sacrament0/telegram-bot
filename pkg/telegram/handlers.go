package telegram

import (
	"context"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
)

const (
	commandStart = "start"
)

// Метод для обработки команд
func (b *Bot) handleCommand(message *tgbotapi.Message) error {

	// получения значения команды
	switch message.Command() {

	case commandStart:

		return b.handleStartCommand(message)

	default:

		return b.handleUnknownCommand(message)

	}

}

// Метод для обработки сообщений
func (b *Bot) handleMessage(message *tgbotapi.Message) error {

	// Проверка валидности формата url
	_, err := url.ParseRequestURI(message.Text)
	if err != nil {
		return errInvalidURL
	}

	// Получение access токена и проверка авторизации
	accessToken, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return errUnauthorized
	}

	// Сохранение ссылки в покете
	if err := b.pocketClient.Add(context.Background(), pocket.AddInput{
		AccessToken: accessToken,
		URL:         message.Text,
	}); err != nil {
		return errUnableToSave
	}

	// Создание структуры для ответоного сообщения с указанием:
	// Chat ID куда отправляется сообщение и текст сообщения
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.SavedSuccessfully)

	// Указание параметра, что отправляемое сообщение является ответом на полученное сообщение
	msg.ReplyToMessageID = message.MessageID

	// Отправка сообщения
	_, err = b.bot.Send(msg)

	return err
}

// Метод для обработки стартовой команды
func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {

	_, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return b.initAuthorizationProcess(message)
	}

	//если пользователь авторизирован, сообщаем ему об этом
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.AlreadyUnauthorized)

	// Отправка сообщения
	b.bot.Send(msg)

	return err
}

// Метод для обработки неизвестных команд
func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {

	// Создание структуры для ответоного сообщения с указанием:
	// Chat ID куда отправляется сообщение и текст сообщения (message.Text)
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.UnknownCommand)

	// Указание параметра, что отправляемое сообщение является ответом на полученное сообщение
	msg.ReplyToMessageID = message.MessageID

	// Отправка сообщения
	_, err := b.bot.Send(msg)

	return err
}
