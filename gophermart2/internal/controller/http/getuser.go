package http

import (
	"github.com/gin-gonic/gin"
	"gophermart/internal/entity"
	"net/http"
)

func getUser(c *gin.Context) (userID entity.UserID, ok bool) {
	value, ok := c.Get("UserID")
	if ok {
		userID, ok = value.(entity.UserID)
		if ok {
			return userID, true
		}
	}
	errorResponse(c, entity.ErrUnauthorized, http.StatusUnauthorized)
	return
}
