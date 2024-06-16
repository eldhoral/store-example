package security

import (
	"math/rand"
	"time"
)

const (
	ALPHANUMERIC = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
	NUMERIC      = "1234567890"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func sringWithCharset(length int, charset string) string {
	if length < 1 {
		return ""
	}

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateRandomStringYii(length int) string {
	return sringWithCharset(length, ALPHANUMERIC)
}

func GenerateRandomStringNumeric(length int) string {
	return sringWithCharset(length, NUMERIC)
}
