package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/core/factory"
)

// SetupRoutes configura todas as rotas HTTP
func SetupRoutes(router *gin.Engine, handlers *factory.HTTPHandlers) {
	// API v1 routes
	v1 := router.Group("/api/v1")

	// Register user routes
	RegisterUserRoutes(v1, handlers)

	// Register listing routes
	RegisterListingRoutes(v1, handlers)
}

// RegisterUserRoutes registers all user-related routes
func RegisterUserRoutes(router *gin.RouterGroup, handlers *factory.HTTPHandlers) {
	// TODO: Implement user routes in Etapa 2
	// This is a placeholder for now

	// Authentication routes (public - without auth middleware)
	auth := router.Group("/auth")
	{
		// CreateOwner godoc
		//	@Summary		Create owner account
		//	@Description	Create a new owner account with user information
		//	@Tags			Authentication
		//	@Accept			json
		//	@Produce		json
		//	@Param			request	body		dto.CreateOwnerRequest	true	"Owner creation data"
		//	@Success		201		{object}	dto.CreateOwnerResponse
		//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
		//	@Failure		409		{object}	dto.ErrorResponse	"User already exists"
		//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
		//	@Failure		501		{object}	dto.ErrorResponse	"Endpoint not implemented yet"
		//	@Router			/auth/owner [post]
		auth.POST("/owner", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })            // CreateOwner

		// CreateRealtor godoc
		//	@Summary		Create realtor account
		//	@Description	Create a new realtor account with user and CRECI information
		//	@Tags			Authentication
		//	@Accept			json
		//	@Produce		json
		//	@Param			request	body		dto.CreateRealtorRequest	true	"Realtor creation data"
		//	@Success		201		{object}	dto.CreateRealtorResponse
		//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
		//	@Failure		409		{object}	dto.ErrorResponse	"User already exists"
		//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
		//	@Failure		501		{object}	dto.ErrorResponse	"Endpoint not implemented yet"
		//	@Router			/auth/realtor [post]
		auth.POST("/realtor", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })          // CreateRealtor

		// CreateAgency godoc
		//	@Summary		Create agency account
		//	@Description	Create a new agency account with user information
		//	@Tags			Authentication
		//	@Accept			json
		//	@Produce		json
		//	@Param			request	body		dto.CreateAgencyRequest	true	"Agency creation data"
		//	@Success		201		{object}	dto.CreateAgencyResponse
		//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
		//	@Failure		409		{object}	dto.ErrorResponse	"User already exists"
		//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
		//	@Failure		501		{object}	dto.ErrorResponse	"Endpoint not implemented yet"
		//	@Router			/auth/agency [post]
		auth.POST("/agency", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })           // CreateAgency

		// SignIn godoc
		//	@Summary		User sign in
		//	@Description	Authenticate user with national ID and password
		//	@Tags			Authentication
		//	@Accept			json
		//	@Produce		json
		//	@Param			request	body		dto.SignInRequest	true	"Sign in credentials"
		//	@Success		200		{object}	dto.SignInResponse
		//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
		//	@Failure		401		{object}	dto.ErrorResponse	"Invalid credentials"
		//	@Failure		429		{object}	dto.ErrorResponse	"Too many attempts"
		//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
		//	@Failure		501		{object}	dto.ErrorResponse	"Endpoint not implemented yet"
		//	@Router			/auth/signin [post]
		auth.POST("/signin", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })           // SignIn

		// RefreshToken godoc
		//	@Summary		Refresh access token
		//	@Description	Generate new access token using refresh token
		//	@Tags			Authentication
		//	@Accept			json
		//	@Produce		json
		//	@Param			request	body		dto.RefreshTokenRequest	true	"Refresh token data"
		//	@Success		200		{object}	dto.RefreshTokenResponse
		//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
		//	@Failure		401		{object}	dto.ErrorResponse	"Invalid refresh token"
		//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
		//	@Failure		501		{object}	dto.ErrorResponse	"Endpoint not implemented yet"
		//	@Router			/auth/refresh [post]
		auth.POST("/refresh", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })          // RefreshToken

		// RequestPasswordChange godoc
		//	@Summary		Request password change
		//	@Description	Initiate password change process by sending verification code
		//	@Tags			Authentication
		//	@Accept			json
		//	@Produce		json
		//	@Param			request	body		dto.RequestPasswordChangeRequest	true	"Password change request data"
		//	@Success		200		{object}	dto.RequestPasswordChangeResponse
		//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
		//	@Failure		404		{object}	dto.ErrorResponse	"User not found"
		//	@Failure		429		{object}	dto.ErrorResponse	"Too many requests"
		//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
		//	@Failure		501		{object}	dto.ErrorResponse	"Endpoint not implemented yet"
		//	@Router			/auth/password/request [post]
		auth.POST("/password/request", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // RequestPasswordChange

		// ConfirmPasswordChange godoc
		//	@Summary		Confirm password change
		//	@Description	Confirm password change using verification code
		//	@Tags			Authentication
		//	@Accept			json
		//	@Produce		json
		//	@Param			request	body		dto.ConfirmPasswordChangeRequest	true	"Password change confirmation data"
		//	@Success		200		{object}	dto.ConfirmPasswordChangeResponse
		//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
		//	@Failure		401		{object}	dto.ErrorResponse	"Invalid verification code"
		//	@Failure		404		{object}	dto.ErrorResponse	"User not found"
		//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
		//	@Failure		501		{object}	dto.ErrorResponse	"Endpoint not implemented yet"
		//	@Router			/auth/password/confirm [post]
		auth.POST("/password/confirm", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // ConfirmPasswordChange

		// SignOut godoc (public endpoint for token invalidation)
		//	@Summary		Sign out (public)
		//	@Description	Sign out user from public endpoint (alternative to authenticated /user/signout)
		//	@Tags			Authentication
		//	@Accept			json
		//	@Produce		json
		//	@Param			request	body		dto.SignOutRequest	true	"Sign out data"
		//	@Success		200		{object}	dto.SignOutResponse
		//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
		//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
		//	@Failure		501		{object}	dto.ErrorResponse	"Endpoint not implemented yet"
		//	@Router			/auth/signout [post]
		auth.POST("/signout", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })          // SignOut
	}

	// User routes (authenticated - Owner, Realtor, Agency and Admin)
	user := router.Group("/user")
	{
		// Profile management
		user.GET("/profile", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })    // GetProfile
		user.PUT("/profile", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })    // UpdateProfile
		user.DELETE("/account", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // DeleteAccount
		user.GET("/onboarding", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // GetOnboardingStatus
		user.GET("/roles", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })      // GetUserRoles
		user.GET("/home", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })       // GoHome
		user.PUT("/opt-status", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // UpdateOptStatus

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
	{
		agency.POST("/invite-realtor", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // InviteRealtor
		agency.GET("/realtors", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })        // GetRealtorsByAgency
		agency.GET("/realtors/:id", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })    // GetRealtorByID
		agency.DELETE("/realtors/:id", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // DeleteRealtorByID
	}

	// Realtor routes (Realtor only)
	realtor := router.Group("/realtor")
	{
		realtor.POST("/creci/verify", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })      // VerifyCreciImages
		realtor.POST("/creci/upload-url", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })  // GetCreciUploadURL
		realtor.POST("/invitation/accept", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // AcceptInvitation
		realtor.POST("/invitation/reject", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // RejectInvitation
		realtor.GET("/agency", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })             // GetAgencyOfRealtor
		realtor.DELETE("/agency", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })          // DeleteAgencyOfRealtor
	}
}

