// Common tools and helper functions
package utils

import (
	"math/rand"
	"time"
)

func GenCaptcha() string {
	digits := []byte("0123456789")
	b := make([]byte, 6)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = digits[rand.Intn(10)]
	}
	return string(b)
}
