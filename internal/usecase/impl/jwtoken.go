package impl

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"gophermart/internal/config"
	"gophermart/internal/entity"
	"time"
)

type JwtImpl struct {
	jwtSecret       string
	jwtExpiresHours int
}

func NewJwt(cfg config.JWT) *JwtImpl {
	return &JwtImpl{
		jwtSecret:       cfg.Secret,
		jwtExpiresHours: cfg.ExpiresHours,
	}
}

type Claims struct {
	jwt.RegisteredClaims
	UserID entity.UserID
}

func (ah *JwtImpl) Get(j entity.TokenData) (entity.Token, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(ah.jwtExpiresHours))),
		},
		UserID: j.UserID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(ah.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("PwdImpl - Get - token.SignedString: %w", err)
	}
	return entity.Token(tokenStr), nil
}
