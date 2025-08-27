package config

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/routes"
	"github.com/giulio-alfieri/toq_server/internal/core/factory"
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
	c.ginRouter.Use(middlewares.AuthMiddleware(c.activityTracker))
	c.ginRouter.Use(middlewares.AccessControlMiddleware())

	// Initialize HTTP handlers with dependency injection
	c.initializeHTTPHandlers()

	// Setup routes with handlers
	routes.SetupRoutes(c.ginRouter, &c.httpHandlers)

	slog.Info("HTTP/Gin server initialized successfully")
}

func (c *config) initializeHTTPHandlers() {
	slog.Debug("Initializing HTTP handlers with dependency injection")

	// Create HTTP handlers using factory pattern
	c.httpHandlers = c.adapterFactory.CreateHTTPHandlers(
		c.userService,
		c.globalService,
		c.listingService,
		c.complexService,
	)

	slog.Debug("HTTP handlers initialized successfully")
}

func (c *config) GetGinRouter() *gin.Engine {
	return c.ginRouter
}

func (c *config) GetHTTPHandlers() *factory.HTTPHandlers {
	return &c.httpHandlers
}
