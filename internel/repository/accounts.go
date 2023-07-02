package repository

import (
	"context"

	"github.com/ak-karimzai/bank-api/internel/db"
)

type AccountRepository interface {
	AddAccountBalance(ctx context.Context, arg db.AddAccountBalanceParams) (db.Account, error)
	CreateAccount(ctx context.Context, arg db.CreateAccountParams) (db.Account, error)
	DeleteAccount(ctx context.Context, id int64) error
	GetAccount(ctx context.Context, id int64) (db.Account, error)
	GetAccountForUpdate(ctx context.Context, id int64) (db.Account, error)
	ListAccounts(ctx context.Context, arg db.ListAccountsParams) ([]db.Account, error)
	UpdateAccount(ctx context.Context, arg db.UpdateAccountParams) (db.Account, error)
}
