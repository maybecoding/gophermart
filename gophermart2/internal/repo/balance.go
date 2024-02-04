package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gophermart/internal/entity"
	"gophermart/pkg/postgres"
)

type BalanceRepo struct {
	*postgres.Postgres
}

func NewBalance(pg *postgres.Postgres) *BalanceRepo {
	return &BalanceRepo{pg}
}

func (br *BalanceRepo) Get(ctx context.Context, userID entity.UserID) (*entity.BonusBalance, error) {
	balance := entity.BonusBalance{}
	err := br.Pool(ctx).QueryRow(ctx, `select available, withdrawn from balance where user_id = @user_id`, pgx.NamedArgs{
		"user_id": userID,
	}).Scan(&balance.Available, &balance.Withdrawn)
	if err != nil {
		return nil, fmt.Errorf("BalanceRepo - Get - br.Pool(ctx).QueryRow: %w", err)
	}
	return &balance, nil
}

func (br *BalanceRepo) Set(ctx context.Context, userID entity.UserID, balance entity.BonusBalance) error {
	_, err := br.Pool(ctx).Exec(ctx,
		`update balance set available = @available, withdrawn = @withdrawn where user_id = @user_id`,
		pgx.NamedArgs{
			"user_id":   userID,
			"available": balance.Available,
			"withdrawn": balance.Withdrawn,
		})
	if err != nil {
		return fmt.Errorf("BalanceRepo - Set - br.Pool(ctx).Exec: %w", err)
	}
	return nil
}
