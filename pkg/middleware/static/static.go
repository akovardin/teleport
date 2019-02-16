package static

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	Config struct {
		Skipper    middleware.Skipper
		Handler    http.Handler
		Extensions ExtSet
		Prefixes   PrefixSet
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
		config.Extensions = ExtSet{".js": {}, ".css": {}, ".ico": {}, ".map": {}, ".png": {}}
	}

	if len(config.Prefixes) == 0 {
		config.Prefixes = PrefixSet{"/api", "/metrics"}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			p := c.Request().URL.Path
			if config.Prefixes.Exists(p) {
				next(c)
				return nil
			}

			if !config.Extensions.Exists(filepath.Ext(p)) && p != "/" {
				c.Request().URL.Path = "/"
			}

			config.Handler.ServeHTTP(c.Response(), c.Request())
			return nil
		}
	}
}

type PrefixSet []string

func (ps PrefixSet) Exists(path string) bool {
	for _, p := range ps {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

type ExtSet map[string]struct{}

func (fs ExtSet) Exists(ext string) bool {
	_, ok := fs[ext]
	return ok
}
