package telegram

import (
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/telegram-bot-api.v4"

	"github.com/horechek/poster/app/di"
)

type Bot struct {
	client *tgbotapi.BotAPI
	app    *di.Services
}

func NewBot(app *di.Services, proxy *http.Client, token string, debug bool) (*Bot, error) {
	bot := &Bot{
		app: app,
	}
	var err error
	bot.client, err = tgbotapi.NewBotAPIWithClient(token, proxy)
	if err != nil {
		return nil, err
	}

	bot.app.Logger.Infow("telegram authorized", "account", bot.client.Self.UserName)
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

		b.app.Logger.Infow("msg from telegram", "msg", update.Message)

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