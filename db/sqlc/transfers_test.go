package db

import (
	"context"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {

	args_a := ListAccountsParams{Limit: 5, Offset: 0}
	accounts, err := testQueries.ListAccounts(context.Background(), args_a)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	af := rand.IntN(len(accounts))
	at := rand.IntN(len(accounts))

	args_t := CreateTransferParams{FromAccountID: accounts[af].ID, ToAccountID: accounts[at].ID, Amount: (accounts[af].Balance / 100) * 8}
	transfer, err := testQueries.CreateTransfer(context.Background(), args_t)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, transfer.FromAccountID, accounts[af].ID)
	require.Equal(t, transfer.ToAccountID, accounts[at].ID)
	require.Equal(t, transfer.Amount, (accounts[af].Balance/100)*8)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestListTransfer(t *testing.T) {
	transfers1 := []Transfer{}
	for range 100 {
		tranfer := createRandomTransfer(t)
		transfers1 = append(transfers1, tranfer)
	}

	tf := rand.IntN(len(transfers1))
	tt := rand.IntN(len(transfers1))

	args := ListTransfersParams{FromAccountID: transfers1[tf].FromAccountID, ToAccountID: transfers1[tt].ID, Limit: 2, Offset: 0}
	transfers, err := testQueries.ListTransfers(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, transfers, 2)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}

func TestGetTransfer(t *testing.T) {
	transfer1 := createRandomTransfer(t)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt.Time, transfer2.CreatedAt.Time, time.Second)
}
