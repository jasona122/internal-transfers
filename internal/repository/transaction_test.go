package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"internal-transfers/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionRepository_CreateTransaction(t *testing.T) {
	// given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &transactionRepository{db: db}
	ctx := context.Background()

	transaction := &model.Transaction{
		SourceAccountID:      1,
		DestinationAccountID: 2,
		Amount:               50.0,
		CreatedAt:            time.Now(),
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		txObj, err := db.Begin()
		require.NoError(t, err)

		mock.ExpectQuery(`INSERT INTO transactions`).
			WithArgs(transaction.SourceAccountID, transaction.DestinationAccountID, transaction.Amount, transaction.CreatedAt).
			WillReturnRows(sqlmock.NewRows([]string{"transaction_id"}).AddRow(123))

		err = repo.CreateTransaction(ctx, txObj, transaction)
		assert.NoError(t, err)
		assert.Equal(t, int64(123), transaction.TransactionID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectBegin()
		txObj, err := db.Begin()
		require.NoError(t, err)

		mock.ExpectQuery(`INSERT INTO transactions`).
			WithArgs(transaction.SourceAccountID, transaction.DestinationAccountID, transaction.Amount, transaction.CreatedAt).
			WillReturnError(assert.AnError)

		// when
		err = repo.CreateTransaction(ctx, txObj, transaction)

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create transaction failed")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestTransactionRepository_GetTransaction(t *testing.T) {
	// given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &transactionRepository{db: db}
	ctx := context.Background()
	transactionID := int64(10)
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"transaction_id",
			"source_account_id",
			"destination_account_id",
			"amount",
			"created_at",
		}).AddRow(transactionID, 1, 2, 100.0, now)

		mock.ExpectQuery(`SELECT transaction_id, source_account_id, destination_account_id, amount, created_at FROM transactions WHERE transaction_id = \$1`).
			WithArgs(transactionID).
			WillReturnRows(rows)

		tx, err := repo.GetTransaction(ctx, transactionID)
		assert.NoError(t, err)
		require.NotNil(t, tx)
		assert.Equal(t, transactionID, tx.TransactionID)
		assert.Equal(t, int64(1), tx.SourceAccountID)
		assert.Equal(t, int64(2), tx.DestinationAccountID)
		assert.Equal(t, 100.0, tx.Amount)
		assert.WithinDuration(t, now, tx.CreatedAt, time.Second)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT transaction_id, source_account_id, destination_account_id, amount, created_at FROM transactions WHERE transaction_id = \$1`).
			WithArgs(transactionID).
			WillReturnError(sql.ErrNoRows)

		tx, err := repo.GetTransaction(ctx, transactionID)
		assert.NoError(t, err)
		assert.Nil(t, tx)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT transaction_id, source_account_id, destination_account_id, amount, created_at FROM transactions WHERE transaction_id = \$1`).
			WithArgs(transactionID).
			WillReturnError(assert.AnError)

		tx, err := repo.GetTransaction(ctx, transactionID)
		assert.Error(t, err)
		assert.Nil(t, tx)
		assert.Contains(t, err.Error(), "get transaction failed")

		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestTransactionRepository_ListTransactions(t *testing.T) {
	// given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &transactionRepository{db: db}
	ctx := context.Background()

	t.Run("success with multiple rows", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"transaction_id",
			"source_account_id",
			"destination_account_id",
			"amount",
			"created_at",
		}).
			AddRow(1, 1, 2, 100.0, now).
			AddRow(2, 3, 4, 200.0, now.Add(time.Minute))

		mock.ExpectQuery(`SELECT transaction_id, source_account_id, destination_account_id, amount, created_at FROM transactions`).
			WillReturnRows(rows)

		txs, err := repo.ListTransactions(ctx)
		assert.NoError(t, err)
		assert.Len(t, txs, 2)
		assert.Equal(t, int64(1), txs[0].TransactionID)
		assert.Equal(t, int64(2), txs[1].TransactionID)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db query error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT transaction_id, source_account_id, destination_account_id, amount, created_at FROM transactions`).
			WillReturnError(assert.AnError)

		txs, err := repo.ListTransactions(ctx)
		assert.Error(t, err)
		assert.Nil(t, txs)
		assert.Contains(t, err.Error(), "list transactions failed")

		require.NoError(t, mock.ExpectationsWereMet())
	})
}
