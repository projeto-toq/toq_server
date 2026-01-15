package globalmodel

import (
	"sync"
	"time"
)

type contextKey string

const (
	// Access token expiration default (override via env) - reduce to 15m
	DefaultAccessTokenExpiration = 15 * time.Minute

	// Default refresh token expiration (override via env)
	DefaultRefreshTokenExpiration = 30 * 24 * time.Hour

	//context keys
	RequestIDKey               contextKey = "requestID"
	TokenKey                   contextKey = "infos"
	UserAgentKey               contextKey = "userAgent"
	ClientIPKey                contextKey = "clientIP"
	DeviceIDKey                contextKey = "deviceID"
	AccessTokenJTIKey          contextKey = "accessTokenJTI"
	AccessTokenExpiresAtKey    contextKey = "accessTokenExpiresAt"
	SessionAbsoluteExpiryKey   contextKey = "sessionAbsoluteExpiry"
	SessionRotationCounterKey  contextKey = "sessionRotationCounter"
	MaxSessionRotationsDefault            = 10
)

var once sync.Once
var jwtSecret string
var refreshTTL = DefaultRefreshTokenExpiration
var accessTTL = DefaultAccessTokenExpiration
var maxSessionRotations = MaxSessionRotationsDefault

func SetJWTSecret(secret string) {
	once.Do(func() {
		jwtSecret = secret
	})
}

func GetJWTSecret() string {
	return jwtSecret
}

func SetRefreshTTL(days int) {
	if days <= 0 {
		return
	}
	refreshTTL = time.Duration(days) * 24 * time.Hour
}

func GetRefreshTTL() time.Duration { return refreshTTL }

func SetAccessTTL(minutes int) {
	if minutes <= 0 {
		return
	}
	accessTTL = time.Duration(minutes) * time.Minute
}

func GetAccessTTL() time.Duration { return accessTTL }

func SetMaxSessionRotations(n int) {
	if n > 0 {
		maxSessionRotations = n
	}
}

func GetMaxSessionRotations() int { return maxSessionRotations }
