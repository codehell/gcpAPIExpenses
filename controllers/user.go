package controllers

import (
	"context"
	"net/http"
	"github.com/codehell/gcpApiExpenses/models"
	"fmt"
	"encoding/json"
	"github.com/codehell/gcpApiExpenses/responses"
)

func GetUsers(ctx context.Context, w http.ResponseWriter) {
	users, err := models.GetUsers(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}
	jsonUsers, _ := json.Marshal(users)
	responses.OkJsonRespond(w, jsonUsers)
}

func StoreUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if err := user.StoreUser(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusCreated)
}