package domain

import "errors"

var (
	ErrAccountDuplicate  = errors.New("account already exists")
	ErrAccountNotFound   = errors.New("account not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
)
