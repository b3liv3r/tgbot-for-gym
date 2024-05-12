package handlers

import (
	"github.com/b3liv3r/tgbot-for-gym/modules/users/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramUserHandler struct {
	userService service.Userer
	bot         *tgbotapi.BotAPI
}
