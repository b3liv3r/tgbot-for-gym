package main

import (
	loggerx "github.com/b3liv3r/logger"
	"github.com/b3liv3r/tgbot-for-gym/config"
	"github.com/b3liv3r/tgbot-for-gym/modules"
	"github.com/b3liv3r/tgbot-for-gym/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"

	"os"
)

func main() {
	AppConfig := config.MustLoadConfig(".env")

	logger := loggerx.InitLogger(AppConfig.Name, AppConfig.Production)

	clients := modules.NewClients()
	services := modules.NewServices(logger, clients)

	botToken := os.Getenv("TOKEN")
	if botToken == "" {
		log.Fatal("Не задан токен бота TOKEN")
	}
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal("Ошибка создания бота:", err)
	}

	// Создание роутера для обработки обновлений от бота
	router := telegram.NewTelegramRouter(services, bot, logger)

	// Запуск роутера
	router.Start()
}
