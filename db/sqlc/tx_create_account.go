package db

import "context"

type CreateAccountTxParams struct {
	CreateAccountParams
	AfterCreate func(account Account) error
}

type CreateAccountTxResult struct {
	Account Account
}

func (store *SQLStore) CreateAccountTx(ctx context.Context, arg CreateAccountTxParams) (CreateAccountTxResult, error) {
	var result CreateAccountTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Account, err = q.CreateAccount(ctx, arg.CreateAccountParams)
		if err != nil {
			return err
		}
		return arg.AfterCreate(result.Account)

	})
	return result, err
}
