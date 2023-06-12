package db

import (
	"context"
	"testing"
	"time"

	"github.com/ak-karimzai/ak-karimzai/simpleb/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPwd, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPwd,
		Fullname:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Fullname, user.Fullname)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordLastChanged.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user2.Username, user1.Username)
	require.Equal(t, user2.HashedPassword, user1.HashedPassword)
	require.Equal(t, user2.Fullname, user1.Fullname)
	require.Equal(t, user2.Email, user1.Email)
	require.WithinDuration(t, user2.PasswordLastChanged, user1.PasswordLastChanged, time.Second)
	require.WithinDuration(t, user2.CreatedAt, user1.CreatedAt, time.Second)
}
