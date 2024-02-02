package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gophermart/internal/entity"
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
	query := `select user_id, order_nr, status, created_at from "order" where order_nr = $1`
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

func (orr *OrderRepo) Add(ctx context.Context, order entity.Order) error {
	tx, err := orr.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return fmt.Errorf("OrderRepo - Add - orr.Pool.BeginTx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Добавляем новый заказ
	query := `insert into "order" (user_id, order_nr, status, accrual, accrual_at, created_at)
        values (@user_id, @order_nr, @status, @accrual, @accrual_at, @created_at)`
	_, err = tx.Exec(ctx, query, pgx.NamedArgs{
		"user_id":    order.UserID,
		"order_nr":   order.Number,
		"status":     order.Status,
		"accrual":    order.Accrual,
		"accrual_at": order.AccrualProcessedAt,
		"created_at": order.UploadedAt,
	})
	if err != nil {
		return fmt.Errorf("OrderRepo - Add - tx.Exec order: %w", err)
	}

	// Обновляем баланс
	var withdrawn entity.BonusAmount
	if order.Accrual < 0 {
		withdrawn = -order.Accrual
	}
	query = `update balance set
		available = available + @available, withdrawn = withdrawn + @withdrawn
    	where user_id = @user_id`
	_, err = tx.Exec(ctx, query, pgx.NamedArgs{
		"available": order.Accrual,
		"withdrawn": withdrawn,
		"user_id":   order.UserID,
	})
	if err != nil {
		return fmt.Errorf("OrderRepo - Add - tx.Exec balance: %w", err)
	}

	// Завершаем транзакцию
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("OrderRepo - Add - tx.Commit: %w", err)
	}
	return nil
}

func (orr *OrderRepo) GetByUser(ctx context.Context, userID entity.UserID) ([]entity.Order, error) {
	query := `select user_id, order_nr, status, case when accrual >  0 then accrual else 0 end, accrual_at, created_at from "order" where user_id = $1`
	rows, _ := orr.Pool.Query(ctx, query, userID)
	orders, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Order])
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - GetByUser - pgx.CollectRows: %w", err)
	}
	return orders, nil
}

func (orr *OrderRepo) GetUnAccrued(ctx context.Context) ([]entity.Order, error) {
	// Берем самые последние сначала, чтобы старые не отодвигали новые
	query := `select user_id, order_nr, status, accrual, accrual_at, created_at
			from "order" where status in ($1, $2)
			order by created_at`
	rows, _ := orr.Pool.Query(ctx, query, entity.OrderStatusNew, entity.OrderStatusProcessing)
	orders, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Order])
	if err != nil {
		return nil, fmt.Errorf("OrderRepo - GetByUser - pgx.CollectRows: %w", err)
	}
	return orders, nil
}

func (orr *OrderRepo) Accrual(ctx context.Context, accrual entity.AccrualInfo) error {
	tx, err := orr.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return fmt.Errorf("OrderRepo - Accrual - orr.Pool.BeginTx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// Получаем предыдущее значение accrual и user_id
	var lastAccrual entity.BonusAmount
	var userID entity.UserID
	err = tx.QueryRow(ctx, `select accrual, user_id from "order" where order_nr = $1`, accrual.Order).Scan(&lastAccrual, &userID)
	if err != nil {
		return fmt.Errorf("OrderRepo - Accrual - tx.QueryRow lastAccrual, userID: %w", err)
	}

	// Изменяем заказ
	_, err = tx.Exec(ctx,
		`update "order" set accrual = @accrual, accrual_at = @accrual_at, status = @status where order_nr = @order_nr`,
		pgx.NamedArgs{
			"order_nr":   accrual.Order,
			"status":     accrual.Status,
			"accrual":    accrual.Accrual,
			"accrual_at": time.Now(),
		})
	if err != nil {
		return fmt.Errorf("OrderRepo - Accrual - tx.Exec update order: %w", err)
	}

	// Корректируем бонусы
	// Если было какое-то начисление, а это корректирующее учитываем это (но в общем случае lastAccrual = 0)
	deltaAvailable := accrual.Accrual - lastAccrual
	_, err = tx.Exec(ctx, `update balance set available = available + @available where user_id = @user_id`,
		pgx.NamedArgs{
			"available": deltaAvailable,
			"user_id":   userID,
		})
	if err != nil {
		return fmt.Errorf("OrderRepo - Accrual - tx.Exec update balance: %w", err)
	}

	// Завершаем транзакцию
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("OrderRepo - Accrual - tx.Commit: %w", err)
	}
	return nil
}
