package userservices

import (
	"math/rand"
	"time"
)

const charset = "ABCDEFGHIJKLMNPQRSTUVWXYZ123456789"

func (us *userService) random6Digits() (random string) {
	seededRand := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}
