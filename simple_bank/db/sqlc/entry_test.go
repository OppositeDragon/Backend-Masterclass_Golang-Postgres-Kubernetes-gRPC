package db

import (
	"context"
	"database/sql"
	"simple_bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateRandomEntry(t *testing.T) Entry {
	account:= CreateRandomAccount(t)
	arg := CreateEntryParams{
		AccountId:   account.ID,
		Amount:  util.RandomMoney(), 
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountId, entry.AccountId)
	require.Equal(t, arg.Amount, entry.Amount) 
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	return entry;
}


func TestCreateEntry(t *testing.T) {
	entry:= CreateRandomEntry(t);
require.IsType(t, entry, Entry{})
}

func TestGetEntry(t *testing.T) {
	entry1 := CreateRandomEntry(t)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountId, entry2.AccountId)
	require.Equal(t, entry1.Amount, entry2.Amount) 
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestUpdateAEntry(t *testing.T) {
	account1 := CreateRandomEntry(t)
	arg := UpdateEntryParams{
		ID:      account1.ID,
		Amount: util.RandomMoney(),
	}
	account2, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.AccountId, account2.AccountId)
	require.Equal(t, arg.Amount, account2.Amount) 
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleleEntry(t *testing.T) {
	entry1 := CreateRandomEntry(t)
	err := testQueries.DeleteEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry2)
}

func TestGetEntries(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomEntry(t)
	}

	arg := GetEntriesParams{
		Limit:  5,
		Offset: 5,
	}

	entries, err := testQueries.GetEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
