package server

import (
	"github.com/horechek/teleport/pkg/middleware/static"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"

	"github.com/horechek/teleport/app/database"
	"github.com/horechek/teleport/app/di"
	"github.com/horechek/teleport/app/server/controllers"
	_ "github.com/horechek/teleport/statik"
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
	server := echo.New()
	server.Logger.SetOutput(ioutil.Discard)
	server.Use(middleware.Recover())
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,
			echo.HeaderAuthorization},
		ExposeHeaders:    []string{"X-Total-Count"},
		AllowCredentials: true,
	}))

	skiped := map[string]struct{}{
		"/server/users/login": {},
	}

	server.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper: func(c echo.Context) bool {
			if !strings.HasPrefix(c.Path(), "/server") {
				return true
			}

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

	statik, err := fs.New()
	if err != nil {
		s.services.Logger.Fatal(err)
	}

	// serve static files
	server.Use(static.Static(static.Config{
		Handler: http.FileServer(statik),
	}))

	users := controllers.NewUsersController(s.services)
	posts := controllers.NewPostsController(s.services)
	callabck := controllers.NewCallbackController(s.services)
	integrations := controllers.NewIntegrationsController(s.services)

	server.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	server.POST("/callback/:id", callabck.Callback)

	api := server.Group("api")

	api.POST("/callback/:id", callabck.Callback)

	api.POST("/users/login", users.Login)
	api.POST("/users/update", users.Update)

	api.GET("/posts", posts.List)
	api.POST("/posts/:id", posts.Update)
	api.POST("/posts", posts.Create)
	api.DELETE("/posts/:id", posts.Remove)

	api.GET("/integrations", integrations.List)
	api.POST("/integrations/:id", integrations.Update)
	api.POST("/integrations", integrations.Create)
	api.DELETE("/integrations/:id", integrations.Remove)

	// Start server
	s.services.Logger.Infow("start server server", "port", s.port)
	s.services.Logger.Fatalw("error on starting server", "err", server.Start(":"+s.port))

}
