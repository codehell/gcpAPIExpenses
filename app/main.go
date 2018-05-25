package main

import (
	"github.com/codehell/gcpApiExpenses/auth"
	"github.com/codehell/gcpApiExpenses/controllers"
	"github.com/codehell/gcpApiExpenses/responses"
	"github.com/codehell/gcpApiExpenses/structs"
	"google.golang.org/appengine"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", mainHandler)
	appengine.Main()
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	method := r.Method
	a := false
	var claims structs.Claim
	//w.Header().Set("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	splited := strings.Split(r.URL.Path, "/")

	//Check correct content type
	/*if contentType := r.Header.Get("Content-Type"); contentType != "application/json" && method != "OPTIONS" {
		responses.BadRequestApiError(w)
		return
	}*/

	if method != "OPTIONS" {
		a = tokenValidation(r, &claims)
	} else {
		a = true
	}

	switch path := splited[1]; path {

	case "movements":
		if !a {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if method == "GET" {
			controllers.GetMovements(ctx, w)
		} else if method == "POST" {
			controllers.StoreMovement(ctx, w, r, claims.Email)
		} else if method != "OPTIONS" {
			responses.MethodNotAllowedApiError(w)
		}

	case "ok":
		w.WriteHeader(http.StatusOK)

	case "login":
		if method == "POST" {
			auth.Login(ctx, w, r)
		} else if method != "OPTIONS" {
			responses.MethodNotAllowedApiError(w)
		}

	case "users":
		if method == "GET" {
			controllers.GetUsers(ctx, w)
		} else if method == "POST" {
			controllers.StoreUser(ctx, w, r)
		} else if method != "OPTIONS" {
			responses.MethodNotAllowedApiError(w)
		}

	default:
		http.NotFound(w, r)
	}
}

func tokenValidation(r *http.Request, claims *structs.Claim) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		splAuthHeader := strings.Split(authHeader, " ")
		if len(splAuthHeader) == 2 {
			if splAuthHeader[0] == "Bearer" {
				isValidate, err := auth.ValidateToken(splAuthHeader[1], claims)
				if err != nil {
					return false
				} else {
					if isValidate {
						return true
					}
				}
			}
		}
	}
	return false
}
