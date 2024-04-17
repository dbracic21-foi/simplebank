package db

import "context"

type ListAccountsTxParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListAccountsTxResult struct {
	Accounts []Account `json:"accounts"`
}

func (store *SQLStore) ListAccountsTx(ctx context.Context, arg ListAccountsTxParams) (ListAccountsTxResult, error) {
	var result ListAccountsTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Accounts, err = q.ListAccounts(ctx, ListAccountsParams{
			Limit:  arg.Limit,
			Offset: arg.Offset,
		})
		if err != nil {
			return err
		}

		return nil

	})
	return result, err
}
