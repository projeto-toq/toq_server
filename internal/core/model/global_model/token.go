package globalmodel

import (
	"sync"
	"time"
)

type contextKey string //TODO ajustar criondo tipos para todas as constantes

const (

	//expiraton time for validation tokens
	TokenExpiration = 5 * time.Hour

	//context keys
	RequestIDKey contextKey = "requestID"
	TokenKey     contextKey = "infos"
)

var once sync.Once
var jwtSecret string

func SetJWTSecret(secret string) {
	once.Do(func() {
		jwtSecret = secret
	})
}

func GetJWTSecret() string {
	return jwtSecret
}
