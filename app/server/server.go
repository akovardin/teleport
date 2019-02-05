package server

import (
	"io/ioutil"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/horechek/poster/app/database"
	"github.com/horechek/poster/app/di"
	"github.com/horechek/poster/app/server/controllers"
)

type Server struct {
	services *di.Services
	port     string
}

func NewServer(services *di.Services, port string) *Server {
	return &Server{
		services: services,
		port:     port,
	}
}

func (s *Server) Run() {
	e := echo.New()
	e.Logger.SetOutput(ioutil.Discard)
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,
			echo.HeaderAuthorization},
		ExposeHeaders:    []string{"X-Total-Count"},
		AllowCredentials: true,
	}))

	skiped := map[string]struct{}{
		"/api/users/login":    {},
		"/api/users/register": {},
		"/api/users/restore":  {},
		"/metrics":            {},
	}

	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper: func(c echo.Context) bool {
			_, ok := skiped[c.Path()]
			return ok
		},
		Validator: func(token string, context echo.Context) (bool, error) {
			u := database.User{}
			err := s.services.Database.Model(u).
				Where("token = ?", token).
				First(&u).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return false, err
			}

			if err == gorm.ErrRecordNotFound {
				return false, nil
			}

			context.Set("user", u.ID)

			return true, nil
		},
	}))

	users := controllers.NewUsersController(s.services)
	posts := controllers.NewPostsController(s.services)
	vk := controllers.NewVKVontroller(s.services)

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	api := e.Group("api")
	// api
	api.POST("/users/login", users.Login)
	api.POST("/users/register", users.Register)
	api.POST("/users/restore", users.Restore)
	api.POST("/users/update", users.Update)

	api.GET("/posts", posts.List)
	api.POST("/posts/:id", posts.Update)
	api.POST("/posts", posts.Create)
	api.DELETE("/posts/:id", posts.Remove)

	api.POST("/callback", vk.Callback)

	// Start server
	s.services.Logger.Infow("start api server", "port", s.port)
	s.services.Logger.Fatalw("error on starting server", "err", e.Start(":"+s.port))
}
