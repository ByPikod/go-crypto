package helpers

import (
	"github.com/golang-jwt/jwt/v4"
)

func GenerateUserToken(secret string, id uint) (string, error) {

	claims := struct {
		jwt.StandardClaims
		UserID uint
	}{
		UserID: id,
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &claims)
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil

}
