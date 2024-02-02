package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"gophermart/internal/config"
	"gophermart/internal/entity"
	"time"
)

type Jwt struct {
	jwtSecret       string
	jwtExpiresHours int
}

func New(cfg config.JWT) *Jwt {
	return &Jwt{
		jwtSecret:       cfg.Secret,
		jwtExpiresHours: cfg.ExpiresHours,
	}
}

type claims struct {
	jwt.RegisteredClaims
	UserID entity.UserID
}

func (ah *Jwt) Get(j entity.TokenData) (entity.Token, error) {
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

func (ah *Jwt) Use(t entity.Token) (entity.TokenData, error) {
	td := entity.TokenData{}
	c := &claims{}
	token, err := jwt.ParseWithClaims(string(t), c, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ah.jwtSecret), nil
	})
	if err != nil {
		return td, fmt.Errorf("Jwt - Use - jwt.Parse: %w", err)
	}

	if !token.Valid {
		return td, errors.New("Jwt - Use - token is not valid")
	}
	td.UserID = c.UserID
	return td, nil
}
