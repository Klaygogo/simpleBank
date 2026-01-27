package db

import (
	"context"
	"testing"

	"github.com/Klaygogo/simplebank/util"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	passwd := util.RandomString(6)
	hashPasswd, _ := util.HashPassword(passwd)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashPasswd,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	err = util.CheckPasswordHash(passwd, user.HashedPassword)
	require.NoError(t, err)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.PasswordChangedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
}
