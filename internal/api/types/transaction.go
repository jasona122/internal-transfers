package types

type TransactionRequest struct {
	SourceAccountID      int64         `json:"source_account_id"`
	DestinationAccountID int64         `json:"destination_account_id"`
	Amount               FlexibleFloat `json:"amount"`
}

type TransactionResponse struct {
	TransactionID        int64         `json:"transaction_id"`
	SourceAccountID      int64         `json:"source_account_id"`
	DestinationAccountID int64         `json:"destination_account_id"`
	Amount               FlexibleFloat `json:"amount"`
	Timestamp            string        `json:"timestamp"`
}
