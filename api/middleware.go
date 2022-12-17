package api

import (
	"github.com/labstack/echo/v4"
)

func skipper(c echo.Context) bool {
	if c.Request().Method == "POST" {
		switch c.Path() {
		case "/user":
			return true
		case "/login":
			return true
		case "/refresh_token":
			return true
		}
	}
	return false
}
