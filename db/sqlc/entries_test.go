package db

import (
	"context"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {

	args := ListAccountsParams{Limit: 5, Offset: 0}
	accounts, err := testQueries.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, accounts, 5)
	k := rand.IntN(len(accounts))

	args_ent := CreateEntryParams{AccountID: accounts[k].ID, Amount: (accounts[k].Balance / 100) * 13}
	entry, err := testQueries.CreateEntry(context.Background(), args_ent)
	require.NoError(t, err)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	require.Less(t, entry.Amount, accounts[k].Balance)
	require.Equal(t, entry.AccountID, accounts[k].ID)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry1 := createRandomEntry(t)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
	require.WithinDuration(t, entry1.CreatedAt.Time, entry1.CreatedAt.Time, time.Second)
}

func TestListEntries(t *testing.T) {
	entries1 := []Entry{}
	for i := 0; i < 5; i++ {
		entry := createRandomEntry(t)
		entries1 = append(entries1, entry)
	}

	k := len(entries1)
	args := ListEntriesParams{AccountID: entries1[rand.IntN(k)].AccountID, Limit: 1, Offset: 0}

	entries, err := testQueries.ListEntries(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, entries, 1)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}

}
