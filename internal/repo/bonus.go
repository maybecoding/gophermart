package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gophermart/internal/entity"
	"gophermart/pkg/postgres"
)

type BonusRepo struct {
	*postgres.Postgres
}

func NewBonus(pg *postgres.Postgres) *BonusRepo {
	return &BonusRepo{pg}
}

//
//

func (br *BonusRepo) GetBalance(ctx context.Context, userID entity.UserID) (*entity.BonusBalance, error) {
	query := `select available, withdrawn from balance where user_id = @user_id`
	balance := entity.BonusBalance{}
	err := br.Pool.QueryRow(ctx, query, pgx.NamedArgs{
		"user_id": userID,
	}).Scan(&balance.Available, &balance.Withdrawn)
	if err != nil {
		return nil, fmt.Errorf("BonusRepo - GetBalance - br.Pool.QueryRow: %w", err)
	}
	return &balance, nil
}
func (br *BonusRepo) GetWithdrawals(ctx context.Context, userID entity.UserID) ([]entity.BonusWithdraw, error) {
	query := `select -accrual, order_nr, accrual_at
from "order"
where user_id = @user_id
    and accrual < 0
    and accrual_at is not null;`

	rows, _ := br.Pool.Query(ctx, query, pgx.NamedArgs{"user_id": userID})
	withdrawals, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.BonusWithdraw])
	if err != nil {
		return nil, fmt.Errorf("BonusRepo - GetBalance - br.Pool.Query: %w", err)
	}
	return withdrawals, nil
}
