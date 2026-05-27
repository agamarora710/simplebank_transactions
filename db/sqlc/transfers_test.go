package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"agam/simplebank/util"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, fromAccount Account, toAccount Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	createRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	transfer1 := createRandomTransfer(t, account1, account2)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, account1, account2)
	}

	// ListTransfers uses OR — matches rows where account is sender OR receiver
	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		Limit:         5,
		Offset:        0,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account1.ID)
	}
}

func TestListTransfersByFromAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, account1, account2)
	}

	arg := ListTransfersByFromAccountParams{
		FromAccountID: account1.ID,
		Limit:         5,
		Offset:        0,
	}

	transfers, err := testQueries.ListTransfersByFromAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
	}
}

func TestListTransfersByToAccount(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, account1, account2)
	}

	arg := ListTransfersByToAccountParams{
		ToAccountID: account2.ID,
		Limit:       5,
		Offset:      0,
	}

	transfers, err := testQueries.ListTransfersByToAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, account2.ID, transfer.ToAccountID)
	}
}

func TestDeleteTransfer(t *testing.T) {
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)
	transfer1 := createRandomTransfer(t, account1, account2)

	err := testQueries.DeleteTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, sql.ErrNoRows)
	require.Empty(t, transfer2)
}
