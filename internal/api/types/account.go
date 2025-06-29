package types

type CreateAccountRequest struct {
	AccountID      int64         `json:"account_id"`
	InitialBalance FlexibleFloat `json:"initial_balance"`
}

type AccountResponse struct {
	AccountID int64   `json:"account_id"`
	Balance   float64 `json:"balance"`
}
