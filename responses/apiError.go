package responses

import (
	"net/http"
	"encoding/json"
)

type ApiError struct {
	Code   string `json:"code"`
	Reason string `json:"reason"`
}

//Api Codes
const (
	MethodNotAllowed = "methodNotAllowed"
	BadRequest       = "badRequest"
	LoginFailed      = "loginFailed"
)

func (e ApiError) NewApiError(code int, w http.ResponseWriter) {
	jsonErr, _ := json.Marshal(e)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(jsonErr))
}

func MethodNotAllowedApiError(w http.ResponseWriter) {
	apiError := ApiError{
		MethodNotAllowed,
		"verb is not allowed",
	}
	apiError.NewApiError(
		http.StatusMethodNotAllowed, w)
}

func BadRequestApiError(w http.ResponseWriter) {
	apiError := ApiError{
		Code:   BadRequest,
		Reason: "I can't do nothing with this shit",
	}
	apiError.NewApiError(http.StatusBadRequest, w)
}
