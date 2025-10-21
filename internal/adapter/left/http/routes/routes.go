package routes

import (
	"github.com/gin-gonic/gin"
	adminhandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/admin_handlers"
	authhandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/auth_handlers"
	complexhandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/complex_handlers"
	holidayhandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/holiday_handlers"
	listinghandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/listing_handlers"
	schedulehandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/schedule_handlers"
	userhandlers "github.com/projeto-toq/toq_server/internal/adapter/left/http/handlers/user_handlers"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	"github.com/projeto-toq/toq_server/internal/core/factory"
	goroutines "github.com/projeto-toq/toq_server/internal/core/go_routines"
	httpport "github.com/projeto-toq/toq_server/internal/core/port/left/http"
	metricsport "github.com/projeto-toq/toq_server/internal/core/port/right/metrics"
	permissionservice "github.com/projeto-toq/toq_server/internal/core/service/permission_service"
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

	// Simple test route
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Serve swagger.json with CORS middleware applied
	router.GET("/docs/swagger.json", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.File("/codigos/go_code/toq_server/docs/swagger.json")
	})

	// Swagger documentation routes
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // Convert handlers to typed handlers
	authHandler := handlers.AuthHandler.(*authhandlers.AuthHandler)
	userHandler := handlers.UserHandler.(*userhandlers.UserHandler)
	listingHandler := handlers.ListingHandler.(*listinghandlers.ListingHandler)
	adminHandler := handlers.AdminHandler.(*adminhandlers.AdminHandler)
	complexHandler := handlers.ComplexHandler.(*complexhandlers.ComplexHandler)
	scheduleHandler := handlers.ScheduleHandler.(*schedulehandlers.ScheduleHandler)
	holidayHandler := handlers.HolidayHandler.(*holidayhandlers.HolidayHandler)
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

	// Register admin routes with dependencies
	RegisterAdminRoutes(v1, adminHandler, holidayHandler, activityTracker, permissionService)

	// Register complex routes (authenticated)
	RegisterComplexRoutes(v1, complexHandler, activityTracker, permissionService)

	// Register schedule routes (authenticated)
	RegisterScheduleRoutes(v1, scheduleHandler, activityTracker, permissionService)
}

