package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gophermart/internal/entity"
	"gophermart/pkg/logger"
	"gophermart/pkg/postgres"
	"time"
)

type OrderRepo struct {
	*postgres.Postgres
}

func NewOrder(pg *postgres.Postgres) *OrderRepo {
	return &OrderRepo{pg}
}

func (orr *OrderRepo) Get(ctx context.Context, number entity.OrderNumber) (*entity.Order, error) {
	query := `select user_id, number, status, created_at from user_order where number = $1`
	o := &entity.Order{}
	err := orr.Pool.QueryRow(ctx, query, number).Scan(&o.UserID, &o.Number, &o.Status, &o.UploadedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrOrderNotFound
		}
		return nil, fmt.Errorf("OrderRepo - Get - orr.Pool.QueryRow: %w", err)
	}
	return o, nil
}

func (orr *OrderRepo) Add(ctx context.Context, order entity.Order) (*entity.Order, error) {
	query := `with us_or as (
    insert into user_order (user_id, number, status, accrual, accrual_at, created_at)
        values (@user_id, @number, @status, @accrual, @accrual_at, @created_at)
        returning user_id, number, status, accrual, accrual_at, created_at
), _ as (
    update user_bonus_balance as bal set available = available + us_or.accrual
        from us_or
        where us_or.user_id = bal.user_id
)
select user_id, number, status, accrual, accrual_at, created_at from us_or;`
	o := &entity.Order{}
	var resAccrualAt sql.NullTime
	err := orr.Pool.QueryRow(ctx, query, pgx.NamedArgs{
		"user_id":    order.UserID,
		"number":     order.Number,
		"status":     order.Status,
		"accrual":    order.Accrual,
		"accrual_at": order.AccrualProcessedAt,
		"created_at": order.UploadedAt,
	}).
		Scan(&o.UserID, &o.Number, &o.Status, &o.Accrual, &resAccrualAt, &o.UploadedAt)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - Add - orr.Pool.QueryRow: %w", err)
	}
	if resAccrualAt.Valid {
		o.AccrualProcessedAt = &resAccrualAt.Time
	}
	return o, nil
}

func (orr *OrderRepo) GetByUser(ctx context.Context, userID entity.UserID) ([]entity.Order, error) {
	query := `select user_id, number, status, case when accrual >  0 then accrual else 0 end, accrual_at, created_at from user_order where user_id = $1`
	rows, _ := orr.Pool.Query(ctx, query, userID)
	orders, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Order])
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - GetByUser - pgx.CollectRows: %w", err)
	}
	return orders, nil
}

func (orr *OrderRepo) GetUnAccrued(ctx context.Context) ([]entity.Order, error) {
	// Берем самые последние сначала, чтобы старые не отодвигали новые
	query := `select user_id, number, status, accrual, accrual_at, created_at
			from user_order where status in ($1, $2)
			order by created_at`
	rows, _ := orr.Pool.Query(ctx, query, entity.OrderStatusNew, entity.OrderStatusProcessing)
	orders, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Order])
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - GetByUser - pgx.CollectRows: %w", err)
	}
	return orders, nil
}

func (orr *OrderRepo) Accrual(ctx context.Context, accrual entity.AccrualInfo) (*entity.Order, error) {
	query := `with us_or as (
    select user_id, number, status, accrual, accrual_at, created_at
    from user_order
    where number = @number
), us_or_new as (
    update user_order set accrual = @accrual, accrual_at = @accrual_at, status = @status
    where number = @number
    returning user_id, number, status, accrual, accrual_at, created_at
), _ as (
    update user_bonus_balance as bal set available = available - us_or.accrual + us_or_new.accrual
        from us_or, us_or_new
        where us_or.user_id = bal.user_id
)
select user_id, number, status, accrual, accrual_at, created_at from us_or_new;`
	o := &entity.Order{}
	var resAccrualAt sql.NullTime
	logger.Debug().Interface("AccrualInfo", accrual).Msg("accrual order")
	err := orr.Pool.QueryRow(ctx, query, pgx.NamedArgs{
		"number":     accrual.Order,
		"status":     accrual.Status,
		"accrual":    accrual.Accrual,
		"accrual_at": time.Now(),
	}).
		Scan(&o.UserID, &o.Number, &o.Status, &o.Accrual, &resAccrualAt, &o.UploadedAt)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - Accrual - orr.Pool.QueryRow: %w", err)
	}
	if resAccrualAt.Valid {
		o.AccrualProcessedAt = &resAccrualAt.Time
	}
	return o, nil
}
