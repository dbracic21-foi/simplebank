package db

import (
	"context"
	"testing"
	"time"

	"github.com/dbracic21-foi/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) Users {

	hashedPassword, err := u
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	Users, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, Users)

	require.Equal(t, arg.Username, Users.Username)
	require.Equal(t, arg.HashedPassword, Users.HashedPassword)
	require.Equal(t, arg.FullName, Users.FullName)
	require.Equal(t, arg.Email, Users.Email)

	require.True(t, Users.PasswordChangedAt.IsZero())
	require.NotZero(t, Users.CreatedAt)

	return Users

}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)

}

func TestGetUser(t *testing.T) {
	User1 := createRandomUser(t)
	User2, err := testQueries.GetUser(context.Background(), User1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, User2)

	require.Equal(t, User1.Username, User2.Username)
	require.Equal(t, User1.HashedPassword, User2.HashedPassword)
	require.Equal(t, User1.FullName, User2.FullName)
	require.Equal(t, User1.Email, User2.Email)

	require.WithinDuration(t, User1.CreatedAt, User2.CreatedAt, time.Second)
	require.WithinDuration(t, User1.PasswordChangedAt, User2.PasswordChangedAt, time.Second)

}
