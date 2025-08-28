package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes"
	"github.com/giulio-alfieri/toq_server/internal/core/factory"
	goroutines "github.com/giulio-alfieri/toq_server/internal/core/go_routines"
	permissionservice "github.com/giulio-alfieri/toq_server/internal/core/service/permission_service"
)

// ServerConfig holds the HTTP server configuration
type ServerConfig struct {
	Port              string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	MaxHeaderBytes    int
	ActivityTracker   *goroutines.ActivityTracker
	HTTPHandlers      *factory.HTTPHandlers
	PermissionService permissionservice.PermissionServiceInterface
	// Cache will be added later when access control middleware is implemented
}

// Server represents the HTTP server
type Server struct {
	config     *ServerConfig
	httpServer *http.Server
	router     *gin.Engine
}

// NewServer creates a new HTTP server instance
func NewServer(config *ServerConfig) *Server {
	// Set Gin mode to release for production
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	router := gin.New()

	// Add global middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(requestid.New())
	router.Use(middlewares.TelemetryMiddleware())
	router.Use(middlewares.CORSMiddleware())

	// Health check endpoints (no auth required)
	router.GET("/healthz", healthzHandler)
	router.GET("/readyz", readyzHandler)

	// API routes with authentication and permissions
	api := router.Group("/api/v1")
	api.Use(middlewares.AuthMiddleware(config.ActivityTracker))
	api.Use(middlewares.PermissionMiddleware(config.PermissionService))

	// Register API routes
	routes.RegisterUserRoutes(api, config.HTTPHandlers)
	routes.RegisterListingRoutes(api, config.HTTPHandlers)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:           config.Port,
		Handler:        router,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}

	return &Server{
		config:     config,
		httpServer: httpServer,
		router:     router,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	slog.Info("Starting HTTP server", "port", s.config.Port)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shutdowns the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("Shutting down HTTP server")
	return s.httpServer.Shutdown(ctx)
}

// GetServer returns the underlying HTTP server
func (s *Server) GetServer() *http.Server {
	return s.httpServer
}

// GetRouter returns the Gin router
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

// healthzHandler handles liveness probe
func healthzHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "toq-server",
	})
}

// readyzHandler handles readiness probe
func readyzHandler(c *gin.Context) {
	// TODO: Add actual readiness checks (database, cache, etc.)
	c.JSON(http.StatusOK, gin.H{
		"status":    "ready",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "toq-server",
		"checks": gin.H{
			"database": "ok",
			"cache":    "ok",
		},
	})
}