// RegisterListingRoutes registers all listing-related routes
func RegisterListingRoutes(router *gin.RouterGroup, handlers *factory.HTTPHandlers) {
	// TODO: Implement listing routes in Etapa 3
	// This is a placeholder for now

	// Main listing routes
	listings := router.Group("/listings")
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

	// Visit management
	visits := router.Group("/visits")
	{
		visits.GET("", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })              // GetAllVisits
		visits.DELETE("/:id", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })       // CancelVisit
		visits.POST("/:id/confirm", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // ConfirmVisitDone
		visits.POST("/:id/approve", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // ApproveVisiting
		visits.POST("/:id/reject", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })  // RejectVisiting
	}

	// Offer management
	offers := router.Group("/offers")
	{
		offers.GET("", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })              // GetAllOffers
		offers.PUT("/:id", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })          // UpdateOffer
		offers.DELETE("/:id", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })       // CancelOffer
		offers.POST("/:id/send", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })    // SendOffer
		offers.POST("/:id/approve", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // ApproveOffer
		offers.POST("/:id/reject", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })  // RejectOffer
	}

	// Evaluation routes
	realtors := router.Group("/realtors")
	{
		realtors.POST("/:id/evaluate", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // EvaluateRealtor
	}

	owners := router.Group("/owners")
	{
		owners.POST("/:id/evaluate", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) }) // EvaluateOwner
	}
}
