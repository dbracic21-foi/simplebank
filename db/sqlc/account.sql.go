// Code generated by sqlc. DO NOT EDIT.
// source: account.sql

package db

import (
	"context"
)

const addAccountsBalance = `-- name: AddAccountsBalance :one
UPDATE accounts
  set balance = balance + $1
WHERE id = $2
RETURNING id, owner, balance, currency, created_at, country_code
`

type AddAccountsBalanceParams struct {
	Amount int64 `json:"amount"`
	ID     int64 `json:"id"`
}

func (q *Queries) AddAccountsBalance(ctx context.Context, arg AddAccountsBalanceParams) (Accounts, error) {
	row := q.db.QueryRowContext(ctx, addAccountsBalance, arg.Amount, arg.ID)
	var i Accounts
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.CountryCode,
	)
	return i, err
}

const createAccount = `-- name: CreateAccount :one
INSERT INTO accounts (
  owner,
  balance,
  currency
) VALUES (
  $1, $2, $3
) RETURNING id, owner, balance, currency, created_at, country_code
`

type CreateAccountParams struct {
	Owner    string `json:"owner"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Accounts, error) {
	row := q.db.QueryRowContext(ctx, createAccount, arg.Owner, arg.Balance, arg.Currency)
	var i Accounts
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.CountryCode,
	)
	return i, err
}

const deleteAccounts = `-- name: DeleteAccounts :exec
DELETE FROM accounts
WHERE id = $1
`

func (q *Queries) DeleteAccounts(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAccounts, id)
	return err
}

const getAccounts = `-- name: GetAccounts :one
SELECT id, owner, balance, currency, created_at, country_code FROM accounts
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAccounts(ctx context.Context, id int64) (Accounts, error) {
	row := q.db.QueryRowContext(ctx, getAccounts, id)
	var i Accounts
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.CountryCode,
	)
	return i, err
}

const getAccountsForUpdate = `-- name: GetAccountsForUpdate :one
SELECT id, owner, balance, currency, created_at, country_code FROM accounts
WHERE id = $1 LIMIT 1
FOR  NO KEY UPDATE
`

func (q *Queries) GetAccountsForUpdate(ctx context.Context, id int64) (Accounts, error) {
	row := q.db.QueryRowContext(ctx, getAccountsForUpdate, id)
	var i Accounts
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.CountryCode,
	)
	return i, err
}

const listAccounts = `-- name: ListAccounts :many
SELECT id, owner, balance, currency, created_at, country_code FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListAccountsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Accounts, error) {
	rows, err := q.db.QueryContext(ctx, listAccounts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Accounts
	for rows.Next() {
		var i Accounts
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
			&i.CountryCode,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAccounts = `-- name: UpdateAccounts :one
UPDATE accounts
  set balance = $2
WHERE id = $1
RETURNING id, owner, balance, currency, created_at, country_code
`

type UpdateAccountsParams struct {
	ID      int64 `json:"id"`
	Balance int64 `json:"balance"`
}

func (q *Queries) UpdateAccounts(ctx context.Context, arg UpdateAccountsParams) (Accounts, error) {
	row := q.db.QueryRowContext(ctx, updateAccounts, arg.ID, arg.Balance)
	var i Accounts
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.CountryCode,
	)
	return i, err
}
