package repository

import (
	"context"

	"github.com/ak-karimzai/bank-api/internel/db"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	UpdateUser(ctx context.Context, arg db.UpdateUserParams) (db.User, error)
	GetUser(ctx context.Context, username string) (db.User, error)
	CreateUserTx(ctx context.Context, arg db.CreateUserTxParams) (db.CreateUserTxResult, error)
}
