// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Accounts struct {
	ID          int64         `json:"id"`
	Owner       string        `json:"owner"`
	Balance     int64         `json:"balance"`
	Currency    string        `json:"currency"`
	CreatedAt   time.Time     `json:"created_at"`
	CountryCode sql.NullInt32 `json:"country_code"`
}

type Entries struct {
	ID        int64     `json:"id"`
	AccountID int64     `json:"account_id"`
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type Sessions struct {
	ID              uuid.UUID `json:"id"`
	Username        string    `json:"username"`
	RefreshToken    string    `json:"refresh_token"`
	UserAgent       string    `json:"user_agent"`
	ClientIp        string    `json:"client_ip"`
	IsEmailVerified bool      `json:"is_email_verified"`
	IsBlocked       bool      `json:"is_blocked"`
	ExpiresAt       time.Time `json:"expires_at"`
	CreatedAt       time.Time `json:"created_at"`
}

type Transfers struct {
	ID            int64     `json:"id"`
	FromAccountID int64     `json:"from_account_id"`
	ToAccountID   int64     `json:"to_account_id"`
	Amount        int64     `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

type Users struct {
	Username          string    `json:"username"`
	HashedPassword    string    `json:"hashed_password"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	IsEmailVerified   bool      `json:"is_email_verified"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}
