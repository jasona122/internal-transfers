package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"internal-transfers/internal/domain"
	"internal-transfers/internal/model"
	"internal-transfers/internal/repository/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestSetup(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *mocks.TransactionRepository, *mocks.AccountRepository, TransactionService) {
	db, mockSql, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}

	txRepo := mocks.NewTransactionRepository(t)
	accRepo := mocks.NewAccountRepository(t)
	service := NewTransactionService(txRepo, accRepo, db)

	return db, mockSql, txRepo, accRepo, service
}

func TestTransactionService_ProcessTransaction(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		db, mockSql, txRepo, accRepo, service := newTestSetup(t)
		defer db.Close()

		source := &model.Account{AccountID: 1, Balance: 200}
		dest := &model.Account{AccountID: 2, Balance: 50}
		amount := 50.0

		mockSql.ExpectBegin()
		mockSql.ExpectCommit()

		accRepo.EXPECT().GetAccount(ctx, source.AccountID).Return(source, nil)
		accRepo.EXPECT().GetAccount(ctx, dest.AccountID).Return(dest, nil)

		accRepo.EXPECT().
			UpdateBalance(mock.Anything, mock.AnythingOfType("*sql.Tx"), source.AccountID, source.Balance-amount).
			Return(nil)
		accRepo.EXPECT().
			UpdateBalance(mock.Anything, mock.AnythingOfType("*sql.Tx"), dest.AccountID, dest.Balance+amount).
			Return(nil)

		txRepo.EXPECT().
			CreateTransaction(mock.Anything, mock.AnythingOfType("*sql.Tx"), mock.MatchedBy(func(tx *model.Transaction) bool {
				return tx.SourceAccountID == source.AccountID &&
					tx.DestinationAccountID == dest.AccountID &&
					tx.Amount == amount
			})).
			Return(nil)

		err := service.ProcessTransaction(ctx, source.AccountID, dest.AccountID, amount)
		assert.NoError(t, err)

		assert.NoError(t, mockSql.ExpectationsWereMet())
	})

	t.Run("insufficient funds from source account", func(t *testing.T) {
		db, mockSql, _, accRepo, service := newTestSetup(t)
		defer db.Close()

		mockSql.ExpectBegin()
		mockSql.ExpectRollback()

		source := &model.Account{AccountID: 1, Balance: 10}
		accRepo.EXPECT().GetAccount(ctx, source.AccountID).Return(source, nil)

		err := service.ProcessTransaction(ctx, source.AccountID, 2, 50)
		assert.ErrorIs(t, err, domain.ErrInsufficientFunds)

		assert.NoError(t, mockSql.ExpectationsWereMet())
	})

	t.Run("source account not found", func(t *testing.T) {
		db, mockSql, _, accRepo, service := newTestSetup(t)
		defer db.Close()

		mockSql.ExpectBegin()
		mockSql.ExpectRollback()

		accRepo.EXPECT().GetAccount(ctx, int64(1)).Return(nil, nil)

		err := service.ProcessTransaction(ctx, 1, 2, 50)
		assert.ErrorIs(t, err, domain.ErrAccountNotFound)

		assert.NoError(t, mockSql.ExpectationsWereMet())
	})

	t.Run("destination account not found", func(t *testing.T) {
		db, mockSql, _, accRepo, service := newTestSetup(t)
		defer db.Close()

		mockSql.ExpectBegin()
		mockSql.ExpectRollback()

		source := &model.Account{AccountID: 1, Balance: 100}
		accRepo.EXPECT().GetAccount(ctx, source.AccountID).Return(source, nil)
		accRepo.EXPECT().GetAccount(ctx, int64(2)).Return(nil, nil)

		err := service.ProcessTransaction(ctx, source.AccountID, 2, 50)
		assert.ErrorIs(t, err, domain.ErrAccountNotFound)

		assert.NoError(t, mockSql.ExpectationsWereMet())
	})

	t.Run("transaction create error when inserting", func(t *testing.T) {
		db, mockSql, txRepo, accRepo, service := newTestSetup(t)
		defer db.Close()

		source := &model.Account{AccountID: 1, Balance: 200}
		dest := &model.Account{AccountID: 2, Balance: 50}
		amount := 50.0

		mockSql.ExpectBegin()
		mockSql.ExpectRollback()

		accRepo.EXPECT().GetAccount(ctx, source.AccountID).Return(source, nil)
		accRepo.EXPECT().GetAccount(ctx, dest.AccountID).Return(dest, nil)

		accRepo.EXPECT().
			UpdateBalance(ctx, mock.AnythingOfType("*sql.Tx"), source.AccountID, source.Balance-amount).
			Return(nil)
		accRepo.EXPECT().
			UpdateBalance(ctx, mock.AnythingOfType("*sql.Tx"), dest.AccountID, dest.Balance+amount).
			Return(nil)

		txRepo.EXPECT().
			CreateTransaction(mock.Anything, mock.AnythingOfType("*sql.Tx"), mock.Anything).
			Return(errors.New("insert error"))

		err := service.ProcessTransaction(ctx, source.AccountID, dest.AccountID, amount)
		assert.ErrorContains(t, err, "insert error")

		assert.NoError(t, mockSql.ExpectationsWereMet())
	})
}
