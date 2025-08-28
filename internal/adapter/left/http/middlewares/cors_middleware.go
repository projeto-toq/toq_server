package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware configures Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins: []string{
			"https://gca.dev.br",
			"https://www.gca.dev.br",
			"http://localhost:3000", // Development
			"http://localhost:5173", // Vite dev server
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-Request-ID",
			"User-Agent",
			"X-Forwarded-For",
			"X-Real-IP",
		},
		ExposeHeaders: []string{
			"X-Request-ID",
			"X-API-Version",
		},
		AllowCredentials: true,
		MaxAge:           43200, // 12 hours
	}

	return cors.New(config)
}
