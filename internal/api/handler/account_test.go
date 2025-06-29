package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"internal-transfers/internal/api/types"
	"internal-transfers/internal/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"internal-transfers/internal/domain"
	"internal-transfers/internal/service/mocks" // import path to your generated mocks

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAccountHandler_CreateAccount(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		mockSvc := mocks.NewAccountService(t)
		h := NewAccountHandler(mockSvc)

		reqBody := `{"account_id": 1, "initial_balance": 100}`
		req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		mockSvc.EXPECT().
			CreateAccount(mock.Anything, int64(1), 100.0).
			Return(nil).
			Once()

		// when
		h.CreateAccount(w, req)

		// then
		resp := w.Result()
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("invalid body", func(t *testing.T) {
		// given
		mockSvc := mocks.NewAccountService(t)
		h := NewAccountHandler(mockSvc)

		req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(`{invalid json`))
		w := httptest.NewRecorder()

		// when
		h.CreateAccount(w, req)

		// then
		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("negative initial balance", func(t *testing.T) {
		// given
		mockSvc := mocks.NewAccountService(t)
		h := NewAccountHandler(mockSvc)

		reqBody := `{"account_id": 1, "initial_balance": -10}`
		req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		// when
		h.CreateAccount(w, req)

		// then
		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("duplicate account", func(t *testing.T) {
		// given
		mockSvc := mocks.NewAccountService(t)
		h := NewAccountHandler(mockSvc)

		reqBody := `{"account_id": 123, "initial_balance": 10}`
		req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBufferString(reqBody))
		w := httptest.NewRecorder()

		mockSvc.EXPECT().
			CreateAccount(mock.Anything, int64(123), 10.0).
			Return(domain.ErrAccountDuplicate).
			Once()

		// when
		h.CreateAccount(w, req)

		// then
		resp := w.Result()
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		mockSvc.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		// given
		mockSvc := mocks.NewAccountService(t)
		h := NewAccountHandler(mockSvc)

		reqBody := `{"account_id": 1, "initial_balance": 100}`
		req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		mockSvc.EXPECT().
			CreateAccount(mock.Anything, int64(1), 100.0).
			Return(errors.New("some db error")).
			Once()

		// when
		h.CreateAccount(w, req)

		// then
		resp := w.Result()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestAccountHandler_GetAccount(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		mockSvc := mocks.NewAccountService(t)
		h := NewAccountHandler(mockSvc)

		accountID := int64(123)
		account := &model.Account{
			AccountID: accountID,
			Balance:   100.23344,
		}

		mockSvc.EXPECT().
			GetAccount(mock.Anything, accountID).
			Return(account, nil).
			Once()

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%d", accountID), nil)
		w := httptest.NewRecorder()

		// when
		h.GetAccount(w, req)

		// then
		resp := w.Result()
		fmt.Printf("RESPONSE BODY IS: %v\n", resp)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var gotResp struct {
			Code    int                   `json:"code"`
			Message string                `json:"message"`
			Data    types.AccountResponse `json:"data"`
		}

		err := json.NewDecoder(resp.Body).Decode(&gotResp)
		assert.NoError(t, err)
		assert.Equal(t, account.AccountID, gotResp.Data.AccountID)
		assert.Equal(t, account.Balance, gotResp.Data.Balance)

		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid account id", func(t *testing.T) {
		// given
		mockSvc := mocks.NewAccountService(t)
		h := NewAccountHandler(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/accounts/abc", nil)
		w := httptest.NewRecorder()

		// when
		h.GetAccount(w, req)

		// then
		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		mockSvc.AssertNotCalled(t, "GetAccount", mock.Anything, mock.Anything)
	})

	t.Run("account not found", func(t *testing.T) {
		// given
		mockSvc := mocks.NewAccountService(t)
		h := NewAccountHandler(mockSvc)

		accountID := int64(123)
		mockSvc.EXPECT().
			GetAccount(mock.Anything, accountID).
			Return(nil, domain.ErrAccountNotFound).
			Once()

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%d", accountID), nil)
		w := httptest.NewRecorder()

		// when
		h.GetAccount(w, req)

		// then
		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		// given
		mockSvc := mocks.NewAccountService(t)
		h := NewAccountHandler(mockSvc)

		accountID := int64(123)
		mockSvc.EXPECT().
			GetAccount(mock.Anything, accountID).
			Return(nil, errors.New("service failure")).
			Once()

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%d", accountID), nil)
		w := httptest.NewRecorder()

		// when
		h.GetAccount(w, req)

		// then
		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})
}
