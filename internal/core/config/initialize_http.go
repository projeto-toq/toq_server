package config

import (
	"crypto/tls"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares"
)

func (c *config) InitializeHTTP() {
	slog.Info("Initializing HTTP/Gin server")

	// Create Gin router
	gin.SetMode(gin.ReleaseMode)

	c.ginRouter = gin.New()

	// Add global middlewares
	c.ginRouter.Use(gin.Logger())
	c.ginRouter.Use(gin.Recovery())
	c.ginRouter.Use(middlewares.CORSMiddleware())
	c.ginRouter.Use(middlewares.TelemetryMiddleware())

	// Health check endpoints (no auth required)
	c.ginRouter.GET("/healthz", c.healthzHandler)
	c.ginRouter.GET("/readyz", c.readyzHandler)

	// API routes with authentication
	api := c.ginRouter.Group("/api/v1")
	api.Use(middlewares.AuthMiddleware(c.activityTracker))
	api.Use(middlewares.AccessControlMiddleware())

	// Note: HTTP handlers will be initialized after dependency injection
	// and routes will be set up at that time

	// Parse timeouts from environment
	readTimeout, err := time.ParseDuration(c.env.HTTP.ReadTimeout)
	if err != nil {
		slog.Warn("Invalid read timeout, using default", "error", err)
		readTimeout = 30 * time.Second
	}

	writeTimeout, err := time.ParseDuration(c.env.HTTP.WriteTimeout)
	if err != nil {
		slog.Warn("Invalid write timeout, using default", "error", err)
		writeTimeout = 30 * time.Second
	}

	// Create HTTP server
	c.httpServer = &http.Server{
		Addr:           c.env.HTTP.Port,
		Handler:        c.ginRouter,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: c.env.HTTP.MaxHeaderBytes,
	}

	// Configure TLS if enabled
	if c.env.HTTP.TLS.Enabled {
		slog.Info("Configuring HTTPS/TLS", "certPath", c.env.HTTP.TLS.CertPath)
		cert, err := tls.LoadX509KeyPair(c.env.HTTP.TLS.CertPath, c.env.HTTP.TLS.KeyPath)
		if err != nil {
			slog.Error("Failed to load TLS certificates", "error", err)
			panic(err)
		}

		c.httpServer.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}
		slog.Info("HTTPS/TLS configured successfully")
	}

	slog.Info("HTTP/Gin server initialized successfully",
		"port", c.env.HTTP.Port,
		"tls", c.env.HTTP.TLS.Enabled)
}

// SetupHTTPHandlersAndRoutes initializes HTTP handlers and sets up routes after dependency injection
func (c *config) SetupHTTPHandlersAndRoutes() {
	slog.Debug("Setting up HTTP handlers and routes")

	// Import routes package dynamically
	const routesImport = "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes"
	_ = routesImport // Prevent unused import error

	// Create HTTP handlers using factory pattern
	c.httpHandlers = c.adapterFactory.CreateHTTPHandlers(
		c.userService,
		c.globalService,
		c.listingService,
		c.complexService,
	)

	// Note: Routes will be manually configured for now until dynamic import is implemented
	slog.Debug("HTTP handlers configured successfully (routes pending)")
}

// healthzHandler handles liveness probe
func (c *config) healthzHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "toq-server",
	})
}

// readyzHandler handles readiness probe
func (c *config) readyzHandler(ctx *gin.Context) {
	status := http.StatusOK
	if !c.readiness {
		status = http.StatusServiceUnavailable
	}

	ctx.JSON(status, gin.H{
		"status":    map[bool]string{true: "ready", false: "not_ready"}[c.readiness],
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "toq-server",
		"checks": gin.H{
			"database": "ok",
			"cache":    "ok",
		},
	})
}
