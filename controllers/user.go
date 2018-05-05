package controllers

import (
	"context"
	"net/http"
	"github.com/codehell/gcpApiExpenses/models"
	"fmt"
	"encoding/json"
	"github.com/codehell/gcpApiExpenses/responses"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(ctx context.Context, w http.ResponseWriter) {
	users, err := models.GetUsers(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}
	jsonUsers, _ := json.Marshal(users)
	responses.JsonResponse(w, jsonUsers, http.StatusOK)
}

func StoreUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		responses.BadRequestApiError(w)
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		responses.BadRequestApiError(w)
		return
	}
	user.Password = string(hashed)
	if err := user.StoreUser(ctx); err != nil {
		responses.BadRequestApiError(w)
	}
	w.WriteHeader(http.StatusCreated)
}