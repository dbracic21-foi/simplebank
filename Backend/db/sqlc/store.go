package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	TransfersTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailsTx(ctx context.Context, arg VerifyEmailsTxParams) (VerifyEmailsTxResult, error)
}

type SQLStore struct {
	*Queries
	conpool *pgxpool.Pool
}

func NewStore(conpool *pgxpool.Pool) Store {
	return &SQLStore{
		conpool: conpool,
		Queries: New(conpool),
	}

}
