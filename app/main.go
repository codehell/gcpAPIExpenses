package main

import (
	"github.com/codehell/gcpApiExpenses/models"
	"google.golang.org/appengine"
	"net/http"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func main() {
	http.HandleFunc("/", expensesHandler)
	appengine.Main()
}

func expensesHandler(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	switch path := r.URL.Path; path {
	case "/expenses":
		expenses, err := models.GetExpenses(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
		}
		jsonExpenses, _ := json.Marshal(expenses)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(jsonExpenses))
	case "/users":
		users, err := models.GetUsers(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
		}
		jsonUsers, _:= json.Marshal(users)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(jsonUsers))
	case "/test":
		expense := models.Expense{
			Username: "codehell",
			Amount: 3223,
			Description: "Compra en carrefour incluyendo una sarten",
			Tags: []string{"comida", "cocina"},
			CreateAt: time.Now().Unix(),
			UpdateAt: time.Now().Unix(),
		}
		err := models.StoreExpense(ctx, &expense)
		log.Println(err)
	}
}
