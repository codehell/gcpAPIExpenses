package responses

import (
	"net/http"
)

func OkJsonRespond(w http.ResponseWriter, json []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(json))
}