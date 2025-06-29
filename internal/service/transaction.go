package service

import (
	"context"
	"database/sql"
	"fmt"
	"internal-transfers/internal/domain"
	"time"

	"internal-transfers/internal/model"
	"internal-transfers/internal/repository"
)

//go:generate mockery --name=TransactionService --filename=transaction_mock.go --output=./mocks --with-expecter
type TransactionService interface {
	ProcessTransaction(ctx context.Context, sourceID, destID int64, amount float64) error
}

type transactionService struct {
	txRepo  repository.TransactionRepository
	accRepo repository.AccountRepository
	db      *sql.DB // for transaction control
}

func NewTransactionService(txRepo repository.TransactionRepository, accRepo repository.AccountRepository, db *sql.DB) TransactionService {
	return &transactionService{
		txRepo:  txRepo,
		accRepo: accRepo,
		db:      db,
	}
}

// ProcessTransaction processes a funds transfer between accounts ensuring atomicity
func (s *transactionService) ProcessTransaction(ctx context.Context, sourceID, destID int64, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	sourceAcc, err := s.accRepo.GetAccount(ctx, sourceID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	if sourceAcc == nil {
		_ = tx.Rollback()
		return domain.ErrAccountNotFound
	}
	if sourceAcc.Balance < amount {
		_ = tx.Rollback()
		return domain.ErrInsufficientFunds
	}

	destAcc, err := s.accRepo.GetAccount(ctx, destID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	if destAcc == nil {
		_ = tx.Rollback()
		return domain.ErrAccountNotFound
	}

	if err := s.accRepo.UpdateBalance(ctx, tx, sourceID, sourceAcc.Balance-amount); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to update source balance: %w", err)
	}
	if err := s.accRepo.UpdateBalance(ctx, tx, destID, destAcc.Balance+amount); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to update destination balance: %w", err)
	}

	transaction := &model.Transaction{
		SourceAccountID:      sourceID,
		DestinationAccountID: destID,
		Amount:               amount,
		CreatedAt:            time.Now(),
	}

	if err := s.txRepo.CreateTransaction(ctx, tx, transaction); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to insert transaction record: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}
