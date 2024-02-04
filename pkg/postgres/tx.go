package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"gophermart/pkg/logger"
)

// DB - интерфейс, содержащий методы имеющиеся и в пуле соединений и в транзакции
// Таким образом мы работаем с этим интерфейсом для выполнения операций с БД, но не знаем выполняем действия в транзакции или нет,
// поскольку это не ответственность кода выполняющего запросы, а уровня выше который решает для каких операций необходимы транзакции
type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

var txKey = struct{}{}

// injectTx - Внедряет транзакцию в контекст
func (pg *Postgres) injectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

// extractTx - Извлекает транзакцию из контекста если она есть
func (pg *Postgres) extractTx(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(txKey).(pgx.Tx); ok {
		return tx
	}
	return nil
}

// Pool - Принимает контекст и возвращает транзакцию, если она встроена в передаваемый контекст, иначе пул соединений.
// Возвращается интерфейс, содержащий методы доступные как у транзакции так и у пула соединений
func (pg *Postgres) Pool(ctx context.Context) DB {
	tx := pg.extractTx(ctx)

	if tx != nil {
		return tx
	}
	return pg.pool

}

// WithTx - Принимает контекст и функцию принимающую контекст, начинает транзакцию и внедряет ее в новый контекст, который передается функции принятой в качестве параметра. Таким образом функция получит контекст, в который внедрена транзакция.
// Таким образом когда вызывающий код будет получать пул соединений функцией Pool он получит транзакцию.
// То есть внутри передаваемой функции можно вызвать несколько операций, которые будут работать в одной транзакции.
func (pg *Postgres) WithTx(ctx context.Context, fn func(context.Context) error) error {
	tx, err := pg.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})

	if err != nil {
		return fmt.Errorf("postgres - WithTx - pg.pool.BeginTx: %w", err)
	}

	err = fn(pg.injectTx(ctx, tx))
	if err != nil {
		rErr := tx.Rollback(ctx)
		if rErr != nil {
			logger.Error().Err(rErr).Msg("postgres - WithTx - defer tx.Rollback")
		}
		logger.Error().Err(rErr).Msg("postgres - WithTx - fn")
		return fmt.Errorf("postgres - WithTx - fn: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("postgres - WithTx - tx.Commit: %w", err)
	}
	return nil
}
