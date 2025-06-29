package model

import "time"

type Transaction struct {
	TransactionID        int64
	SourceAccountID      int64
	DestinationAccountID int64
	Amount               float64
	CreatedAt            time.Time
}
