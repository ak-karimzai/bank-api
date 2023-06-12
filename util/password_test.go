package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hashedPwd, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPwd)

	err = CheckPassword(password, hashedPwd)
	require.NoError(t, err)

	wrongPwd := RandomString(6)
	err = CheckPassword(wrongPwd, hashedPwd)
	require.EqualError(t,
		err,
		bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPwd2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPwd)
	require.NotEqual(t, hashedPwd, hashedPwd2)
}
