/*Unit тесты создания денежного трансфера в рамках транзакции БД.*/
package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTX(t *testing.T) {
	errs := make(chan error)
	results := make(chan TransferTxResult)

	amount := int64(10)

	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5

	fmt.Printf("Balance before transactions: account 1=%d, accont 2=%d\n", account1.Balance, account2.Balance)
	for i := 0; i < n; i++ {

		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)

			result, err := store.TransferTX(ctx, TransferTxParams{A1_ID: account1.ID,
				A2_ID:  account2.ID,
				Amount: amount,
			})

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		r := <-results
		require.NotEmpty(t, r)

		//Check transfer
		transfer := r.Transfer
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		_, err = store.GetTransfer(context.Background(), r.Transfer.ID)
		require.NoError(t, err)

		//Check FromEntry
		fromEntry := r.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		//Check ToEntry
		toEntry := r.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//Check balance of account From
		updatedaccount1 := r.FromAccount
		require.Equal(t, updatedaccount1.ID, account1.ID)
		require.NotEmpty(t, r.FromAccount)
		//require.Equal(t, account1.Balance-amount, faccount.Balance) //Тоже рабочий вариант проверки правильности баланса акаунта после обновления

		updatedaccount2 := r.ToAccount
		require.NotEmpty(t, r.ToAccount)
		require.Equal(t, updatedaccount2.ID, account2.ID)
		//require.Equal(t, account2.Balance+amount, taccount.Balance) //Тоже рабочий вариант проверки правильности баланса акаунта после обновления

		fmt.Printf("Current balance: account 1=%d, accont 2=%d\n", updatedaccount1.Balance, updatedaccount2.Balance)

		diff1 := account1.Balance - r.FromAccount.Balance
		diff2 := r.ToAccount.Balance - account2.Balance
		fmt.Println("Amount:", amount, "Diff1:", diff1)
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	final_f_account, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	final_t_account, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Printf("Balance after transactions: account 1=%d, accont 2=%d\n", final_f_account.Balance, final_t_account.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, final_f_account.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, final_t_account.Balance)
}

func TestTransferTXDeadlock(t *testing.T) {
	errs := make(chan error)

	amount := int64(10)

	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 10
	fmt.Printf("Balance before transactions: account 1=%d, accont 2=%d\n", account1.Balance, account2.Balance)
	for i := 0; i < n; i++ {
		fromAccountId := account1.ID
		toAccountId := account2.ID

		txName := fmt.Sprintf("tx %d", i+1)

		if i%2 == 1 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
		}
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)

			_, err := store.TransferTX(ctx, TransferTxParams{A1_ID: fromAccountId,
				A2_ID:  toAccountId,
				Amount: amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
