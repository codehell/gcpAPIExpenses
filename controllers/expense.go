package controllers

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/codehell/gcpApiExpenses/responses"
	"github.com/codehell/gcpApiExpenses/models"
	"context"
	"time"
	"log"
)

func GetMovements(ctx context.Context, w http.ResponseWriter)  {
	movements, err := models.GetMovements(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
	}
	jsonExpenses, _ := json.Marshal(movements)
	responses.OkJsonRespond(w, jsonExpenses)
}

func StoreMovement(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	expense := models.Movement{
		Username:    "codehell",
		Amount:      -3223,
		Description: "Compra en carrefour incluyendo una sarten",
		Tags:        []string{"comida", "cocina"},
		CreateAt:    time.Now().Unix(),
		UpdateAt:    time.Now().Unix(),
	}
	err := models.StoreMovement(ctx, &expense)
	log.Println(err)

}
