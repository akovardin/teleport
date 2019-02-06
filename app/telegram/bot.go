package telegram

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/horechek/poster/app/di"
)

type Bot struct {
	client   *tgbotapi.BotAPI
	services *di.Services
}

func NewBot(services *di.Services, proxy *http.Client, token string, debug bool) (*Bot, error) {
	var (
		client *tgbotapi.BotAPI
		err    error
	)

	if proxy != nil {
		client, err = tgbotapi.NewBotAPIWithClient(token, proxy)
		if err != nil {
			return nil, err
		}
	} else {
		client, err = tgbotapi.NewBotAPI(token)
		if err != nil {
			return nil, err
		}
	}

	bot := &Bot{
		services: services,
		client:   client,
	}

	bot.services.Logger.Infow("telegram authorized", "account", bot.client.Self.UserName)
	bot.client.Debug = debug

	return bot, nil
}

func (b *Bot) Run() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.client.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		b.services.Logger.Infow("msg from telegram", "msg", update.Message)

		if !strings.HasPrefix(update.Message.Text, "/start") {
			continue
		}

		// welcome message
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Приветствую! Ваше имя пользователе в телеграме: "+
			update.Message.From.UserName+" и чат: "+
			strconv.Itoa(int(update.Message.Chat.ID))+".")

		b.client.Send(msg)
	}

	return nil
}
