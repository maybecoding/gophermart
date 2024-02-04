package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gophermart/internal/entity"
	"gophermart/pkg/postgres"
	"time"
)

type BonusRepo struct {
	*postgres.Postgres
}

func NewBonus(pg *postgres.Postgres) *BonusRepo {
	return &BonusRepo{pg}
}

func (bns *BonusRepo) GetWithdrawals(ctx context.Context, userID entity.UserID) ([]entity.BonusWithdraw, error) {
	query := `select -accrual, order_nr, accrual_at
from "order"
where user_id = @user_id
    and accrual < 0
    and accrual_at is not null;`

	rows, _ := bns.Pool(ctx).Query(ctx, query, pgx.NamedArgs{"user_id": userID})
	withdrawals, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.BonusWithdraw])
	if err != nil {
		return nil, fmt.Errorf("BonusRepo - GetWithdrawals - bns.Pool(ctx).Query: %w", err)
	}
	return withdrawals, nil
}

func (bns *BonusRepo) Set(ctx context.Context, order entity.OrderNumber, amount entity.BonusAmount, status entity.OrderStatus) error {
	_, err := bns.Pool(ctx).Exec(ctx,
		`update "order" set accrual = @accrual, accrual_at = @accrual_at, status = @status where order_nr = @order_nr`,
		pgx.NamedArgs{
			"order_nr":   order,
			"accrual":    amount,
			"status":     status,
			"accrual_at": time.Now(),
		})
	if err != nil {
		return fmt.Errorf("BonusRepo - Set - bns.Pool(ctx).Exec: %w", err)
	}
	return nil
}
