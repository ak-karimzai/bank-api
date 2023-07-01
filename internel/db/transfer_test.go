package db

import (
	"context"
	"testing"
	"time"

	"github.com/ak-karimzai/bank-api/internel/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, acc1, acc2 Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(
		context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.NotZero(t, transfer.ID)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	createRandomTransfer(t, acc1, acc2)
}

func TestGetTransfer(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	transfer1 := createRandomTransfer(t, acc1, acc2)
	transfer2, err := testQueries.
		GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t,
		transfer1.CreatedAt.Time,
		transfer2.CreatedAt.Time,
		time.Second)
}

func TestListTransfers(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	numberOfTransfers := 5
	for i := 0; i < numberOfTransfers; i++ {
		createRandomTransfer(t, acc1, acc2)
	}

	arg := ListTransfersParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Limit:         int32(numberOfTransfers),
		Offset:        0,
	}
	transfers, err := testQueries.
		ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, numberOfTransfers)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
		require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	}
}

func TestDeleteTransfer(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	transfer1 := createRandomTransfer(t, acc1, acc2)
	err := testQueries.DeleteTransfer(
		context.Background(), transfer1.ID)
	require.NoError(t, err)

	transfer2, err := testQueries.
		GetTransfer(context.Background(), transfer1.ID)

	require.Error(t, err)
	require.Empty(t, transfer2)
}