// setupGlobalMiddlewares configura middlewares aplicados a todas as rotas
func setupGlobalMiddlewares(router *gin.Engine, metricsAdapter *factory.MetricsAdapter) {
	// Ordem específica dos middlewares para otimização e segurança

	// 1. RequestIDMiddleware - Gera ID único para cada request
	router.Use(middlewares.RequestIDMiddleware())

	// 2. StructuredLoggingMiddleware - Log estruturado JSON com separação stdout/stderr
	router.Use(middlewares.StructuredLoggingMiddleware())

	// 3. CORSMiddleware - Configuração CORS
	router.Use(middlewares.CORSMiddleware())

	// 4. TelemetryMiddleware - Tracing OpenTelemetry + Métricas
	var metricsPort metricsport.MetricsPortInterface
	if metricsAdapter != nil {
		metricsPort = metricsAdapter.Prometheus
		// Store metrics adapter reference in Gin context for other middlewares
		router.Use(func(c *gin.Context) {
			c.Set("metricsAdapter", metricsPort)
			c.Next()
		})
	}
	router.Use(middlewares.TelemetryMiddleware(metricsPort))

	// 5. ErrorRecoveryMiddleware - Captura panics (após Telemetry para marcar spans)
	router.Use(middlewares.ErrorRecoveryMiddleware())

	// 6. DeviceContextMiddleware - injeta DeviceID no contexto
	router.Use(middlewares.DeviceContextMiddleware())

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
		// Validation endpoints (public signed requests)
		auth.POST("/validate/cpf", authHandler.ValidateCPF)
		auth.POST("/validate/cnpj", authHandler.ValidateCNPJ)
		auth.POST("/validate/cep", authHandler.ValidateCEP)

		// CreateOwner
		auth.POST("/owner", authHandler.CreateOwner) // CreateOwner

		// CreateRealtor
		auth.POST("/realtor", authHandler.CreateRealtor) // CreateRealtor

		// CreateAgency
		auth.POST("/agency", authHandler.CreateAgency) // CreateAgency

		// SignIn
		auth.POST("/signin", middlewares.RequireDeviceIDMiddleware(), authHandler.SignIn) // SignIn
	}

	// User routes (authenticated)
	user := router.Group("/user")
	user.Use(middlewares.AuthMiddleware(activityTracker))
	user.Use(middlewares.PermissionMiddleware(permissionService))
	{
		user.GET("/roles", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // GetUserRoles
		user.GET("/home", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })  // GoHome
		user.PUT("/opt-status", userHandler.UpdateOptStatus)                                            // UpdateOptStatus

		// SignOut (authenticated endpoint)
		user.POST("/signout", userHandler.SignOut) // SignOut

		// Photo management
		user.POST("/photo/upload-url", userHandler.PostPhotoUploadURL)     // PostPhotoUploadURL
		user.POST("/photo/download-url", userHandler.PostPhotoDownloadURL) // PostPhotoDownloadURL

		// Email change workflow
		user.POST("/email/request", userHandler.RequestEmailChange)
		user.POST("/email/confirm", userHandler.ConfirmEmailChange)
		user.POST("/email/resend", userHandler.ResendEmailChangeCode)

		// Phone change workflow
		user.POST("/phone/request", userHandler.RequestPhoneChange)
		user.POST("/phone/confirm", userHandler.ConfirmPhoneChange)
		user.POST("/phone/resend", userHandler.ResendPhoneChangeCode)

		// Role management (Owner and Realtor only)
		user.POST("/role/alternative", userHandler.AddAlternativeUserRole) // AddAlternativeUserRole
		user.POST("/role/switch", userHandler.SwitchUserRole)              // SwitchUserRole
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
		realtor.POST("/creci/verify", userHandler.VerifyCreciDocuments)                                                 // VerifyCreciDocuments
		realtor.POST("/creci/upload-url", userHandler.PostCreciUploadURL)                                               // PostCreciUploadURL
		realtor.POST("/invitation/accept", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // AcceptInvitation
		realtor.POST("/invitation/reject", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // RejectInvitation
		realtor.GET("/agency", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })             // GetAgencyOfRealtor
		realtor.DELETE("/agency", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })          // DeleteAgencyOfRealtor
	}
}

