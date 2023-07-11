package db

import (
	"context"
	"testing"
	"time"

	"github.com/ak-karimzai/bank-api/internel/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
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

func TestUpdateUserFullName(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := util.RandomOwner()

	user, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		FullName: pgtype.Text{
			String: newFullName,
			Valid:  true,
		},
		Username: oldUser.Username,
	})
	require.NoError(t, err)

	require.Equal(t, newFullName, user.FullName)
	require.Equal(t, oldUser.Username, user.Username)
	require.Equal(t, oldUser.Email, user.Email)
}

func TestUpdateUserEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmailAddress := util.RandomEmail()

	user, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Email: pgtype.Text{
			String: newEmailAddress,
			Valid:  true,
		},
		Username: oldUser.Username,
	})
	require.NoError(t, err)

	require.Equal(t, newEmailAddress, user.Email)
	require.Equal(t, oldUser.FullName, user.FullName)
	require.Equal(t, oldUser.Username, user.Username)
}

func TestUpdateUserPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	randomPassword := util.RandomString(32)
	hashedPwd, err := util.HashPasswrod(randomPassword)
	require.NoError(t, err)

	user, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		HashedPwd: pgtype.Text{
			String: hashedPwd,
			Valid:  true,
		},
		Username: oldUser.Username,
	})
	require.NoError(t, err)

	require.Equal(t, hashedPwd, user.HashedPwd)
	require.Equal(t, oldUser.FullName, user.FullName)
	require.Equal(t, oldUser.Username, user.Username)
	require.NoError(t,
		bcrypt.CompareHashAndPassword([]byte(user.HashedPwd), []byte(randomPassword)))
}
