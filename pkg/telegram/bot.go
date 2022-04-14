package telegram

import (
	"log"

	"github.com/Sacrament0/telegram-bot/pkg/config"
	"github.com/Sacrament0/telegram-bot/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
)

// Структура описывающая бота
type Bot struct {
	bot *tgbotapi.BotAPI
	pocketClient *pocket.Client
	TokenRepository repository.TokenRepository
	redirectURL string
	messages config.Messages
}

// Конструктор для бота. Создает переменную структуру типа Bot и помещает туда созданного бота типа *tgbotapi.BotAPI
// Почему принимает ссылку - Чтобы не копировать всего бота в метод и не выделять память
// а также чтобы была возможность изменять данные
func NewBot(bot *tgbotapi.BotAPI, pocketClient *pocket.Client, tr repository.TokenRepository, redirectURL string, messages config.Messages) *Bot {
	return &Bot{bot: bot, pocketClient: pocketClient, redirectURL: redirectURL, TokenRepository: tr, messages: messages}
}

// Start создает канал для отправки сообщений и обрабатывает входящие сообщения
func (b *Bot) Start() error {

	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	// Cоздание канала для получения сообщений
	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	// Обработка сообщений
	b.handleUpdates(updates)

	return nil

}

// handleUpdates получает и отправляет сообщения пользователю
func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	// Итерация по каналу
	for update := range updates {
		// если сообщение пустое, скипаем цикл
		if update.Message == nil {
			continue
		}
		// проверка является ли сообщение командой
		if update.Message.IsCommand() {

			//обработка команды
			if err := b.handleCommand(update.Message); err != nil {
				b.HandleError(update.Message.Chat.ID, err)
			}
			continue

		}

		// обработка сообщения
		if err := b.handleMessage(update.Message); err != nil {
			b.HandleError(update.Message.Chat.ID, err)
		}


	}

}

// initUpdatesChannel инициализирует канал для получения-передачи сообщений
func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	// Создание конфигурации для получения обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Создание канала для получения значеий от API
	return b.bot.GetUpdatesChan(u)
}
