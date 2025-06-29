package types

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func WriteResponseError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	resp := ErrorResponse{
		Code:    code,
		Message: msg,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func WriteResponseSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := SuccessResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	}
	_ = json.NewEncoder(w).Encode(resp)
}
