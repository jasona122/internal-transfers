package handler

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"

	"internal-transfers/internal/api/types"
	"internal-transfers/internal/domain"
	"internal-transfers/internal/service"
)

type AccountHandler struct {
	accountService service.AccountService
}

func NewAccountHandler(svc service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: svc}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req types.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("error decoding body")
		types.WriteResponseError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.InitialBalance < 0 {
		log.Warn().Msg("attempt to create account with negative balance")
		types.WriteResponseError(w, http.StatusBadRequest, "initial balance cannot be negative")
		return
	}

	err := h.accountService.CreateAccount(r.Context(), req.AccountID, float64(req.InitialBalance))
	if err != nil {
		if errors.Is(err, domain.ErrAccountDuplicate) {
			log.Warn().Err(err).Msg("attempt to create account that already exists")
			types.WriteResponseError(w, http.StatusConflict, "account has already been created")
			return
		}
		log.Error().Err(err).Msg("error creating account")
		types.WriteResponseError(w, http.StatusInternalServerError, "failed to create account")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	accountIDStr := r.URL.Path[len("/accounts/"):]
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse account id")
		types.WriteResponseError(w, http.StatusBadRequest, "invalid account id")
		return
	}
	acc, err := h.accountService.GetAccount(r.Context(), accountID)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			log.Warn().Msg("account not found")
			types.WriteResponseError(w, http.StatusNotFound, "account not found")
			return
		}
		log.Error().Err(err).Msg("failed to get account")
		types.WriteResponseError(w, http.StatusInternalServerError, "failed to get account")
		return
	}

	resp := types.AccountResponse{
		AccountID: acc.AccountID,
		Balance:   acc.Balance,
	}

	types.WriteResponseSuccess(w, resp)
}
