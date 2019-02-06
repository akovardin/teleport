package controllers

import (
	"github.com/horechek/poster/app/telegram"
	"go.uber.org/zap"
	"net/http"

	"github.com/labstack/echo"

	"github.com/horechek/poster/app/database"
	"github.com/horechek/poster/app/di"
)

type VKController struct {
	services *di.Services
	tg       *telegram.Telegram
}

func NewVKСontroller(services *di.Services, tg *telegram.Telegram) *VKController {
	return &VKController{
		services: services,
		tg:       tg,
	}
}

func (c *VKController) Callback(ctx echo.Context) error {
	var resp struct {
		Object struct {
			Id   int    `json:"id"`
			Text string `json:"text"`
		} `json:"object"`
	}

	if err := ctx.Bind(&resp); err != nil {
		return err
	}

	post := &database.Post{
		Body:   resp.Object.Text,
		Status: true,
	}

	if err := c.services.Database.Model(post).Save(post).Error; err != nil {
		return err
	}

	if err := c.tg.Send(&database.Post{Title: "test", Body: "test"}); err != nil {
		c.services.Logger.Warnw("error on send message to telegram", zap.Error(err))
	}

	return ctx.String(http.StatusOK, "5d2c1139")
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
