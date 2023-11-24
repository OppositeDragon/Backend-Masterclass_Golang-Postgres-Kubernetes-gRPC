// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package db

import (
	"database/sql"
	"time"
)

type Account struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Currency  string    `json:"currency"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

type Entry struct {
	ID        int64 `json:"id"`
	AccountId int64 `json:"accountId"`
	// can be positive or negative
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"createdAt"`
}

type Transfer struct {
	ID            int64 `json:"id"`
	FromAccountId int64 `json:"fromAccountId"`
	ToAccountId   int64 `json:"toAccountId"`
	// must be positive
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"createdAt"`
}

type User struct {
	Username          string         `json:"username"`
	Name1             string         `json:"name1"`
	Name2             sql.NullString `json:"name2"`
	Lastname1         string         `json:"lastname1"`
	Lastname2         sql.NullString `json:"lastname2"`
	Email             string         `json:"email"`
	HashedPassword    string         `json:"hashedPassword"`
	PasswordChangedAt time.Time      `json:"passwordChangedAt"`
	CreatedAt         time.Time      `json:"createdAt"`
}
