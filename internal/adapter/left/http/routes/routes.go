package routes

import (
	"github.com/gin-gonic/gin"
	authhandlers "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/handlers/auth_handlers"
	listinghandlers "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/handlers/listing_handlers"
	userhandlers "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/handlers/user_handlers"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/middlewares"
	"github.com/giulio-alfieri/toq_server/internal/core/factory"
	goroutines "github.com/giulio-alfieri/toq_server/internal/core/go_routines"
	httpport "github.com/giulio-alfieri/toq_server/internal/core/port/left/http"
	metricsport "github.com/giulio-alfieri/toq_server/internal/core/port/right/metrics"
	permissionservice "github.com/giulio-alfieri/toq_server/internal/core/service/permission_service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes configura todas as rotas HTTP com middlewares e dependências injetadas
func SetupRoutes(
	router *gin.Engine,
	handlers *factory.HTTPHandlers,
	activityTracker *goroutines.ActivityTracker,
	permissionService permissionservice.PermissionServiceInterface,
	metricsAdapter *factory.MetricsAdapter,
	versionProvider httpport.APIVersionProvider,
) {
	// Configurar middlewares globais na ordem correta
	setupGlobalMiddlewares(router, metricsAdapter)

	// Swagger documentation routes
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Convert handlers to typed handlers
	authHandler := handlers.AuthHandler.(*authhandlers.AuthHandler)
	userHandler := handlers.UserHandler.(*userhandlers.UserHandler)
	listingHandler := handlers.ListingHandler.(*listinghandlers.ListingHandler)

	// API base routes (v2)
	base := "/api/v2"
	if versionProvider != nil {
		base = versionProvider.BasePath()
	}
	v1 := router.Group(base)

	// Register user routes with dependencies
	RegisterUserRoutes(v1, authHandler, userHandler, activityTracker, permissionService)

	// Register listing routes with dependencies
	RegisterListingRoutes(v1, listingHandler, activityTracker, permissionService)
}

// setupGlobalMiddlewares configura middlewares aplicados a todas as rotas
func setupGlobalMiddlewares(router *gin.Engine, metricsAdapter *factory.MetricsAdapter) {
	// Ordem específica dos middlewares para otimização e segurança

	// 1. RequestIDMiddleware - Gera ID único para cada request
	router.Use(middlewares.RequestIDMiddleware())

	// 2. Recovery - Captura panics
	router.Use(gin.Recovery())

	// 3. StructuredLoggingMiddleware - Log estruturado JSON com separação stdout/stderr
	router.Use(middlewares.StructuredLoggingMiddleware())

	// 4. CORSMiddleware - Configuração CORS
	router.Use(middlewares.CORSMiddleware())

	// 5. TelemetryMiddleware - Tracing OpenTelemetry + Métricas
	var metricsPort metricsport.MetricsPortInterface
	if metricsAdapter != nil {
		metricsPort = metricsAdapter.Prometheus
	}
	router.Use(middlewares.TelemetryMiddleware(metricsPort))

	// Nota: AuthMiddleware e PermissionMiddleware são aplicados apenas em rotas específicas
}

// RegisterUserRoutes registers all user-related routes with middleware dependencies
func RegisterUserRoutes(
	router *gin.RouterGroup,
	authHandler *authhandlers.AuthHandler,
	userHandler *userhandlers.UserHandler,
	activityTracker *goroutines.ActivityTracker,
	permissionService permissionservice.PermissionServiceInterface,
) {
	// Authentication routes (public - without auth middleware)
	auth := router.Group("/auth")
	{
		// CreateOwner
		auth.POST("/owner", authHandler.CreateOwner) // CreateOwner

		// CreateRealtor
		auth.POST("/realtor", authHandler.CreateRealtor) // CreateRealtor

		// CreateAgency
		auth.POST("/agency", authHandler.CreateAgency) // CreateAgency

		// SignIn
		auth.POST("/signin", authHandler.SignIn) // SignIn

		// RefreshToken
		auth.POST("/refresh", authHandler.RefreshToken) // RefreshToken

		// RequestPasswordChange
		auth.POST("/password/request", authHandler.RequestPasswordChange) // RequestPasswordChange

		// ConfirmPasswordChange
		auth.POST("/password/confirm", authHandler.ConfirmPasswordChange) // ConfirmPasswordChange
	}

	// User routes (authenticated - Owner, Realtor, Agency and Admin)
	user := router.Group("/user")
	// Apply security middlewares to authenticated routes
	user.Use(middlewares.AuthMiddleware(activityTracker))
	user.Use(middlewares.PermissionMiddleware(permissionService))
	{
		// Profile management
		user.GET("/profile", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })    // GetProfile
		user.PUT("/profile", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })    // UpdateProfile
		user.DELETE("/account", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // DeleteAccount
		user.GET("/onboarding", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // GetOnboardingStatus
		user.GET("/roles", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })      // GetUserRoles
		user.GET("/home", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })       // GoHome
		user.PUT("/opt-status", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // UpdateOptStatus

		// SignOut (authenticated endpoint)
		user.POST("/signout", userHandler.SignOut) // SignOut

		// Photo management
		user.POST("/photo/upload-url", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })  // GetPhotoUploadURL
		user.GET("/profile/thumbnails", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // GetProfileThumbnails

		// Email change workflow
		user.POST("/email/request", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // RequestEmailChange
		user.POST("/email/confirm", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // ConfirmEmailChange
		user.POST("/email/resend", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })  // ResendEmailChangeCode

		// Phone change workflow
		user.POST("/phone/request", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // RequestPhoneChange
		user.POST("/phone/confirm", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // ConfirmPhoneChange
		user.POST("/phone/resend", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })  // ResendPhoneChangeCode

		// Role management (Owner and Realtor only)
		user.POST("/role/alternative", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // AddAlternativeUserRole
		user.POST("/role/switch", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })      // SwitchUserRole
	}

	// Agency routes (Agency only)
	agency := router.Group("/agency")
	// Apply security middlewares to authenticated routes
	agency.Use(middlewares.AuthMiddleware(activityTracker))
	agency.Use(middlewares.PermissionMiddleware(permissionService))
	{
		agency.POST("/invite-realtor", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // InviteRealtor
		agency.GET("/realtors", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })        // GetRealtorsByAgency
		agency.GET("/realtors/:id", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })    // GetRealtorByID
		agency.DELETE("/realtors/:id", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // DeleteRealtorByID
	}

	// Realtor routes (Realtor only)
	realtor := router.Group("/realtor")
	// Apply security middlewares to authenticated routes
	realtor.Use(middlewares.AuthMiddleware(activityTracker))
	realtor.Use(middlewares.PermissionMiddleware(permissionService))
	{
		realtor.POST("/creci/verify", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })      // VerifyCreciImages
		realtor.POST("/creci/upload-url", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })  // GetCreciUploadURL
		realtor.POST("/invitation/accept", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // AcceptInvitation
		realtor.POST("/invitation/reject", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // RejectInvitation
		realtor.GET("/agency", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })             // GetAgencyOfRealtor
		realtor.DELETE("/agency", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })          // DeleteAgencyOfRealtor
	}
}

