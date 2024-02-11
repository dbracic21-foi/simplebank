package db

import "context"

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user Users) error
}

type CreateUserTxResult struct {
	Users Users
}

func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Users, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}
		return arg.AfterCreate(result.Users)

	})
	return result, err
}
