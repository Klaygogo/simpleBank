package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool" //  导入pgxpool包，这是一个PostgreSQL连接池库
)

type Store interface {
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	Querier
}

type SQLStore struct {
	*Queries
	db *pgxpool.Pool
}

var txKey = struct{}{}

func NewStore(db *pgxpool.Pool) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

type TransferTxParams struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
}

type TransferTxResult struct {
	Transfer    Transfer
	FromAccount Account
	ToAccount   Account
	FromEntry   Entry
	ToEntry     Entry
}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	// 相当于是实现了具体的callback回调函数
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		update1 := AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		}
		update2 := AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, err = q.AddAccountBalance(ctx, update1)
			if err != nil {
				return err
			}
			result.ToAccount, err = q.AddAccountBalance(ctx, update2)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, err = q.AddAccountBalance(ctx, update2)
			if err != nil {
				return err
			}
			result.FromAccount, err = q.AddAccountBalance(ctx, update1)
			if err != nil {
			}
		}

		return nil
	})
	return result, err
}
