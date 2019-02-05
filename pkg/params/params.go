package params

import (
	"strconv"

	"github.com/labstack/echo"
)

func GetInt(c echo.Context, name string, def int) int {
	result, _ := strconv.Atoi(c.QueryParam(name))
	if result == 0 {
		return def
	}
	return result
}

func GetString(c echo.Context, name string, def string) string {
	result := c.QueryParam(name)
	if result == "" {
		return def
	}
	return result
}
