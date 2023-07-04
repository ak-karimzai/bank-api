package token

import (
	"testing"
	"time"

	"github.com/ak-karimzai/bank-api/internel/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	secretKey := util.RandomString(32)
	maker, err := NewPasetoMaker(secretKey)
	require.NoError(t, err)

	username := util.RandomOwner()
	token, err := maker.CreateToken(username, time.Minute)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, payload.Username, username)
	require.True(t, time.Now().Before(payload.ExpiredAt))
}

func TestPasetoMakerWrongToken(t *testing.T) {
	secretKey := util.RandomString(32)
	maker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)

	username := util.RandomOwner()
	token, err := maker.CreateToken(username, time.Minute)
	require.NoError(t, err)

	randomToken := token[:len(token)-1] + util.RandomString(1)
	payload, err := maker.VerifyToken(randomToken)
	require.Error(t, err)
	require.Nil(t, payload)
}

func TestPasetoMakerExpiredToken(t *testing.T) {
	secretKey := util.RandomString(32)
	maker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)

	username := util.RandomOwner()
	token, err := maker.CreateToken(username, -time.Minute)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidToken)
	require.Nil(t, payload)
}
