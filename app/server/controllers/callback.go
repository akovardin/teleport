package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"go.uber.org/zap"

	db "github.com/horechek/teleport/app/database"
	"github.com/horechek/teleport/app/di"
	"github.com/horechek/teleport/app/telegram"
)

type CallbackController struct {
	services *di.Services
}

func NewCallbackController(services *di.Services) *CallbackController {
	return &CallbackController{
		services: services,
	}
}

func (c *CallbackController) Callback(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return err
	}

	model := &db.Integration{}
	model, err = model.FindOne(c.services.Database, db.Condition{
		Params: []db.Param{
			{
				Field: "id",
				Value: id,
			},
		},
	})

	if err != nil {
		return err
	}

	fmt.Println("address", model.ProxyAddress)

	proxy, err := telegram.NewProxy(model.ProxyAddress, model.ProxyUser, model.ProxyPass)
	if err != nil {
		c.services.Logger.Error("error on create proxy", zap.Error(err))
		return err
	}

	tg, err := telegram.NewTelegram(c.services, proxy, model.Token, true)
	if err != nil {
		c.services.Logger.Error("error on create telegram", zap.Error(err))
		return err
	}

	var resp struct {
		Type   string `json:"type"`
		Object struct {
			Id   int    `json:"id"`
			Text string `json:"text"`
		} `json:"object"`
	}

	if err := ctx.Bind(&resp); err != nil {
		return err
	}

	if resp.Type == "confirmation" {
		return ctx.String(http.StatusOK, model.Secret)
	}

	if resp.Object.Text == "" {
		return ctx.String(http.StatusOK, "ok")
	}

	post := &db.Post{
		Body:   resp.Object.Text,
		Status: true,
	}

	if err := c.services.Database.Model(post).Save(post).Error; err != nil {
		return err
	}

	if err := tg.Send(model.Channel, post); err != nil {
		c.services.Logger.Error("error on send message to telegram", zap.Error(err))
		return err
	}

	return ctx.String(http.StatusOK, "ok")
}

/*
{
   "type":"wall_post_new",
   "object":{
      "id":125,
      "from_id":-169563832,
      "owner_id":-169563832,
      "date":1549378680,
      "marked_as_ads":0,
      "post_type":"post",
      "text":"Реклама убивает пуш-уведомления\n\nВеб-уведомления поддерживаются во всех современных браузерах: Хроме, Огнелисе и т.д. Это часть Progressive Web Apps - прогрессивных веб-приложений или PWAs. \nК сожалению, время их использования подходит к концу. Пользователи их блокируют или совсем отключают.\n\nПуши стали такой раздражающей штукой, что разработчики браузеров начали задумываться. Пока пуши еще работают, но уже сейчас заблокировать или отписаться от \nуведомлений становится проще с каждым релизом браузера. \n\nhttp:\/\/vancelucas.com\/blog\/the-adtech-industry-is-ruining-web-push-notifications-for-pwas\/",
      "can_edit":1,
      "created_by":188859983,
      "can_delete":1,
      "attachments":[
         {
            "type":"link",
            "link":{
               "url":"http:\/\/vancelucas.com\/blog\/the-adtech-industry-is-ruining-web-push-notifications-for-pwas\/",
               "title":"The AdTech Industry is Ruining Web Push Notifications for PWAs | Vance Lucas",
               "caption":"vancelucas.com",
               "description":"",
               "photo":{
                  "id":456260428,
                  "album_id":-27,
                  "owner_id":2000060696,
                  "sizes":[
                     {
                        "type":"l",
                        "url":"https:\/\/pp.userapi.com\/c850224\/v850224317\/d3a1c\/l5PaD4u6KSA.jpg",
                        "width":537,
                        "height":240
                     },
                     {
                        "type":"m",
                        "url":"https:\/\/pp.userapi.com\/c850224\/v850224317\/d3a19\/rAyUkFrBqjc.jpg",
                        "width":130,
                        "height":130
                     },
                     {
                        "type":"p",
                        "url":"https:\/\/pp.userapi.com\/c850224\/v850224317\/d3a1b\/mD9BPXI1Dxc.jpg",
                        "width":260,
                        "height":260
                     },
                     {
                        "type":"s",
                        "url":"https:\/\/pp.userapi.com\/c850224\/v850224317\/d3a18\/BiUMVw8gB4A.jpg",
                        "width":75,
                        "height":75
                     },
                     {
                        "type":"x",
                        "url":"https:\/\/pp.userapi.com\/c850224\/v850224317\/d3a1a\/uCYs-N1yt8E.jpg",
                        "width":150,
                        "height":150
                     }
                  ],
                  "text":"",
                  "date":1549361236
               }
            }
         }
      ],
      "comments":{
         "count":0
      }
   },
   "group_id":169563832,
   "secret":"w45gt5"
}
*/
