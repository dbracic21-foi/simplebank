package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">>before", account1.Balance, account2.Balance)
	//run n concurent

	n := 5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := testStore.TransfersTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}
	//check errors

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check transfer

		transfer := result.Transfers
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTransfers(context.Background(), transfer.ID)
		require.NoError(t, err)
		//check entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testStore.GetEntries(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		ToEntry := result.ToEntry
		require.NotEmpty(t, ToEntry)
		require.Equal(t, account2.ID, ToEntry.AccountID)
		require.Equal(t, +amount, ToEntry.Amount)
		require.NotZero(t, ToEntry.ID)
		require.NotZero(t, ToEntry.CreatedAt)

		_, err = testStore.GetEntries(context.Background(), ToEntry.ID)
		require.NoError(t, err)

		//check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)
		fmt.Println(">>tx", fromAccount.Balance, toAccount.Balance)

		//check Accounts balance

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

	}

	//check final updated accounts

	UpdateAccount1, err := testStore.GetAccounts(context.Background(), account1.ID)
	require.NoError(t, err)

	UpdateAccount2, err := testStore.GetAccounts(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">>after", UpdateAccount1.Balance, UpdateAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, UpdateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, UpdateAccount2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">>before", account1.Balance, account2.Balance)
	//run n concurent

	n := 10
	amount := int64(10)
	errs := make(chan error)
	for i := 0; i < n; i++ {
		FromAccountID := account1.ID
		ToAccountID := account2.ID

		if i%2 == 1 {
			FromAccountID = account2.ID
			ToAccountID = account1.ID

		}
		go func() {
			ctx := context.Background()
			_, err := testStore.TransfersTx(ctx, TransferTxParams{
				FromAccountID: FromAccountID,
				ToAccountID:   ToAccountID,
				Amount:        amount,
			})
			errs <- err

		}()
	}
	//check errors

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}

	//check final updated accounts

	UpdateAccount1, err := testStore.GetAccounts(context.Background(), account1.ID)
	require.NoError(t, err)

	UpdateAccount2, err := testStore.GetAccounts(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">>after", UpdateAccount1.Balance, UpdateAccount2.Balance)

	require.Equal(t, account1.Balance, UpdateAccount1.Balance)
	require.Equal(t, account2.Balance, UpdateAccount2.Balance)

}
