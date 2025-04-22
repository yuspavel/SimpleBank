package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var txKey = struct{}{}

type Store struct {
	*Queries
	dbt *pgxpool.Pool
}

type TransferTxParams struct {
	A1_ID  int64 `json:"from_account_id"`
	A2_ID  int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		dbt:     db,
		Queries: New(db),
	}
}
func (store *Store) execTX(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.dbt.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx) //получаем набор запросов, которые будут работать в рамках транзакции tx
	err = fn(q)  //отложенное выполнение запросов (коллбэк)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("TX err: %v\nTX rollback err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

func (store *Store) TransferTX(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTX(ctx, func(q *Queries) error { //начало выполнения коллбэка
		var err error

		txName := ctx.Value(txKey)

		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{FromAccountID: args.A1_ID, ToAccountID: args.A2_ID, Amount: args.Amount})
		if err != nil {
			return err
		}
		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{AccountID: args.A1_ID, Amount: -args.Amount})
		if err != nil {
			return err
		}
		fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{AccountID: args.A2_ID, Amount: args.Amount})
		if err != nil {
			return err
		}

		if args.A1_ID < args.A2_ID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, args.A1_ID, -args.Amount, args.A2_ID, args.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, args.A2_ID, args.Amount, args.A1_ID, -args.Amount)
		}
		return nil
	}) //конец выполнения коллбэка

	return result, err
}

func addMoney(ctx context.Context, q *Queries, account1ID, amount1, account2ID, amount2 int64) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{Amount: amount1, ID: account1ID})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{Amount: amount2, ID: account2ID})
	if err != nil {
		return
	}
	return
}
