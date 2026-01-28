package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	// Test with a simple password
	password := RandomString(10)
	hashedPassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	err = CheckPasswordHash(password, hashedPassword1)
	require.NoError(t, err)

	wrongPassword := RandomString(10)
	err = CheckPasswordHash(wrongPassword, hashedPassword1)
	require.Equal(t, err, bcrypt.ErrMismatchedHashAndPassword)

	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)

	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
