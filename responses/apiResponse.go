package responses

import (
	"net/http"
	"encoding/json"
)

type Login struct {
	AccessToken string `json:"access_token"`
	ExpiresIn int64 `json:"expires_in"`
}

func JsonResponse(w http.ResponseWriter, json []byte, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(json))
}

func LoginResponse(w http.ResponseWriter, token string) {

	login := Login{
		AccessToken: token,
		ExpiresIn: 60 * 120,
	}
	json, _:= json.Marshal(login)
	JsonResponse(w, json, http.StatusOK)
}
