package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	pwd := RandomString(6)

	hashedPwd, err := HashPasswrod(pwd)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPwd)

	err = CheckPwd(pwd, hashedPwd)
	require.NoError(t, err)

	wrongPwd := RandomString(6)
	err = CheckPwd(wrongPwd, hashedPwd)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPwd2, err := HashPasswrod(pwd)
	require.NoError(t, err)
	require.NotEqual(t, hashedPwd, hashedPwd2)
}

func TestHashFailed(t *testing.T) {
	pwd := RandomString(73)

	hashedPwd, err := HashPasswrod(pwd)
	require.Error(t, err)
	require.Empty(t, hashedPwd)
}
