package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	//run n concurent

	n := 5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransfersTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
		//check errors

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

			_, err = store.GetTransfers(context.Background(), transfer.ID)
			require.NoError(t, err)
			//check entry
			fromEntry := result.FromEntry
			require.NotEmpty(t, fromEntry)
			require.Equal(t, account1.ID, fromEntry.AccountID)
			require.Equal(t, -amount, fromEntry.Amount)
			require.NotZero(t, fromEntry.ID)
			require.NotZero(t, fromEntry.CreatedAt)

			_, err = store.GetEntries(context.Background(), fromEntry.ID)
			require.NoError(t, err)

			ToEntry := result.ToEntry
			require.NotEmpty(t, ToEntry)
			require.Equal(t, account2.ID, ToEntry.AccountID)
			require.Equal(t, amount, ToEntry.Amount)
			require.NotZero(t, ToEntry.ID)
			require.NotZero(t, ToEntry.CreatedAt)

			_, err = store.GetEntries(context.Background(), ToEntry.ID)
			require.NoError(t, err)

			//TODO: check Accounts balance

		}

	}

}