// RegisterListingRoutes registers all listing-related routes with middleware dependencies
func RegisterListingRoutes(
	router *gin.RouterGroup,
	listingHandler *listinghandlers.ListingHandler,
	activityTracker *goroutines.ActivityTracker,
	permissionService permissionservice.PermissionServiceInterface,
) {
	// Main listing routes (all require authentication)
	listings := router.Group("/listings")
	// Apply security middlewares to authenticated routes
	listings.Use(middlewares.AuthMiddleware(activityTracker))
	listings.Use(middlewares.PermissionMiddleware(permissionService))
	{
		// Basic CRUD
		listings.GET("", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })               // GetAllListings
		listings.POST("", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })              // StartListing
		listings.GET("/search", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })        // SearchListing
		listings.GET("/options", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })       // GetOptions
		listings.GET("/features/base", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // GetBaseFeatures

		// Favorites (Realtor side)
		listings.GET("/favorites", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // GetFavoriteListings

		// Individual listing operations
		listings.GET("/:id", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })             // GetListing
		listings.PUT("/:id", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })             // UpdateListing
		listings.DELETE("/:id", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })          // DeleteListing
		listings.POST("/:id/end-update", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // EndUpdateListing
		listings.GET("/:id/status", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })      // GetListingStatus

		// Owner operations
		listings.POST("/:id/approve", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // ApproveListing
		listings.POST("/:id/reject", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })  // RejectListing
		listings.POST("/:id/suspend", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // SuspendListing
		listings.POST("/:id/release", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // ReleaseListing
		listings.POST("/:id/copy", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })    // CopyListing

		// Realtor operations
		listings.POST("/:id/share", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })      // ShareListing
		listings.POST("/:id/favorite", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })   // AddFavoriteListing
		listings.DELETE("/:id/favorite", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // RemoveFavoriteListing

		// Visit requests
		listings.POST("/:id/visit/request", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // RequestVisit
		listings.GET("/:id/visits", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })         // GetVisits

		// Offers
		listings.POST("/:id/offers", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // CreateOffer
		listings.GET("/:id/offers", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })  // GetOffers
	}

	// Visit management (all require authentication)
	visits := router.Group("/visits")
	// Apply security middlewares to authenticated routes
	visits.Use(middlewares.AuthMiddleware(activityTracker))
	visits.Use(middlewares.PermissionMiddleware(permissionService))
	{
		visits.GET("", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })              // GetAllVisits
		visits.DELETE("/:id", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })       // CancelVisit
		visits.POST("/:id/confirm", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // ConfirmVisitDone
		visits.POST("/:id/approve", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // ApproveVisiting
		visits.POST("/:id/reject", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })  // RejectVisiting
	}

	// Offer management (all require authentication)
	offers := router.Group("/offers")
	// Apply security middlewares to authenticated routes
	offers.Use(middlewares.AuthMiddleware(activityTracker))
	offers.Use(middlewares.PermissionMiddleware(permissionService))
	{
		offers.GET("", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })              // GetAllOffers
		offers.PUT("/:id", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })          // UpdateOffer
		offers.DELETE("/:id", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })       // CancelOffer
		offers.POST("/:id/send", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })    // SendOffer
		offers.POST("/:id/approve", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // ApproveOffer
		offers.POST("/:id/reject", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })  // RejectOffer
	}

	// Evaluation routes (all require authentication)
	realtors := router.Group("/realtors")
	// Apply security middlewares to authenticated routes
	realtors.Use(middlewares.AuthMiddleware(activityTracker))
	realtors.Use(middlewares.PermissionMiddleware(permissionService))
	{
		realtors.POST("/:id/evaluate", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // EvaluateRealtor
	}

	owners := router.Group("/owners")
	// Apply security middlewares to authenticated routes
	owners.Use(middlewares.AuthMiddleware(activityTracker))
	owners.Use(middlewares.PermissionMiddleware(permissionService))
	{
		owners.POST("/:id/evaluate", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // EvaluateOwner
	}
}
