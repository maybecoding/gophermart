package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"gophermart/internal/entity"
	"gophermart/pkg/postgres"
)

type AuthRepo struct {
	*postgres.Postgres
}

func NewAuth(pg *postgres.Postgres) *AuthRepo {
	return &AuthRepo{pg}
}

func (ar *AuthRepo) LoginAvailable(ctx context.Context, login entity.UserLogin) (loginAvail bool, err error) {
	query := `select not exists (select 1 from usr where login = $1) login_available;`
	err = ar.Pool.QueryRow(ctx, query, login).Scan(&loginAvail)
	if err != nil {
		return false, fmt.Errorf("AuthRepo - LoginAvailable - ar.Pool.QueryRow: %w", err)
	}
	return loginAvail, nil
}
func (ar *AuthRepo) UserNew(ctx context.Context, login entity.UserLogin, hash entity.UserPasswordHash) (*entity.User, error) {
	usr := entity.User{}
	query := `insert into usr(login, hash) values($1, $2) returning id, login, hash`
	err := ar.Pool.QueryRow(ctx, query, login, hash).Scan(&usr.UserID, &usr.UserLogin, &usr.UserPasswordHash)
	if err != nil {
		return nil, fmt.Errorf("AuthRepo - UserNew - ar.Pool.QueryRow: %w", err)
	}
	return &usr, nil
}
func (ar *AuthRepo) UserGet(ctx context.Context, login entity.UserLogin) (*entity.User, error) {
	usr := entity.User{}
	query := `select id, login, hash from  usr where login = $1`
	err := ar.Pool.QueryRow(ctx, query, login).Scan(&usr.UserID, &usr.UserLogin, &usr.UserPasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrUserNotFound
		}
		return nil, fmt.Errorf("AuthRepo - UserNew - ar.Pool.QueryRow: %w", err)
	}
	return &usr, nil
}
