package userservices

import (
	"log/slog"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/golang-jwt/jwt"
)

func validateRefreshToken(refresh string) (userID int64, err error) {

	//tenta validar o token
	token, err2 := jwt.Parse(refresh, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			slog.Warn("unexpected signing method", "method", token.Header["alg"])
			return nil, utils.ErrInternalServer
		}
		secret := globalmodel.GetJWTSecret()
		return []byte(secret), nil
	})

	if err2 != nil {
		return 0, utils.ErrInternalServer
	}

	//tenta recuperar os claims
	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, utils.ErrInternalServer
	}

	infosraw, ok := payload[string(globalmodel.TokenKey)].(map[string]interface{})
	if !ok {
		return 0, utils.ErrInternalServer
	}

	idFloat, ok := infosraw["ID"].(float64)
	if !ok {
		return 0, utils.ErrInternalServer
	}
	userID = int64(idFloat)

	return
}
