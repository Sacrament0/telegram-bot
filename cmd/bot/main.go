package main

import (
	"log"

	"github.com/Sacrament0/telegram-bot/pkg/config"
	"github.com/Sacrament0/telegram-bot/pkg/repository"
	"github.com/Sacrament0/telegram-bot/pkg/repository/boltdb"
	"github.com/Sacrament0/telegram-bot/pkg/server"
	"github.com/Sacrament0/telegram-bot/pkg/telegram"
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
)

func main() {
	// Создание конфига
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	// Создание нового бота
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}
	// Флаг для вывода логов в консоли
	bot.Debug = true

	// Создание pocket клинта
	pocketClient, err := pocket.NewClient(cfg.PocketConsumerKey)
	if err != nil {
		log.Fatal(err)
	}

	// Создание БД
	db, err := initDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Оборачивание БД в структуру TokenRepository
	tokenRepository := boltdb.NewTokenRepository(db)

	// Оборачивание бота в структуру Bot
	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, cfg.AuthServerURL, cfg.Messages)

	// Запуск сервера для обработки редиректов от pocket
	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, cfg.TelegramBotURL)

	// Запускаем бота
	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	// Запускаем сервер
	if err := authorizationServer.Start(); err != nil {
		log.Fatal(err)
	}

}

// initDB создаёт базу данных и бакеты для хранения токенов
func initDB(cfg *config.Config) (*bolt.DB, error) {

	// Создание базы данных
	db, err := bolt.Open(cfg.DBParth, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}

// Create check
