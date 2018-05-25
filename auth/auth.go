package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/codehell/gcpApiExpenses/models"
	"github.com/codehell/gcpApiExpenses/responses"
	"github.com/codehell/gcpApiExpenses/structs"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
)

func GenerateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":   "web",
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Unix() + (60 * 120),
	})
	formedToken, err := token.SignedString([]byte("Tyig<Mead1"))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return formedToken, nil
}

func ValidateToken(tokenString string, claims *structs.Claim) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unspected signing method: %v", token.Header["alg"])
		}
		return []byte("Tyig<Mead1"), nil
	})
	if err != nil {
		return false, err
	}
	if res, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		claims.Email = res["email"].(string)
		return true, nil
	}
	return false, err
}

func Login(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	auth := strings.Split(r.Header.Get("Authorization"), " ")
	if len(auth) != 2 {
		responses.BadRequestApiError(w)
		return
	}
	decoded, _ := base64.StdEncoding.DecodeString(auth[1])

	payload := strings.Split(string(decoded), ":")
	email := payload[0]
	password := payload[1]
	user := models.User{
		Email: email,
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

	token, err := GenerateToken(email)
	if err != nil {
		responses.BadRequestApiError(w)
		return
	}

	responses.LoginResponse(w, token)
}
