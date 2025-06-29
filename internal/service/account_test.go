package service

import (
	"context"
	"errors"
	"testing"

	"internal-transfers/internal/domain"
	"internal-transfers/internal/model"
	"internal-transfers/internal/repository/mocks"

	"github.com/stretchr/testify/assert"
)

func TestAccountService_CreateAccount(t *testing.T) {
	ctx := context.Background()
	repo := mocks.NewAccountRepository(t)
	service := NewAccountService(repo)

	t.Run("success", func(t *testing.T) {
		repo.EXPECT().
			CreateAccount(ctx, &model.Account{AccountID: 1, Balance: 100}).
			Return(nil)

		err := service.CreateAccount(ctx, 1, 100)
		assert.NoError(t, err)
	})

	t.Run("invalid balance", func(t *testing.T) {
		err := service.CreateAccount(ctx, 1, 0)
		assert.ErrorIs(t, err, domain.ErrInsufficientFunds)
	})

	t.Run("repo error", func(t *testing.T) {
		repo.EXPECT().
			CreateAccount(ctx, &model.Account{AccountID: 2, Balance: 100}).
			Return(errors.New("db error"))

		err := service.CreateAccount(ctx, 2, 100)
		assert.ErrorContains(t, err, "db error")
	})
}
