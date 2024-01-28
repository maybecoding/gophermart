package impl

import (
	"errors"
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

type claims struct {
	jwt.RegisteredClaims
	UserID entity.UserID
}

func (ah *JwtImpl) Get(j entity.TokenData) (entity.Token, error) {
	c := claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(ah.jwtExpiresHours))),
		},
		UserID: j.UserID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	tokenStr, err := token.SignedString([]byte(ah.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("PwdImpl - Get - token.SignedString: %w", err)
	}
	return entity.Token(tokenStr), nil
}

func (ah *JwtImpl) Use(t entity.Token) (entity.TokenData, error) {
	td := entity.TokenData{}
	c := &claims{}
	token, err := jwt.ParseWithClaims(string(t), c, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ah.jwtSecret), nil
	})
	if err != nil {
		return td, fmt.Errorf("JwtImpl - Use - jwt.Parse: %w", err)
	}

	if !token.Valid {
		return td, errors.New("JwtImpl - Use - token is not valid")
	}
	td.UserID = c.UserID
	return td, nil
}
