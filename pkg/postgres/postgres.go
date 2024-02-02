package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func New(uri string) (*Postgres, error) {
	pg := &Postgres{}
	var err error
	pg.Pool, err = pgxpool.New(context.Background(), uri)
	if err != nil {
		return nil, fmt.Errorf("postgres - New - pgxpool.New: %w", err)
	}
	return pg, nil
}
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
