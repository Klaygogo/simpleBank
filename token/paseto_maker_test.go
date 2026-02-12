package token

import (
	"testing"
	"time"

	"github.com/Klaygogo/simplebank/util"
	required "github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	required.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	token, payload, err := maker.CreateToken(username, duration)
	required.NoError(t, err)
	required.NotEmpty(t, payload)
	required.NotEmpty(t, token)

	claims, err := maker.VerifyToken(token)
	required.NoError(t, err)
	required.NotEmpty(t, claims)

	required.Equal(t, username, claims.Username)
	required.WithinDuration(t, issuedAt, claims.IssuedAt.Time, time.Second)
	required.WithinDuration(t, expiredAt, claims.ExpiredAt.Time, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	required.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	required.NoError(t, err)
	required.NotEmpty(t, payload)
	required.NotEmpty(t, token)

	_, err = maker.VerifyToken(token)
	required.Error(t, err)
	required.Equal(t, ErrExpiredToken, err)
}
