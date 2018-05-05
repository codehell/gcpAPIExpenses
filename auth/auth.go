package auth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"log"
	"fmt"
	"net/http"
	"strings"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"github.com/codehell/gcpApiExpenses/responses"
	"github.com/codehell/gcpApiExpenses/models"
	"context"
)

func GenerateToken() (string, error) {
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

func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unspected signing method: %v", token.Header["alg"])
		}
		return []byte("Tyig<Mead1"), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return fmt.Sprintf("%v", claims["email"]), nil
	}
	return "", err
}

func Login(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	auth := strings.Split(r.Header.Get("Authorization"), " ")
	if len(auth) != 2 {
		responses.BadRequestApiError(w)
		return
	}
	decoded, _ := base64.StdEncoding.DecodeString(auth[1])

	payload := strings.Split(string(decoded), ":")
	username := payload[0]
	password := payload[1]
	user := models.User{
		Email: username,
	}

	getUSerErr := user.GetUser(ctx)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil || getUSerErr != nil {
		apiErr := responses.ApiError{
			Code:   responses.LoginFailed,
			Reason: "Incorrect credentials",
		}
		apiErr.NewApiError(http.StatusUnauthorized, w)
		return
	}

	token, err := GenerateToken()
	if err != nil {
		responses.BadRequestApiError(w)
		return
	}

	responses.LoginResponse(w, token)
}
