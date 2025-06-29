package api

import (
	"net/http"

	"internal-transfers/internal/api/handler"
	"internal-transfers/internal/api/middleware"
	"internal-transfers/internal/service"
)

func NewRouter(
	accountSvc service.AccountService,
	transactionSvc service.TransactionService,
) http.Handler {

	mux := http.NewServeMux()

	accountHandler := handler.NewAccountHandler(accountSvc)
	transactionHandler := handler.NewTransactionHandler(transactionSvc)

	// Account endpoints
	mux.HandleFunc("/accounts/", withMethod(http.MethodGet, accountHandler.GetAccount)) // expects /accounts/{id}
	mux.HandleFunc("/accounts", withMethod(http.MethodPost, accountHandler.CreateAccount))

	// Transaction endpoints
	mux.HandleFunc("/transactions", withMethod(http.MethodPost, transactionHandler.SubmitTransaction))

	return middleware.RecoverPanic(mux)
}

// helper to enforce allowed methods
func withMethod(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}
