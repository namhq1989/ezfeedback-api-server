package manipulation

import (
	"math/rand"
	"time"
)

const (
	charset   = "abcdefghijklmnopqrstuvwxyz0123456789"
	otpDigits = "0123456789"
)

func RandomIntInRange(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandomAlphaNumeric(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}

func GenerateOTP(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = otpDigits[r.Intn(len(otpDigits))]
	}

	return string(result)
}
