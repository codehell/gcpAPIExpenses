package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/codehell/gcpApiExpenses/models"
	"github.com/codehell/gcpApiExpenses/responses"
	"net/http"
	"time"
)

func GetMovements(ctx context.Context, w http.ResponseWriter) {
	movements, err := models.GetMovements(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}
	jsonExpenses, _ := json.Marshal(movements)
	responses.JsonResponse(w, jsonExpenses, http.StatusOK)
}

func StoreMovement(ctx context.Context, w http.ResponseWriter, r *http.Request, email string) {
	movement := models.Movement{
		Email: email,
		CreateAt: time.Now().Unix(),
		UpdateAt: time.Now().Unix(),
	}
	if err := json.NewDecoder(r.Body).Decode(&movement); err != nil {
		responses.BadRequestApiError(w)
		return
	}
	err := models.StoreMovement(ctx, &movement)
	if err != nil {
		responses.BadRequestApiError(w)
	}
}
