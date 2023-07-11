package db

import (
	"context"
	"testing"
	"time"

	"github.com/ak-karimzai/bank-api/internel/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T, currency ...string) Account {
	selectedCurrency := util.RandomCurrency()
	if len(currency) > 0 {
		selectedCurrency = currency[0]
	}

	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: Currency(selectedCurrency),
	}

	account, err := testQueries.
		CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account.Owner, arg.Owner)
	require.Equal(t, account.Balance, arg.Balance)
	require.Equal(t, account.Currency, arg.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2, err := testQueries.
		GetAccount(context.Background(), acc1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.Equal(t, acc1.Balance, acc2.Balance)
	require.WithinDuration(t,
		acc1.CreatedAt,
		acc2.CreatedAt,
		time.Second)
}

func TestUpdateAccount(t *testing.T) {
	acc1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      acc1.ID,
		Balance: util.RandomMoney(),
	}

	acc2, err := testQueries.
		UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, arg.Balance, acc2.Balance)
}

func TestDeleteAccount(t *testing.T) {
	acc1 := createRandomAccount(t)

	err := testQueries.
		DeleteAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	acc2, err := testQueries.
		GetAccount(context.Background(), acc1.ID)
	require.Error(t, err)
	require.Empty(t, acc2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 0,
		Owner:  lastAccount.Owner,
	}

	accounts, err := testQueries.
		ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