// RegisterScheduleRoutes registers routes related to listing schedules and availability.
func RegisterScheduleRoutes(
	router *gin.RouterGroup,
	scheduleHandler *schedulehandlers.ScheduleHandler,
	activityTracker *goroutines.ActivityTracker,
	permissionService permissionservice.PermissionServiceInterface,
) {
	schedules := router.Group("/schedules")
	schedules.Use(middlewares.AuthMiddleware(activityTracker))
	schedules.Use(middlewares.PermissionMiddleware(permissionService))
	{
		schedules.POST("/owner/summary", scheduleHandler.PostOwnerSummary)
		schedules.POST("/listing/detail", scheduleHandler.PostListingAgenda)
		schedules.POST("/listing/block", scheduleHandler.PostCreateBlockEntry)
		schedules.PUT("/listing/block", scheduleHandler.PutUpdateBlockEntry)
		schedules.DELETE("/listing/block", scheduleHandler.DeleteBlockEntry)
		schedules.POST("/listing/availability", scheduleHandler.PostListAvailability)
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
		listings.GET("", listingHandler.GetAllListings)                                                      // GetAllListings (returns 501 until service ready)
		listings.POST("", listingHandler.StartListing)                                                       // StartListing
		listings.GET("/search", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // SearchListing
		listings.POST("/options", listingHandler.PostOptions)                                                // PostOptions
		listings.GET("/features/base", listingHandler.GetBaseFeatures)                                       // GetBaseFeatures
		listings.GET("/catalog", listingHandler.ListCatalogValues)                                           // ListCatalogValues

		// Favorites (Realtor side)
		listings.GET("/favorites", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // GetFavoriteListings

		// Photo session scheduling
		listings.GET("/photo-session/slots", listingHandler.ListPhotographerSlots)
		listings.POST("/photo-session/reserve", listingHandler.ReservePhotoSession)
		listings.POST("/photo-session/confirm", listingHandler.ConfirmPhotoSession)
		listings.POST("/photo-session/cancel", listingHandler.CancelPhotoSession)

		// Individual listing operations
		listings.GET("/detail", listingHandler.GetListing)                                                       // GetListing
		listings.GET("/:id", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })        // GetListing
		listings.PUT("", listingHandler.UpdateListing)                                                           // UpdateListing
		listings.DELETE("/:id", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })     // DeleteListing
		listings.POST("/end-update", listingHandler.EndUpdateListing)                                            // EndUpdateListing
		listings.GET("/:id/status", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // GetListingStatus

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

// RegisterAdminRoutes registers admin-related routes with authentication and admin permission
func RegisterAdminRoutes(
	router *gin.RouterGroup,
	adminHandler *adminhandlers.AdminHandler,
	holidayHandler *holidayhandlers.HolidayHandler,
	activityTracker *goroutines.ActivityTracker,
	permissionService permissionservice.PermissionServiceInterface,
) {
	admin := router.Group("/admin")
	// Apply security middlewares: Auth and Admin permission check
	admin.Use(middlewares.AuthMiddleware(activityTracker))
	admin.Use(middlewares.PermissionMiddleware(permissionService))
	admin.Use(middlewares.RequireAdminPermission(permissionService))

	{
		// GET /admin/users/creci/pending
		admin.GET("/users/creci/pending", adminHandler.GetPendingRealtors)
		admin.GET("/users", adminHandler.GetAdminUsers)
		admin.GET("/listing/catalog", adminHandler.ListListingCatalogValues)
		admin.POST("/listing/catalog/detail", adminHandler.PostAdminGetListingCatalogDetail)
		admin.POST("/listing/catalog", adminHandler.CreateListingCatalogValue)
		admin.PUT("/listing/catalog", adminHandler.UpdateListingCatalogValue)
		admin.POST("/listing/catalog/restore", adminHandler.RestoreListingCatalogValue)
		admin.DELETE("/listing/catalog", adminHandler.DeleteListingCatalogValue)
		admin.GET("/complexes", adminHandler.GetAdminComplexes)
		admin.POST("/complexes", adminHandler.PostAdminCreateComplex)
		admin.PUT("/complexes", adminHandler.PutAdminUpdateComplex)
		admin.DELETE("/complexes", adminHandler.DeleteAdminComplex)
		admin.POST("/complexes/detail", adminHandler.PostAdminGetComplexDetail)
		admin.GET("/complexes/towers", adminHandler.GetAdminComplexTowers)
		admin.POST("/complexes/towers", adminHandler.PostAdminCreateComplexTower)
		admin.PUT("/complexes/towers", adminHandler.PutAdminUpdateComplexTower)
		admin.DELETE("/complexes/towers", adminHandler.DeleteAdminComplexTower)
		admin.POST("/complexes/towers/detail", adminHandler.PostAdminGetComplexTowerDetail)
		admin.GET("/complexes/sizes", adminHandler.GetAdminComplexSizes)
		admin.POST("/complexes/sizes", adminHandler.PostAdminCreateComplexSize)
		admin.PUT("/complexes/sizes", adminHandler.PutAdminUpdateComplexSize)
		admin.DELETE("/complexes/sizes", adminHandler.DeleteAdminComplexSize)
		admin.POST("/complexes/sizes/detail", adminHandler.PostAdminGetComplexSizeDetail)
		admin.GET("/complexes/zip-codes", adminHandler.GetAdminComplexZipCodes)
		admin.POST("/complexes/zip-codes", adminHandler.PostAdminCreateComplexZipCode)
		admin.PUT("/complexes/zip-codes", adminHandler.PutAdminUpdateComplexZipCode)
		admin.DELETE("/complexes/zip-codes", adminHandler.DeleteAdminComplexZipCode)
		admin.POST("/complexes/zip-codes/detail", adminHandler.PostAdminGetComplexZipCodeDetail)
		admin.POST("/users/system", adminHandler.PostAdminCreateSystemUser)
		admin.PUT("/users/system", adminHandler.PutAdminUpdateSystemUser)
		admin.DELETE("/users/system", adminHandler.DeleteAdminSystemUser)
		admin.GET("/roles", adminHandler.GetAdminRoles)
		admin.POST("/roles/detail", adminHandler.PostAdminGetRoleDetail)
		admin.POST("/roles", adminHandler.PostAdminCreateRole)
		admin.PUT("/roles", adminHandler.PutAdminUpdateRole)
		admin.POST("/roles/restore", adminHandler.RestoreAdminRole)
		admin.DELETE("/roles", adminHandler.DeleteAdminRole)
		admin.GET("/permissions", adminHandler.GetAdminPermissions)
		admin.POST("/permissions/detail", adminHandler.PostAdminGetPermissionDetail)
		admin.POST("/permissions", adminHandler.PostAdminCreatePermission)
		admin.PUT("/permissions", adminHandler.PutAdminUpdatePermission)
		admin.DELETE("/permissions", adminHandler.DeleteAdminPermission)
		admin.GET("/role-permissions", adminHandler.GetAdminRolePermissions)
		admin.POST("/role-permissions/detail", adminHandler.PostAdminGetRolePermissionDetail)
		admin.POST("/role-permissions", adminHandler.PostAdminCreateRolePermission)
		admin.PUT("/role-permissions", adminHandler.PutAdminUpdateRolePermission)
		admin.DELETE("/role-permissions", adminHandler.DeleteAdminRolePermission)

		// POST /admin/users/detail
		admin.POST("/users/detail", adminHandler.PostAdminGetUser)

		// POST /admin/users/creci/approve
		admin.POST("/users/creci/approve", adminHandler.PostAdminApproveUser)

		// POST /admin/users/creci/download-url
		admin.POST("/users/creci/download-url", adminHandler.PostAdminCreciDownloadURL)
	}

	holidayGroup := admin.Group("/holidays")
	{
		holidayGroup.GET("/calendars", holidayHandler.ListCalendars)
		holidayGroup.POST("/calendars/detail", holidayHandler.GetCalendarDetail)
		holidayGroup.POST("/calendars", holidayHandler.CreateCalendar)
		holidayGroup.PUT("/calendars", holidayHandler.UpdateCalendar)
		holidayGroup.POST("/dates", holidayHandler.CreateCalendarDate)
		holidayGroup.GET("/dates", holidayHandler.ListCalendarDates)
		holidayGroup.DELETE("/dates", holidayHandler.DeleteCalendarDate)
	}
}

// RegisterComplexRoutes configura rotas autenticadas relacionadas aos empreendimentos.
func RegisterComplexRoutes(
	router *gin.RouterGroup,
	complexHandler *complexhandlers.ComplexHandler,
	activityTracker *goroutines.ActivityTracker,
	permissionService permissionservice.PermissionServiceInterface,
) {
	complex := router.Group("/complex")
	complex.Use(middlewares.AuthMiddleware(activityTracker))
	complex.Use(middlewares.PermissionMiddleware(permissionService))
	{
		complex.GET("/sizes", complexHandler.ListSizesByAddress)
	}
}
