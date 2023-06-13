package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ak-karimzai/ak-karimzai/simpleb/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T, user User) Account {
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	acc, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, acc)
	require.Equal(t, arg.Owner, acc.Owner)
	require.Equal(t, arg.Balance, acc.Balance)
	require.Equal(t, arg.Currency, acc.Currency)

	require.NotZero(t, acc.ID)
	require.NotZero(t, acc.CreatedAt)
	return acc
}

func TestCreateAccount(t *testing.T) {
	user := createRandomUser(t)
	createRandomAccount(t, user)
}

func TestGetAccount(t *testing.T) {
	user := createRandomUser(t)
	acc1 := createRandomAccount(t, user)
	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc2.ID, acc1.ID)
	require.Equal(t, acc2.Owner, acc1.Owner)
	require.Equal(t, acc2.Balance, acc1.Balance)
	require.Equal(t, acc2.Currency, acc1.Currency)
	require.WithinDuration(t, acc2.CreatedAt, acc1.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	user := createRandomUser(t)
	acc1 := createRandomAccount(t, user)

	arg := UpdateAccountParams{
		ID:      acc1.ID,
		Balance: util.RandomMoney(),
	}
	acc2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, acc2)
	require.Equal(t, acc2.ID, acc1.ID)
	require.Equal(t, acc2.Owner, acc1.Owner)
	require.Equal(t, acc2.Balance, arg.Balance)
	require.Equal(t, acc2.Currency, acc1.Currency)
	require.WithinDuration(t, acc2.CreatedAt, acc1.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	user := createRandomUser(t)
	acc1 := createRandomAccount(t, user)
	err := testQueries.DeleteAccount(
		context.Background(), acc1.ID)

	require.NoError(t, err)
	acc2, err := testQueries.GetAccount(
		context.Background(), acc1.ID)

	require.Error(t, err)
	require.Empty(t, acc2)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestListAccounts(t *testing.T) {
	users := make([]User, 10)
	for i := 0; i < 10; i++ {
		users[i] = createRandomUser(t)
		require.NotEmpty(t, users[i])
		createRandomAccount(t, users[i])
	}

	arg := GetAccountsParams{
		Owner:  users[0].Username,
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.GetAccounts(
		context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, acc := range accounts {
		require.NotEmpty(t, acc)
	}
}
