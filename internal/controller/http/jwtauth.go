package http

import (
	"github.com/gin-gonic/gin"
	"gophermart/internal/entity"
	"gophermart/internal/usecase"
	"gophermart/pkg/logger"
	"strings"
)

func JWTAuth(uc usecase.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		authCookie, err := c.Cookie("Authorization")
		if err != nil {
			logger.Error().Err(err).Msg("http - JWTAuth - c.Cookie")
			c.Next()
			return
		}
		logger.Debug().Str("Authorization", authCookie).Msg("http - JWTAuth - c.Cookie")
		auth := strings.Split(authCookie, " ")
		if len(auth) != 2 || auth[0] != "Bearer" {
			c.Next()
			return
		}
		token := entity.Token(auth[1])

		tokenData, err := uc.ValidateToken(token)
		if err != nil {
			logger.Error().Err(err).Msg("http - JWTAuth - uc.ValidateToken")
			c.Next()
			return
		}
		c.Set("UserID", tokenData.UserID)
		logger.Debug().Int32("UserID", int32(tokenData.UserID)).Msg("JWTAuth")
		c.Next()
	}
}
