package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserInfo(c *gin.Context) {
	ui, exists := c.Get("UserInfo")
	if exists {
		c.JSON(http.StatusOK, ui)
	} else {
		c.String(http.StatusOK, "{\"name\": \"Unknown\"}")
	}
}
