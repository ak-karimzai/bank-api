package db

import (
	"context"
	"testing"

	"github.com/ak-karimzai/bank-api/internel/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, acc Account) Entry {
	arg := CreateEntryParams{
		AccountID: acc.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.
		CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.NotZero(t, entry.ID)

	require.Equal(
		t, arg.Amount, entry.Amount)
	require.Equal(
		t, arg.AccountID, entry.AccountID)
	return entry
}

func TestCreateEntry(t *testing.T) {
	acc := createRandomAccount(t)

	createRandomEntry(t, acc)
}

func TestGetEntry(t *testing.T) {
	acc := createRandomAccount(t)

	entry1 := createRandomEntry(t, acc)
	entry2, err := testQueries.GetEntry(
		context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
}

func TestListEntries(t *testing.T) {
	acc := createRandomAccount(t)

	numberOfEntries := 5

	for i := 0; i < numberOfEntries; i++ {
		createRandomEntry(t, acc)
	}

	arg := ListEntriesParams{
		AccountID: acc.ID,
		Limit:     int32(numberOfEntries),
		Offset:    0,
	}
	entries, err := testQueries.
		ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, numberOfEntries)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, entry.AccountID, arg.AccountID)
	}
}

func TestUpdateEntry(t *testing.T) {
	acc := createRandomAccount(t)

	entry1 := createRandomEntry(t, acc)

	arg := UpdateEntryParams{
		Amount: util.RandomMoney(),
		ID:     entry1.ID,
	}
	entry2, err := testQueries.UpdateEntry(
		context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, arg.Amount, entry2.Amount)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
}

func TestDeleteEntry(t *testing.T) {
	acc := createRandomAccount(t)

	entry1 := createRandomEntry(t, acc)

	err := testQueries.DeleteEntry(
		context.Background(), entry1.ID)
	require.NoError(t, err)

	entry2, err := testQueries.GetEntry(
		context.Background(), entry1.ID)
	require.Error(t, err)
	require.Empty(t, entry2)
}
