package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/codehell/gcpApiExpenses/controllers"
	"github.com/codehell/gcpApiExpenses/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/appengine"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/codehell/gcpApiExpenses/responses"
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
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	switch path := r.URL.Path; path {
	case "/movements":
		if method == "GET" {
			controllers.GetMovements(ctx, w)
		} else if method != "OPTIONS" {
			responses.MethodNotAllowedApiError(w)
		}

	case "/users":
		if method == "GET" {
			controllers.GetUsers(ctx, w)
		} else if method == "POST" {
			controllers.StoreUser(ctx, w, r)
		}

	case "/test":
		if method == "GET" {

		} else if method == "POST" {
			controllers.StoreMovement(ctx, w, r)
		}

	case "/ok":
		w.WriteHeader(http.StatusOK)

	case "/login":
		if method == "POST" {
			login(ctx, w, r)
		} else if method != "OPTIONS" {
			responses.MethodNotAllowedApiError(w)
		}

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
		"exp":      time.Now().Unix() + (60 * 120),
	})
	formedToken, err := token.SignedString([]byte("Tyig<Mead1"))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return formedToken, nil
}

func validateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unspected signing method: %v", token.Header["alg"])
		}
		return []byte("Tyig<Mead1"), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return fmt.Sprintf("%v", claims["nickname"]), nil
	}
	return "", err
}

func login(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	auth := strings.Split(r.Header.Get("Authorization"), " ")
	if len(auth) != 2 {
		http.Error(w, "authorization failed", http.StatusUnauthorized)
	}
	decoded, _ := base64.StdEncoding.DecodeString(auth[1])
	payload := strings.Split(string(decoded), ":")
	username := payload[0]
	password := payload[1]
	user := models.User{
		Username: username,
	}

	getUSerErr := user.GetUser(ctx)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil || getUSerErr != nil {
		apiErr := responses.ApiError{
			Code:   responses.LoginFailed,
			Reason: "Incorrect credentials",
		}
		apiErr.NewApiError(http.StatusUnauthorized, w)
	}

	token, err := generateToken()
	if err != nil {
		responses.BadRequestApiError(w)
	}

	responses.LoginResponse(w, token)
}
