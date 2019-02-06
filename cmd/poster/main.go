package main

import (
	"flag"

	"github.com/jonboulle/clockwork"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/horechek/poster/app/database"
	"github.com/horechek/poster/app/di"
	"github.com/horechek/poster/app/server"
	"github.com/horechek/poster/app/telegram"
)

var env = flag.String("env", "dev", "application environment")

func main() {
	flag.Parse()

	// init logger
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer log.Sync()
	shugar := log.Sugar()

	// init db
	dbconf := database.Config{
		Debug:   true,
		Driver:  "sqlite3",
		Connect: "poster.db",
	}
	db := database.NewDatabase(dbconf, shugar)

	// migrations
	db.AutoMigrate(database.User{})
	db.AutoMigrate(database.Post{})

	// init di
	services := &di.Services{
		Database: db,
		Logger:   shugar,
		Clock:    clockwork.NewRealClock(),
	}

	// init telegram messenger
	proxy, err := telegram.NewProxy("188.166.21.43:1111", "artem", "589311")
	if err != nil {
		shugar.Fatalw("error on start proxy", zap.Error(err))
	}

	tg, err := telegram.NewTelegram(services, proxy, "551357910:AAHvqlvWmYZvqhLl_M42qjaG0n3O0jRDqG8", "@adtechbeer", true)
	if err != nil {
		shugar.Fatalw("error on start telegram messenger", zap.Error(err))
	}

	bot, err := telegram.NewBot(services, proxy, "551357910:AAHvqlvWmYZvqhLl_M42qjaG0n3O0jRDqG8", true)
	if err != nil {
		shugar.Fatalw("error on start telegram bot", zap.Error(err))
	}
	go bot.Run()

	server := server.NewServer(services, viper.GetString("server.port"))
	server.Run()
}
