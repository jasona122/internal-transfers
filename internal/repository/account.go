package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"internal-transfers/internal/domain"

	"internal-transfers/internal/model"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
)

// AccountRepository defines db operations for account data
//
//go:generate mockery --name=AccountRepository --output=./mocks --filename=account_mock.go --with-expecter
type AccountRepository interface {
	CreateAccount(ctx context.Context, account *model.Account) error
	GetAccount(ctx context.Context, accountID int64) (*model.Account, error)
	UpdateBalance(ctx context.Context, tx *sql.Tx, accountID int64, newBalance float64) error
}

// accountRepository is the Postgres implementation
type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) CreateAccount(ctx context.Context, account *model.Account) error {
	query := `INSERT INTO accounts (account_id, balance) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, account.AccountID, account.Balance)
	if err != nil {
		// case where account already exists
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return domain.ErrAccountDuplicate
		}
		return fmt.Errorf("create account failed: %w", err)
	}
	return nil
}

func (r *accountRepository) GetAccount(ctx context.Context, accountID int64) (*model.Account, error) {
	query := `SELECT account_id, balance FROM accounts WHERE account_id = $1`
	row := r.db.QueryRowContext(ctx, query, accountID)

	var acc model.Account
	if err := row.Scan(&acc.AccountID, &acc.Balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAccountNotFound
		}
		return nil, fmt.Errorf("get account failed: %w", err)
	}
	return &acc, nil
}

func (r *accountRepository) UpdateBalance(ctx context.Context, tx *sql.Tx, accountID int64, newBalance float64) error {
	query := `UPDATE accounts SET balance = $1 WHERE account_id = $2`
	_, err := tx.ExecContext(ctx, query, newBalance, accountID)
	if err != nil {
		return fmt.Errorf("update balance failed: %w", err)
	}
	return nil
}
