package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

var tokenSecret []byte = []byte("dont read me") // should be something very random and secret instead!!!

type Token struct {
}

/// NewToken to create a new jwt token and return signed string with our secret.
func (t *Token) NewToken(username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject: username, // should associate random ID to the user and used as subject
	})

	tokenString, _ := token.SignedString(tokenSecret)
	return tokenString
}

/// VerifyToken return username if the tokenString is a valid jwt token signed by us
func (t *Token) VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return tokenSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"].(string), nil
	} else {
		return "", err
	}
}
