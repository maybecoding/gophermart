package usecase

import (
	"context"
	"gophermart/internal/entity"
)

type (
	Auth interface {
		Register(ctx context.Context, login entity.UserLogin, pwd entity.UserPassword) (entity.Token, error)
		Login(ctx context.Context, login entity.UserLogin, pwd entity.UserPassword) (entity.Token, error)
		ValidateToken(t entity.Token) (entity.TokenData, error)
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
		Use(t entity.Token) (entity.TokenData, error)
	}
)

type (
	Order interface {
		AddNew(ctx context.Context, userID entity.UserID, number entity.OrderNumber) (*entity.Order, error)
		AddForBonuses(ctx context.Context, userID entity.UserID, number entity.OrderNumber, amount entity.BonusAmount) (*entity.Order, error)
		GetByUser(ctx context.Context, userID entity.UserID) ([]entity.Order, error)
		RunAccrualRefresh(ctx context.Context)
	}

	OrderRepo interface {
		Add(ctx context.Context, order entity.Order) (*entity.Order, error)
		Get(ctx context.Context, number entity.OrderNumber) (*entity.Order, error)
		GetByUser(ctx context.Context, userID entity.UserID) ([]entity.Order, error)
		GetUnAccrued(ctx context.Context) ([]entity.Order, error)
		Accrual(ctx context.Context, accrual entity.AccrualInfo) (*entity.Order, error)
	}

	OrderNumAlg interface {
		Check(num entity.OrderNumber) (isCorrect bool, err error)
	}
	OrderAccrual interface {
		GetStatus(orderNum entity.OrderNumber) (*entity.AccrualInfo, error)
	}
)

type (
	Bonus interface {
		GetBalance(ctx context.Context, userID entity.UserID) (*entity.BonusBalance, error)
		GetWithdrawals(ctx context.Context, userID entity.UserID) ([]entity.BonusWithdraw, error)
	}
	BonusRepo interface {
		GetBalance(ctx context.Context, userID entity.UserID) (*entity.BonusBalance, error)
		GetWithdrawals(ctx context.Context, userID entity.UserID) ([]entity.BonusWithdraw, error)
	}
)
