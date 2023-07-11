package repository

import (
	"context"

	"github.com/ak-karimzai/bank-api/internel/db"
)

type TransferRepository interface {
	TransferTx(ctx context.Context, arg db.TransferTxParams) (db.TransferTxResult, error)
	CreateTransfer(ctx context.Context, arg db.CreateTransferParams) (db.Transfer, error)
	DeleteTransfer(ctx context.Context, id int64) error
	GetTransfer(ctx context.Context, id int64) (db.Transfer, error)
	ListTransfers(ctx context.Context, arg db.ListTransfersParams) ([]db.Transfer, error)
	UpdateTransfer(ctx context.Context, arg db.UpdateTransferParams) (db.Transfer, error)
}
