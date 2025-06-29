package handler

import (
	"encoding/json"
	"errors"
	"internal-transfers/internal/domain"
	"net/http"

	"internal-transfers/internal/api/types"
	"internal-transfers/internal/service"
)

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(svc service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: svc}
}

func (h *TransactionHandler) SubmitTransaction(w http.ResponseWriter, r *http.Request) {
	var req types.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		types.WriteResponseError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Amount <= 0 {
		types.WriteResponseError(w, http.StatusBadRequest, "amount must be positive")
		return
	}
	err := h.transactionService.ProcessTransaction(r.Context(), req.SourceAccountID, req.DestinationAccountID, float64(req.Amount))
	if errors.Is(err, domain.ErrInsufficientFunds) {
		types.WriteResponseError(w, http.StatusBadRequest, "insufficient funds from source account")
		return
	}

	if err != nil {
		types.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
