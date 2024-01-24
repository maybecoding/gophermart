package http

import (
	"github.com/gin-gonic/gin"
	"gophermart/internal/entity"
)

func setAuthCookie(c *gin.Context, token entity.Token) {
	c.SetCookie("Authorization", "Bearer "+string(token), 3600, "/", "localhost", false, true)
}
