package main

import (
	"github.com/codehell/gcpApiExpenses/models"
	"google.golang.org/appengine"
	"net/http"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"google.golang.org/appengine/user"
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
		expense := models.Movement{
			Username: "codehell",
			Amount: -3223,
			Description: "Compra en carrefour incluyendo una sarten",
			Tags: []string{"comida", "cocina"},
			CreateAt: time.Now().Unix(),
			UpdateAt: time.Now().Unix(),
		}
		err := models.StoreExpense(ctx, &expense)
		log.Println(err)
	case "/auth":
		w.Header().Set("Content-type", "text/html; charset=utf-8")
		ctx := appengine.NewContext(r)
		u := user.Current(ctx)
		if u == nil {
			url, _ := user.LoginURL(ctx, "/auth")
			fmt.Fprintf(w, `<a href="%s">Sign in or register</a>`, url)
			return
		}
		url, _ := user.LogoutURL(ctx, "/auth")
		fmt.Fprintf(w, `Welcome, %s! (<a href="%s">sign out</a>)`, u, url)
	}
}
