package tokens

import (
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func GenerateConfirmationToken(length int64) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

var secretKey = []byte(os.Getenv("SECRET_KEY"))

func CreateAccessToken(userId int64) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(int(userId)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func VerifyAccessToken(accessToken string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not parse token")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return &claims, nil
	}
	return nil, fmt.Errorf("could not parse claims")
}
