package db

import (
	"context"
	"database/sql"
)

type VerifyEmailsTxParams struct {
	EmailId    int64
	SecretCode string
}

type VerifyEmailsTxResult struct {
	Users        Users
	VerifyEmails VerifyEmails
}

func (store *SQLStore) VerifyEmailsTx(ctx context.Context, arg VerifyEmailsTxParams) (VerifyEmailsTxResult, error) {
	var result VerifyEmailsTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.VerifyEmails, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}
		result.Users,err = q.UpdateUser(ctx, UpdateUserParams{
			Username: result.VerifyEmails.Username,
			IsEmailVerified: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		})

		return err

	})
	return result, err
}
