package otp

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
)

func GenerateOTP() (string, error) {
	upperBound := big.NewInt(9000)
	randomNumber, err := rand.Int(rand.Reader, upperBound)
	if err != nil {
		return "", fmt.Errorf("error generating random number: %w", err)
	}
	result := new(big.Int).Add(randomNumber, big.NewInt(1000))
	code := strconv.Itoa(int(result.Int64()))
	return code, nil
}
