package repository

import (
	"context"
	"database/sql"
	"fmt"

	"internal-transfers/internal/model"
)

// TransactionRepository defines db operations for transactions
//
//go:generate mockery --name=TransactionRepository --filename=transaction_mock.go --output=./mocks --with-expecter
type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *sql.Tx, transaction *model.Transaction) error
	GetTransaction(ctx context.Context, transactionID int64) (*model.Transaction, error)
	ListTransactions(ctx context.Context) ([]*model.Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreateTransaction(ctx context.Context, tx *sql.Tx, transaction *model.Transaction) error {
	query := `
        INSERT INTO transactions (source_account_id, destination_account_id, amount, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING transaction_id`
	err := tx.QueryRowContext(ctx, query,
		transaction.SourceAccountID, transaction.DestinationAccountID, transaction.Amount, transaction.CreatedAt).
		Scan(&transaction.TransactionID)
	if err != nil {
		return fmt.Errorf("create transaction failed: %w", err)
	}
	return nil
}

func (r *transactionRepository) GetTransaction(ctx context.Context, transactionID int64) (*model.Transaction, error) {
	query := `
        SELECT transaction_id, source_account_id, destination_account_id, amount, created_at
        FROM transactions
        WHERE transaction_id = $1`
	row := r.db.QueryRowContext(ctx, query, transactionID)

	var tx model.Transaction
	if err := row.Scan(
		&tx.TransactionID,
		&tx.SourceAccountID,
		&tx.DestinationAccountID,
		&tx.Amount,
		&tx.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get transaction failed: %w", err)
	}
	return &tx, nil
}

func (r *transactionRepository) ListTransactions(ctx context.Context) ([]*model.Transaction, error) {
	query := `SELECT transaction_id, source_account_id, destination_account_id, amount, created_at
			  FROM transactions`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list transactions failed: %w", err)
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		var tx model.Transaction
		if err := rows.Scan(
			&tx.TransactionID,
			&tx.SourceAccountID,
			&tx.DestinationAccountID,
			&tx.Amount,
			&tx.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		transactions = append(transactions, &tx)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return transactions, nil
}
