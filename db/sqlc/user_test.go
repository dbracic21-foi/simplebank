package db

import (
	"context"
	"testing"
	"time"

	"github.com/dbracic21-foi/simplebank/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {

	hashedPassword, err := util.HashPassword(util.RadnomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	Users, err := testStore.CreateUser(context.Background(), arg)

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
	User2, err := testStore.GetUser(context.Background(), User1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, User2)

	require.Equal(t, User1.Username, User2.Username)
	require.Equal(t, User1.HashedPassword, User2.HashedPassword)
	require.Equal(t, User1.FullName, User2.FullName)
	require.Equal(t, User1.Email, User2.Email)

	require.WithinDuration(t, User1.CreatedAt, User2.CreatedAt, time.Second)
	require.WithinDuration(t, User1.PasswordChangedAt, User2.PasswordChangedAt, time.Second)

}
func TestUpdateOnlyFullNameUser(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := util.RandomOwner()
	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: pgtype.Text{
			String: newFullName,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)

}
func TestUpdateOnlyEmailUser(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := util.RandomEmail()
	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: pgtype.Text{
			String: newEmail,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)

}
func TestUpdateOnlyPasswordUser(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword := util.RadnomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)
	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: pgtype.Text{
			String: newHashedPassword,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)

}
func TestUpdateAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()
	newPassword := util.RadnomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)
	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: pgtype.Text{
			String: newHashedPassword,
			Valid:  true,
		},
		FullName: pgtype.Text{
			String: newFullName,
			Valid:  true,
		},
		Email: pgtype.Text{
			String: newEmail,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.Equal(t, newEmail, updatedUser.Email)

}
