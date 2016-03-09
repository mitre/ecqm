package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

func UserInfo(c *echo.Context) error {
	return c.String(http.StatusOK, c.Request().Header.Get("X-USER"))
}
