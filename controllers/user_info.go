package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserInfo(c *gin.Context) {
	c.String(http.StatusOK, c.Request.Header.Get("X-USER"))
}
