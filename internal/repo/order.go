package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gophermart/internal/entity"
	"gophermart/pkg/postgres"
)

type OrderRepo struct {
	*postgres.Postgres
}

func NewOrder(pg *postgres.Postgres) *OrderRepo {
	return &OrderRepo{pg}
}

func (orr *OrderRepo) Get(ctx context.Context, number entity.OrderNumber) (*entity.Order, error) {
	query := `select user_id, order_nr, status, accrual, created_at from "order" where order_nr = $1`
	o := &entity.Order{}
	err := orr.Pool(ctx).QueryRow(ctx, query, number).Scan(&o.UserID, &o.Number, &o.Status, &o.Accrual, &o.UploadedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrOrderNotFound
		}
		return nil, fmt.Errorf("OrderRepo - Get - orr.Pool(ctx).QueryRow: %w", err)
	}
	return o, nil
}

func (orr *OrderRepo) Add(ctx context.Context, order entity.Order) error {
	// Добавляем новый заказ
	query := `insert into "order" (user_id, order_nr, status, created_at)
        values (@user_id, @order_nr, @status, @created_at)`
	_, err := orr.Pool(ctx).Exec(ctx, query, pgx.NamedArgs{
		"user_id":    order.UserID,
		"order_nr":   order.Number,
		"status":     order.Status,
		"created_at": order.UploadedAt,
	})
	if err != nil {
		return fmt.Errorf("OrderRepo - Add - tx.Exec order: %w", err)
	}
	return nil
}

func (orr *OrderRepo) GetByUser(ctx context.Context, userID entity.UserID) ([]entity.Order, error) {
	query := `select user_id, order_nr, status, case when accrual >  0 then accrual else 0 end, accrual_at, created_at from "order" where user_id = $1`
	rows, _ := orr.Pool(ctx).Query(ctx, query, userID)
	orders, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Order])
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - GetByUser - pgx.CollectRows: %w", err)
	}
	return orders, nil
}

func (orr *OrderRepo) GetUnAccrued(ctx context.Context) ([]entity.OrderAccrual, error) {
	// Берем самые последние сначала, чтобы старые не отодвигали новые
	query := `select order_nr, status
			from "order" where status in ($1, $2)
			order by created_at`
	rows, _ := orr.Pool(ctx).Query(ctx, query, entity.OrderStatusNew, entity.OrderStatusProcessing)
	orders, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.OrderAccrual])
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - GetByUser - pgx.CollectRows: %w", err)
	}
	return orders, nil
}
