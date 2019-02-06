package main

import (
	"github.com/jonboulle/clockwork"
	"go.uber.org/zap"

	"github.com/horechek/teleport/app/database"
	"github.com/horechek/teleport/app/di"
	"github.com/horechek/teleport/app/server"
	"github.com/horechek/teleport/app/telegram"
)

var (
	token      = "551357910:AAHvqlvWmYZvqhLl_M42qjaG0n3O0jRDqG8"
	chanelName = "@adtechbeer"

	proxyAddress = "188.166.21.43:1111"
	proxyUser    = "artem"
	proxyPass    = "589311"

	port = "8080"
)

func main() {
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
		Connect: "teleport.db",
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
	proxy, err := telegram.NewProxy(proxyAddress, proxyUser, proxyPass)
	if err != nil {
		shugar.Fatalw("error on start proxy", zap.Error(err))
	}

	tg, err := telegram.NewTelegram(services, proxy, token, chanelName, true)
	if err != nil {
		shugar.Fatalw("error on start telegram messenger", zap.Error(err))
	}

	bot, err := telegram.NewBot(services, proxy, token, true)
	if err != nil {
		shugar.Fatalw("error on start telegram bot", zap.Error(err))
	}
	go bot.Run()

	server := server.NewServer(services, tg, port)
	server.Run()
}
