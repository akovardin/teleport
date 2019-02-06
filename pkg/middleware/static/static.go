package static

import (
	"net/http"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	Config struct {
		Skipper    middleware.Skipper
		Handler    http.Handler
		Extensions ExtSet
	}
)

var (
	// DefaultStaticConfig is the default Static middleware config.
	DefaultConfig = Config{
		Skipper: middleware.DefaultSkipper,
	}
)

func Static(config Config) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}

	if len(config.Extensions) == 0 {
		config.Extensions = ExtSet{".js":{}, ".css":{}, ".ico":{}, ".map":{}, ".png":{}}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			p := c.Request().URL.Path


			if !config.Extensions.Exists(filepath.Ext(p)) {
				return next(c)
			}

			config.Handler.ServeHTTP(c.Response(), c.Request())
			return nil
		}
	}
}

type ExtSet map[string]struct{}

func (fs ExtSet) Exists(ext string) bool {
	_, ok := fs[ext]
	return ok
}
