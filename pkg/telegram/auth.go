package telegram

import (
	"context"
	"fmt"

	"github.com/Sacrament0/telegram-bot/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)


// Начинает авторизацию пользователя (генерит ссылку на доступ к аккаунту)
func (b *Bot) initAuthorizationProcess(message *tgbotapi.Message) error {
	// Генерация ссылки на аутентификацию
	authLink, err := b.generateAuthorizationLink(message.Chat.ID)
	if err != nil {
		return err
	}
	// Создание структуры для ответного сообщения с указанием:
	// Chat ID куда отправляется сообщение и текст сообщения (message.Text)
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(b.messages.Start, authLink))

	// Указание параметра, что отправляемое сообщение является ответом на полученное сообщение
	msg.ReplyToMessageID = message.MessageID

	// Отправка сообщения
	_, err = b.bot.Send(msg)

	return err
}

// Получает access token из бакета
func (b *Bot) getAccessToken(chatID int64) (string, error) {
	return b.TokenRepository.Get(chatID, repository.AccessTokens)
}

// Создаёт ссылку для авторизации пользоватея
func (b *Bot) generateAuthorizationLink(chatID int64) (string, error) {

	redirectURL := b.generateRedirectURL(chatID)

	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), redirectURL)
	if err != nil {
		return "", err
	}

	// сохранение реквест-токена в бакет
	if err = b.TokenRepository.Save(chatID, requestToken, repository.RequestTokens); err != nil {
		return "", err
	}

	return b.pocketClient.GetAuthorizationURL(requestToken, redirectURL)
}

// Генерирует ссылку на редирект
func (b *Bot) generateRedirectURL(chatID int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.redirectURL, chatID)
}
