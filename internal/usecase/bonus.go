package usecase

import (
	"context"
	"fmt"
	"gophermart/internal/entity"
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
	return balance, err
}

func (uc *BonusUseCase) GetWithdrawals(ctx context.Context, userID entity.UserID) ([]entity.BonusWithdraw, error) {
	withdrawals, err := uc.repo.GetWithdrawals(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("BonusUseCase - GetWithdrawals - uc.repo.GetWithdrawals: %w", err)
	}
	return withdrawals, err
}
