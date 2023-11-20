package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// executes a function within a transaction.
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"fromAccountId"`
	ToAccountID   int64 `json:"toAccountId"`
	Amount        int64 `json:"amount"`
}
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"fromAccount"`
	ToAccount   Account  `json:"toAccount"`
	FromEntry   Entry    `json:"fromEntry"`
	ToEntry     Entry    `json:"toEntry"`
}

var txKey = struct{}{}

// creates a transfer record, account entries and updates accounts' balance within a transaction.
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		fmt.Print(ctx.Value(txKey))
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountId: arg.FromAccountID,
			ToAccountId:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountId: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountId: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		//TODO: update accounts' balance
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, _ = AddAmount(ctx, q, arg.FromAccountID, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, _ = AddAmount(ctx, q, arg.ToAccountID, arg.FromAccountID, -arg.Amount)
		}
		return nil
	})
	return result, err
}

func AddAmount(ctx context.Context,
	q *Queries,
	accountId1 int64,
	accountId2 int64,
	amount int64,
) (
	account1 Account,
	account2 Account,
	err error,
) {
	account1, err = q.AddAmountAccount(ctx, AddAmountAccountParams{
		ID:     accountId1,
		Amount: -amount,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAmountAccount(ctx, AddAmountAccountParams{
		ID:     accountId2,
		Amount: amount,
	})
	return

}
