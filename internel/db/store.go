package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(
	ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

type TransferTxParams struct {
	FromAccountID int64  `json:"from_account_id"`
	ToAccountID   int64  `json:"to_account_id"`
	Amount        int64  `json:"amount"`
	Username      string `json:"sender_username"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *SQLStore) TransferTx(
	ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		acc1, err := store.GetAccount(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		acc2, err := store.GetAccount(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		if err = validTransaction(&acc1, &arg); err != nil {
			return err
		}

		if acc1.Currency != acc2.Currency {
			return fmt.Errorf(
				"account [%d] currency mismatch: %s vs %s",
				acc1.ID, acc1.Currency, acc2.Currency)
		}

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(
				ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.FromAccount, result.ToAccount, err = addMoney(
				ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		return err
	})
	return result, err
}

func validTransaction(
	acc *Account, arg *TransferTxParams) error {
	if acc.Owner != arg.Username {
		return errors.New("account doesn't belong to the sender")
	}
	if acc.Balance-arg.Amount < 0 {
		return fmt.Errorf("no enough balance")
	}
	return nil
}

func addMoney(ctx context.Context,
	q *Queries, accId1, amount1, accId2, amount2 int64) (acc1, acc2 Account, err error) {
	acc1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accId1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	acc2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accId2,
		Amount: amount2,
	})
	if err != nil {
		return
	}
	return
}
