package service

import (
	"context"
	"internal-transfers/internal/domain"

	"internal-transfers/internal/model"
	"internal-transfers/internal/repository"
)

//go:generate mockery --name=AccountService --filename=account_mock.go --output=./mocks --with-expecter
type AccountService interface {
	CreateAccount(ctx context.Context, accountID int64, balance float64) error
	GetAccount(ctx context.Context, accountID int64) (*model.Account, error)
}

type accountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{repo: repo}
}

// CreateAccount creates a new account with initial balance; assumes negative balance is not allowed
func (s *accountService) CreateAccount(ctx context.Context, accountID int64, initialBalance float64) error {
	if initialBalance <= 0 {
		return domain.ErrInsufficientFunds
	}

	acc := &model.Account{
		AccountID: accountID,
		Balance:   initialBalance,
	}
	return s.repo.CreateAccount(ctx, acc)
}

// GetAccount retrieves account details by ID
func (s *accountService) GetAccount(ctx context.Context, accountID int64) (*model.Account, error) {
	acc, err := s.repo.GetAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	return acc, nil
}
