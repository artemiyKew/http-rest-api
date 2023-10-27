package apiserver

import (
	"time"

	"github.com/artemiyKew/http-rest-api/internal/app/model"
	"github.com/golang-jwt/jwt"
)

func generateJWT(user *model.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["sub"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
