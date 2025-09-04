package userservices

import (
	"log/slog"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/golang-jwt/jwt"
)

func validateRefreshToken(refresh string) (userID int64, err error) {
	// Validação do token com verificação explícita de método de assinatura e tipagem
	token, parseErr := jwt.Parse(refresh, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			slog.Warn("jwt.unexpected_signing_method", "alg", token.Header["alg"])
			return nil, utils.WrapDomainErrorWithSource(utils.ErrInvalidRefreshToken)
		}
		secret := globalmodel.GetJWTSecret()
		return []byte(secret), nil
	})

	if parseErr != nil {
		if ve, ok := parseErr.(*jwt.ValidationError); ok {
			if (ve.Errors & jwt.ValidationErrorExpired) != 0 {
				return 0, utils.WrapDomainErrorWithSource(utils.ErrRefreshTokenExpired)
			}
		}
		return 0, utils.WrapDomainErrorWithSource(utils.ErrInvalidRefreshToken)
	}

	// Recupera claims e valida
	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, utils.WrapDomainErrorWithSource(utils.ErrInvalidRefreshToken)
	}

	// Verifica tipagem do token
	if typ, ok := payload["typ"].(string); !ok || typ != "refresh" {
		slog.Warn("jwt.invalid_type_for_refresh", "typ", payload["typ"]) // log de segurança
		return 0, utils.WrapDomainErrorWithSource(utils.ErrInvalidRefreshToken)
	}

	infosraw, ok := payload[string(globalmodel.TokenKey)].(map[string]interface{})
	if !ok {
		return 0, utils.WrapDomainErrorWithSource(utils.ErrInvalidRefreshToken)
	}

	idFloat, ok := infosraw["ID"].(float64)
	if !ok {
		return 0, utils.WrapDomainErrorWithSource(utils.ErrInvalidRefreshToken)
	}
	userID = int64(idFloat)

	return
}
