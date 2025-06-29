package handler

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/mock"
	"internal-transfers/internal/domain"
	"internal-transfers/internal/service/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionHandler_SubmitTransaction(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		mockSvc := &mocks.TransactionService{}
		h := NewTransactionHandler(mockSvc)
		reqBody := `{"source_account_id": 1, "destination_account_id": 2, "amount": 100}`
		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader([]byte(reqBody)))
		w := httptest.NewRecorder()

		// when
		mockSvc.
			On("ProcessTransaction", mock.Anything, int64(1), int64(2), 100.0).
			Return(nil)

		// then
		h.SubmitTransaction(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("invalid json", func(t *testing.T) {
		// given
		mockSvc := &mocks.TransactionService{}
		h := NewTransactionHandler(mockSvc)
		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader([]byte(`invalid`)))
		w := httptest.NewRecorder()

		// when/then
		h.SubmitTransaction(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("negative amount", func(t *testing.T) {
		// given
		mockSvc := &mocks.TransactionService{}
		h := NewTransactionHandler(mockSvc)
		reqBody := `{"source_account_id": 1, "destination_account_id": 2, "amount": -50}`
		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader([]byte(reqBody)))
		w := httptest.NewRecorder()

		// when/then
		h.SubmitTransaction(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("insufficient funds", func(t *testing.T) {
		// given
		mockSvc := &mocks.TransactionService{}
		h := NewTransactionHandler(mockSvc)
		reqBody := `{"source_account_id": 1, "destination_account_id": 2, "amount": 100}`
		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader([]byte(reqBody)))
		w := httptest.NewRecorder()

		// when
		mockSvc.
			On("ProcessTransaction", mock.Anything, int64(1), int64(2), 100.0).
			Return(domain.ErrInsufficientFunds)

		// then
		h.SubmitTransaction(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("generic service error", func(t *testing.T) {
		// given
		mockSvc := &mocks.TransactionService{}
		h := NewTransactionHandler(mockSvc)
		reqBody := `{"source_account_id": 1, "destination_account_id": 2, "amount": 100}`
		req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader([]byte(reqBody)))
		w := httptest.NewRecorder()

		// when
		mockSvc.
			On("ProcessTransaction", mock.Anything, int64(1), int64(2), 100.0).
			Return(errors.New("some db error"))

		// then
		h.SubmitTransaction(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

}
