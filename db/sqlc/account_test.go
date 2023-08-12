package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dbracic21-foi/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Accounts {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RadnomMoney(),
		Currency: util.RandomCurrency(),
	}
	Accounts, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, Accounts)

	require.Equal(t, arg.Owner, Accounts.Owner)
	require.Equal(t, arg.Balance, Accounts.Balance)
	require.Equal(t, arg.Currency, Accounts.Currency)

	require.NotZero(t, Accounts.ID)
	require.NotZero(t, Accounts.CreatedAt)

	return Accounts

}

func TestCreateAccounts(t *testing.T) {
	createRandomAccount(t)

}

func TestGetAccounts(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccounts(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}

func TestUpdateAccounts(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountsParams{
		ID:      account1.ID,
		Balance: util.RadnomMoney(),
	}
	account2, err := testQueries.UpdateAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	err := testQueries.DeleteAccounts(context.Background(), account1.ID)

	require.NoError(t, err)

	account2, err := testQueries.GetAccounts(context.Background(), account1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)

}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	Accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, Accounts, 5)

	for _, account := range Accounts {
		require.NotEmpty(t, account)
	}
}
