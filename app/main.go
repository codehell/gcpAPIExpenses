package main

import (
	"fmt"
	"github.com/codehell/gcpApiExpenses/controllers"
	"google.golang.org/appengine"
	"net/http"
	"github.com/codehell/gcpApiExpenses/responses"
	"github.com/codehell/gcpApiExpenses/auth"
)

func main() {
	http.HandleFunc("/", mainHandler)
	appengine.Main()
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	method := r.Method
	//w.Header().Set("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	if contentType := r.Header.Get("Accept"); contentType != "application/json" && method != "OPTIONS" {
		responses.BadRequestApiError(w)
		return
	}
	
	switch path := r.URL.Path; path {

	case "/movements":
		if method == "GET" {
			controllers.GetMovements(ctx, w)
		} else if method == "POST" {
			controllers.StoreMovement(ctx, w, r)
		} else if method != "OPTIONS" {
			responses.MethodNotAllowedApiError(w)
		}

	case "/ok":
		w.WriteHeader(http.StatusOK)

	case "/login":
		if method == "POST" {
			auth.Login(ctx, w, r)
		} else if method != "OPTIONS" {
			responses.MethodNotAllowedApiError(w)
		}

	case "/users":
		if method == "GET" {
			controllers.GetUsers(ctx, w)
		} else if method == "POST" {
			controllers.StoreUser(ctx, w, r)
		} else if method != "OPTIONS" {
			responses.MethodNotAllowedApiError(w)
		}

	case "/verify":
		http.Error(w, "error de prueba", http.StatusBadRequest)
		result, err := auth.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MjUwMjQyMTIsImlhdCI6MTUyNTAyMzIxMiwiaXNzIjoid2ViIiwibmlja25hbWUiOiJjb2RlaGVsbCJ9.jMenCNxm5KZ0ozRzCHJNm0wFqC5SAtanQBM7jKxFEfk")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
		}
		w.WriteHeader(200)
		w.Write([]byte(result))
	}
}
