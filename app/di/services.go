package di

import (
	"github.com/horechek/poster/app/telegram"
	"github.com/jinzhu/gorm"
	"github.com/jonboulle/clockwork"
	"go.uber.org/zap"
)

type Services struct {
	Database *gorm.DB
	Logger   *zap.SugaredLogger
	Clock    clockwork.Clock
	Telegram *telegram.Telegram
}
