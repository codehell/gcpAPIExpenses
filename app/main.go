package main

import (
	"fmt"
	"github.com/codehell/gcpApiExpenses/controllers"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/appengine"
	"log"
	"net/http"
	"time"
	"github.com/codehell/gcpApiExpenses/structs"
	"encoding/json"
)

func main() {
	http.HandleFunc("/", mainHandler)
	appengine.Main()
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	switch path := r.URL.Path; path {
	case "/movements":
		if r.Method == "GET" {
			controllers.GetMovements(ctx, w)
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

	case "/users":
		method := r.Method
		if method == "GET" {
			controllers.GetUsers(ctx, w)
		} else if method == "POST" {
			controllers.StoreUser(ctx, w, r)
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

	case "/test":
		method := r.Method
		if method == "GET" {

		} else if method == "POST" {
			controllers.StoreMovement(ctx, w, r)
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)

	case "/ok":
		w.WriteHeader(http.StatusOK)

	case "/login":
		login(w, r)

	case "/verify":
		http.Error(w, "error de prueba", http.StatusBadRequest)
		result, err := validateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MjUwMjQyMTIsImlhdCI6MTUyNTAyMzIxMiwiaXNzIjoid2ViIiwibmlja25hbWUiOiJjb2RlaGVsbCJ9.jMenCNxm5KZ0ozRzCHJNm0wFqC5SAtanQBM7jKxFEfk")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
		}
		w.WriteHeader(200)
		w.Write([]byte(result))
	}
}

func generateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":      "web",
		"nickname": "codehell",
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Unix() + 1000,
	})
	res, err := token.SignedString([]byte("sandruky"))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return res, nil
}

func validateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unspected signing method: %v", token.Header["alg"])
		}
		return []byte("sandruky"), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return fmt.Sprintf("%v", claims["nickname"]), nil
	}
	return "", err
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var login structs.Login
		if err := json.NewDecoder(r.Body).Decode(login); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		token, err := generateToken()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write([]byte(token))
	}
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}
