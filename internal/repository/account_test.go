package repository

import (
	"context"
	"database/sql"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"internal-transfers/internal/domain"
	"testing"

	"internal-transfers/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountRepository_CreateAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &accountRepository{db: db}
	ctx := context.Background()
	account := &model.Account{
		AccountID: 123,
		Balance:   100.0,
	}

	t.Run("create account successfully", func(t *testing.T) {
		// given
		mock.ExpectExec(`INSERT INTO accounts`).
			WithArgs(account.AccountID, account.Balance).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// when
		err = repo.CreateAccount(ctx, account)

		// then
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("create account fail due to duplicate", func(t *testing.T) {
		// given
		mock.ExpectExec(`INSERT INTO accounts`).
			WithArgs(account.AccountID, account.Balance).
			WillReturnError(&pq.Error{Code: pgerrcode.UniqueViolation})

		// when
		err = repo.CreateAccount(ctx, account)

		// then
		assert.ErrorIs(t, err, domain.ErrAccountDuplicate)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("create account fail due to database error", func(t *testing.T) {
		// given
		mock.ExpectExec(`INSERT INTO accounts`).
			WithArgs(account.AccountID, account.Balance).
			WillReturnError(assert.AnError) // any unexpected error

		// when
		err = repo.CreateAccount(ctx, account)

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create account failed")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAccountRepository_GetAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &accountRepository{db: db}
	ctx := context.Background()
	accountID := int64(123)

	t.Run("get account successfully", func(t *testing.T) {
		// given
		rows := sqlmock.NewRows([]string{"account_id", "balance"}).
			AddRow(accountID, 100.0)

		mock.ExpectQuery(`SELECT account_id, balance FROM accounts WHERE account_id = \$1`).
			WithArgs(accountID).
			WillReturnRows(rows)

		// when
		account, err := repo.GetAccount(ctx, accountID)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, accountID, account.AccountID)
		assert.Equal(t, 100.0, account.Balance)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("get account fail due to account not found", func(t *testing.T) {
		// given
		mock.ExpectQuery(`SELECT account_id, balance FROM accounts WHERE account_id = \$1`).
			WithArgs(accountID).
			WillReturnError(sql.ErrNoRows)

		// when
		account, err := repo.GetAccount(ctx, accountID)

		// then
		assert.Nil(t, account)
		assert.ErrorIs(t, err, domain.ErrAccountNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("get account fail due to database error", func(t *testing.T) {
		// given
		mock.ExpectQuery(`SELECT account_id, balance FROM accounts WHERE account_id = \$1`).
			WithArgs(accountID).
			WillReturnError(assert.AnError)

		// when
		account, err := repo.GetAccount(ctx, accountID)

		// then
		assert.Nil(t, account)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "get account failed")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAccountRepository_UpdateBalance(t *testing.T) {
	// given
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &accountRepository{db: db}
	ctx := context.Background()
	accountID := int64(123)
	newBalance := 200.0

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		tx, err := db.Begin()
		require.NoError(t, err)

		mock.ExpectExec(`UPDATE accounts SET balance =`).
			WithArgs(newBalance, accountID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// when
		err = repo.UpdateBalance(ctx, tx, accountID, newBalance)

		// then
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		mock.ExpectBegin()
		tx, err := db.Begin()
		require.NoError(t, err)

		mock.ExpectExec(`UPDATE accounts SET balance =`).
			WithArgs(newBalance, accountID).
			WillReturnError(assert.AnError)

		// when
		err = repo.UpdateBalance(ctx, tx, accountID, newBalance)

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update balance failed")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
