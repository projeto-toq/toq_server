package userservices

import (
	"time"

	"golang.org/x/exp/rand"
)

const charset = "ABCDEFGHIJKLMNPQRSTUVWXYZ123456789"

func (us *userService) random6Digits() (random string) {

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(uint64(time.Now().UTC().UnixNano())))

	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]

	}

	return string(b)

}
