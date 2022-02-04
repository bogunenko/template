package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	errorNotEnoughBalance = errors.New("not enough balance")
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

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func (store *Store) Transfer(ctx context.Context, arg CreateTransactionParams) (sql.Result, error) {
	var result sql.Result

	err := store.execTx(ctx, func(q *Queries) error {

		from, err := q.GetAccountForUpdate(ctx, arg.FromID)
		if err != nil {
			return err
		}
		if from.Balance < arg.Amount {
			return errorNotEnoughBalance
		}

		_, err = q.Withdraw(ctx, WithdrawParams{ID: arg.FromID, Amount: arg.Amount})
		if err != nil {
			return err
		}
		_, err = q.Deposit(ctx, DepositParams{ID: arg.ToID, Amount: arg.Amount})
		if err != nil {
			return err
		}

		result, err = q.CreateTransaction(ctx, arg)

		return err
	})

	return result, err
}
