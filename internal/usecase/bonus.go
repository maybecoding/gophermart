package usecase

import (
	"context"
	"fmt"
	"gophermart/internal/entity"
	"gophermart/pkg/logger"
)

type BonusUseCase struct {
	balance BalanceRepo
	bonus   BonusRepo
}

func NewBonus(bal BalanceRepo, bns BonusRepo) *BonusUseCase {
	return &BonusUseCase{bal, bns}
}

func (uc *BonusUseCase) Get(ctx context.Context, userID entity.UserID) (*entity.BonusBalance, error) {
	balance, err := uc.balance.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("BonusUseCase - BonusUseCase - uc.order.Get: %w", err)
	}
	logger.Debug().Interface("balance", balance).Msg("BonusUseCase - Get")
	return balance, err
}

func (uc *BonusUseCase) GetWithdrawals(ctx context.Context, userID entity.UserID) ([]entity.BonusWithdraw, error) {
	withdrawals, err := uc.bonus.GetWithdrawals(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("BonusUseCase - GetWithdrawals - uc.order.GetWithdrawals: %w", err)
	}
	logger.Debug().Interface("withdrawals", withdrawals).Msg("BonusUseCase - GetWithdrawals")
	return withdrawals, err
}
