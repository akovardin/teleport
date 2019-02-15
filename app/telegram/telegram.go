package telegram

import (
	"net/http"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/horechek/teleport/app/database"
	"github.com/horechek/teleport/app/di"
)

type Telegram struct {
	api      *tgbotapi.BotAPI
	services *di.Services
}

func NewTelegram(services *di.Services, proxy *http.Client, token string, debug bool) (*Telegram, error) {
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

	client.Debug = debug

	return &Telegram{
		api:      client,
		services: services,
	}, nil
}

func (t *Telegram) Send(chat string, post *database.Post) error {
	text := post.Body
	if post.Title != "" {
		text = post.Title + " / / " + text
	}

	msg := tgbotapi.NewMessageToChannel(chat, text)
	if _, err := t.api.Send(msg); err != nil {
		return err
	}

	post.Send = true
	post.UpdatedAt = time.Now().Unix()

	return t.services.Database.Save(post).Error
}
