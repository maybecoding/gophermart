package usecase

import (
	"context"
	"fmt"
	"gophermart/internal/entity"
	"gophermart/pkg/logger"
)

type BonusUseCase struct {
	repo BonusRepo
}

func NewBonus(r BonusRepo) *BonusUseCase {
	return &BonusUseCase{r}
}

func (uc *BonusUseCase) GetBalance(ctx context.Context, userID entity.UserID) (*entity.BonusBalance, error) {
	balance, err := uc.repo.GetBalance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("BonusUseCase - BonusUseCase - uc.repo.GetBalance: %w", err)
	}
	logger.Debug().Interface("balance", balance).Msg("BonusUseCase - GetBalance")
	return balance, err
}

func (uc *BonusUseCase) GetWithdrawals(ctx context.Context, userID entity.UserID) ([]entity.BonusWithdraw, error) {
	withdrawals, err := uc.repo.GetWithdrawals(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("BonusUseCase - GetWithdrawals - uc.repo.GetWithdrawals: %w", err)
	}
	logger.Debug().Interface("withdrawals", withdrawals).Msg("BonusUseCase - GetWithdrawals")
	return withdrawals, err
}
