package utils

import (
	"crypto/rand"
	"math/big"
)

const chars = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"0123456789"

func RandomString(n int) (string, error) {
	result := make([]byte, n)
	max := big.NewInt(int64(len(chars)))
	for i := range result {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		result[i] = chars[num.Int64()]
	}
	return string(result), nil
}
