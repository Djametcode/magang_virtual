package helper

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(id uint, email string) (string, error) {
	secretKey := "your-secret-key"

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = id
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte(secretKey))

	return tokenString, err
}