package telegram

import (
	"net/http"
	"strconv"
	"time"

	"gopkg.in/telegram-bot-api.v4"

	"github.com/horechek/poster/app/database"
	"github.com/horechek/poster/app/di"
)

type Telegram struct {
	chat     string
	api      *tgbotapi.BotAPI
	services *di.Services
}

func NewTelegram(services *di.Services, proxy *http.Client, token, chat string) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPIWithClient(token, proxy)
	if err != nil {
		return nil, err
	}

	return &Telegram{
		chat:     chat,
		api:      bot,
		services: services,
	}, nil
}

func (t *Telegram) Send(post *database.Post) error {
	chat, _ := strconv.Atoi(t.chat)
	msg := tgbotapi.NewMessage(int64(chat), post.Title+" / / "+post.Body)
	if _, err := t.api.Send(msg); err != nil {
		return err
	}

	post.Send = true
	post.UpdatedAt = time.Now().Unix()

	return t.services.Database.Save(post).Error
}
