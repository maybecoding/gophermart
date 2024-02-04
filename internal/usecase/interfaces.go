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
		LoginAvailable(ctx context.Context, login entity.UserLogin) (loginAvail bool, err error)
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
		Add(ctx context.Context, userID entity.UserID, number entity.OrderNumber) error
		AddForBonuses(ctx context.Context, userID entity.UserID, number entity.OrderNumber, amount entity.BonusAmount) error
		GetByUser(ctx context.Context, userID entity.UserID) ([]entity.Order, error)
		RunAccrualRefresh(ctx context.Context)
	}

	OrderRepo interface {
		Get(ctx context.Context, number entity.OrderNumber) (*entity.Order, error)
		Add(ctx context.Context, order entity.Order) error
		GetByUser(ctx context.Context, userID entity.UserID) ([]entity.Order, error)
		GetUnAccrued(ctx context.Context) ([]entity.OrderAccrual, error)
	}

	OrderNumAlg interface {
		Check(num entity.OrderNumber) (isCorrect bool, err error)
	}
	OrderAccrual interface {
		GetStatus(orderNum entity.OrderNumber) (*entity.AccrualInfo, error)
	}
)

type TxRepo interface {
	WithTx(ctx context.Context, fn func(context.Context) error) error
}

type (
	Bonus interface {
		Get(ctx context.Context, userID entity.UserID) (*entity.BonusBalance, error)
		GetWithdrawals(ctx context.Context, userID entity.UserID) ([]entity.BonusWithdraw, error)
	}
	BonusRepo interface {
		GetWithdrawals(ctx context.Context, userID entity.UserID) ([]entity.BonusWithdraw, error)
		Set(ctx context.Context, order entity.OrderNumber, amount entity.BonusAmount, status entity.OrderStatus) error
	}
	BalanceRepo interface {
		Get(ctx context.Context, userID entity.UserID) (*entity.BonusBalance, error)
		Set(ctx context.Context, userID entity.UserID, balance entity.BonusBalance) error
	}
)
