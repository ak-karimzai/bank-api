package db

import (
	"context"
	"testing"
	"time"

	"github.com/ak-karimzai/bank-api/internel/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPwd, err := util.HashPasswrod(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:  util.RandomOwner(),
		HashedPwd: hashedPwd,
		FullName:  util.RandomOwner(),
		Email:     util.RandomEmail(),
	}

	user, err := testQueries.
		CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, user.Username, arg.Username)
	require.Equal(t, user.HashedPwd, arg.HashedPwd)
	require.Equal(t, user.FullName, arg.FullName)
	require.Equal(t, user.Email, arg.Email)

	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	user2, err := testQueries.
		GetUser(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user.Username, user2.Username)
	require.Equal(t, user.HashedPwd, user2.HashedPwd)
	require.Equal(t, user.FullName, user2.FullName)
	require.Equal(t, user.Email, user2.Email)
	require.WithinDuration(t,
		user.CreatedAt,
		user2.CreatedAt,
		time.Second)
}
