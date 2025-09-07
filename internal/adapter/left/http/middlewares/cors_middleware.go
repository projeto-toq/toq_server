package middlewares

import (
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware configura CORS aceitando subdomínios HTTPS de gca.dev.br e origens de desenvolvimento.
// Inclui cabeçalho X-Device-Id e permite credenciais.
func CORSMiddleware() gin.HandlerFunc {
	allowedStatic := map[string]struct{}{
		"https://gca.dev.br":            {},
		"https://www.gca.dev.br":        {},
		"https://api.gca.dev.br":        {},
		"https://swagger.gca.dev.br":    {},
		"https://grafana.gca.dev.br":    {},
		"https://jaeger.gca.dev.br":     {},
		"https://prometheus.gca.dev.br": {},
		// Dev
		"http://localhost:3000": {},
		"http://localhost:5173": {},
		"http://127.0.0.1:3000": {},
		"http://127.0.0.1:5173": {},
	}

	cfg := cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if origin == "" {
				return false
			}
			if _, ok := allowedStatic[origin]; ok {
				return true
			}
			if strings.HasPrefix(origin, "https://") && strings.HasSuffix(origin, ".gca.dev.br") {
				return true
			}
			return false
		},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Accept", "Authorization",
			"X-Requested-With", "X-Request-ID", "User-Agent",
			"X-Forwarded-For", "X-Real-IP", "X-Device-Id",
		},
		ExposeHeaders:    []string{"X-Request-ID", "X-API-Version"},
		AllowCredentials: true,
		MaxAge:           43200,
	}
	return cors.New(cfg)
}
