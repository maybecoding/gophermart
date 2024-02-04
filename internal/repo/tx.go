package repo

import (
	"context"
	"gophermart/pkg/postgres"
)

type TxRepo struct {
	pg *postgres.Postgres
}

func NewTx(pg *postgres.Postgres) *TxRepo {
	return &TxRepo{pg}
}

func (tr *TxRepo) WithTx(ctx context.Context, fn func(context.Context) error) error {
	return tr.pg.WithTx(ctx, fn)
}
