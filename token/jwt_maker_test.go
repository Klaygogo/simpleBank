package token

import (
	"testing"
	"time"

	"github.com/Klaygogo/simplebank/util"
	"github.com/dgrijalva/jwt-go"
	required "github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	required.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	token, err := maker.CreateToken(username, duration)
	required.NoError(t, err)
	required.NotEmpty(t, token)

	claims, err := maker.VerifyToken(token)
	required.NoError(t, err)
	required.NotEmpty(t, claims)

	required.Equal(t, username, claims.Username)
	required.WithinDuration(t, issuedAt, claims.IssuedAt, time.Second)
	required.WithinDuration(t, expiredAt, claims.ExpiredAt, time.Second)
}

func TestExpiredJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	required.NoError(t, err)

	username := util.RandomOwner()

	token, err := maker.CreateToken(username, -time.Second)
	required.NoError(t, err)
	required.NotEmpty(t, token)

	_, err = maker.VerifyToken(token)
	required.Error(t, err)
	required.Equal(t, ErrExpiredToken, err)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	required.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute
	payload, err := NewPayload(username, duration)
	required.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	required.NoError(t, err)
	required.NotEmpty(t, token)

	_, err = maker.VerifyToken(token)
	required.Error(t, err)
	required.Equal(t, ErrInvalidToken, err)

}
