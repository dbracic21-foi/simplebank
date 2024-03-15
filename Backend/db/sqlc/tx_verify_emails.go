package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type VerifyEmailsTxParams struct {
	EmailId    int64
	SecretCode string
}

type VerifyEmailsTxResult struct {
	Users        User
	VerifyEmails VerifyEmail
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
		result.Users, err = q.UpdateUser(ctx, UpdateUserParams{
			Username: result.VerifyEmails.Username,
			IsEmailVerified: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		})

		return err

	})
	return result, err
}
