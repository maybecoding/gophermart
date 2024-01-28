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
	query := `insert into user_order(user_id, number, status, created_at)
	values($1, $2, $3, $4)
	returning user_id, number, status, created_at`
	o := &entity.Order{}
	err := orr.Pool.QueryRow(ctx, query, order.UserID, order.Number, order.Status, order.UploadedAt).
		Scan(&o.UserID, &o.Number, &o.Status, &o.UploadedAt)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - Add - orr.Pool.QueryRow: %w", err)
	}
	return o, nil
}

func (orr *OrderRepo) GetByUser(ctx context.Context, userID entity.UserID) ([]entity.Order, error) {
	query := `select user_id, number, status, accrual, created_at from user_order where user_id = $1`
	rows, _ := orr.Pool.Query(ctx, query, userID)
	orders, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Order])
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - Add - orr.Pool.QueryRow: %w", err)
	}
	return orders, nil
}
