package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type TeleBot struct {
	Bot *tgbotapi.BotAPI
}

func NewTeleBot(token string) *TeleBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	return &TeleBot{
		Bot: bot,
	}
}

func (teleBot *TeleBot) SendMessage(text string, groupId int64) {
	msg := tgbotapi.NewMessage(groupId, text)
	teleBot.Bot.Send(msg)
}
