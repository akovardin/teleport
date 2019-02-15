package main

import (
	"fmt"
	"os"

	"github.com/jonboulle/clockwork"
	"go.uber.org/zap"
	"github.com/urfave/cli"


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

func main()  {
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
	db.AutoMigrate(database.Integration{})
	db.AutoMigrate(database.Post{})

	// init di
	services := &di.Services{
		Database: db,
		Logger:   shugar,
		Clock:    clockwork.NewRealClock(),
	}

	app := cli.NewApp()
	app.Name = "cli"
	app.Usage = "CPA network services"
	app.Commands = []cli.Command{
		{
			Name: "users",
			Subcommands: []cli.Command{
				{
					Name: "list",
					Action: func(c *cli.Context) error {
						users := []database.User{}
						if err := db.Model(database.User{}).Find(&users).Error; err != nil {
							return err
						}

						for _, u := range users {
							fmt.Println(u)
						}

						return nil
					},
				},
				{
					Name: "add",
					Flags: []cli.Flag {
						cli.StringFlag{
							Name: "email",
							Value: "admit@adtech.beer",
							Usage: "user email for login",
						},
						cli.StringFlag{
							Name: "pass",
							Value: "admin",
							Usage: "user password for login",
						},
						cli.StringFlag{
							Name: "token",
							Value: "123",
							Usage: "user token for login",
						},
					},
					Action: func(c *cli.Context) error {
						u := database.User{
							Email: c.String("email"),
							Password: c.String("pass"),
							Token: c.String("token"),
						}
						return db.Save(&u).Error
					},
				},
			},
		},
		{
			Name: "server",
			Action: func(c *cli.Context) error {
				// init telegram messenger
				proxy, err := telegram.NewProxy(proxyAddress, proxyUser, proxyPass)
				if err != nil {
					shugar.Fatalw("error on start proxy", zap.Error(err))
					return err
				}

				tg, err := telegram.NewTelegram(services, proxy, token, chanelName, true)
				if err != nil {
					shugar.Fatalw("error on start telegram messenger", zap.Error(err))
					return err
				}

				bot, err := telegram.NewBot(services, proxy, token, true)
				if err != nil {
					shugar.Fatalw("error on start telegram bot", zap.Error(err))
					return err
				}
				go bot.Run()

				server := server.NewServer(services, tg, port)
				server.Run()

				return nil
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

