package usecase

import (
	"context"
	"gophermart/internal/entity"
)

type (
	Auth interface {
		Register(ctx context.Context, login entity.UserLogin, pwd entity.UserPassword) (entity.Token, error)
		Login(ctx context.Context, login entity.UserLogin, pwd entity.UserPassword) (entity.Token, error)
	}

	AuthRepo interface {
		LoginAvailable(ctx context.Context, login entity.UserLogin) (bool, error)
		UserNew(ctx context.Context, login entity.UserLogin, hash entity.UserPasswordHash) (*entity.User, error)
		UserGet(ctx context.Context, login entity.UserLogin) (*entity.User, error)
	}

	Pwd interface {
		Hash(pwd entity.UserPassword) (entity.UserPasswordHash, error)
		IsCorrect(pwd entity.UserPassword, hash entity.UserPasswordHash) bool
	}
	Jwt interface {
		Get(jwt entity.TokenData) (entity.Token, error)
	}
)
