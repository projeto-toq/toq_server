package routes

import (
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes registers all user-related routes
func RegisterUserRoutes(router *gin.RouterGroup) {
	// TODO: Implement user routes in Etapa 2
	// This is a placeholder for now
	
	// Authentication routes (public)
	auth := router.Group("/auth")
	{
		auth.POST("/owner", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		auth.POST("/realtor", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		auth.POST("/agency", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		auth.POST("/signin", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		auth.POST("/refresh", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		auth.POST("/password/request", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		auth.POST("/password/confirm", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		auth.POST("/signout", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
	}

	// User routes (authenticated)
	user := router.Group("/user")
	{
		user.GET("/profile", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		user.PUT("/profile", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		user.DELETE("/account", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		user.GET("/onboarding", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		user.GET("/roles", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		user.GET("/home", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		user.PUT("/opt-status", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		user.POST("/photo/upload-url", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		user.GET("/profile/thumbnails", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
	}

	// Agency routes
	agency := router.Group("/agency")
	{
		agency.POST("/invite-realtor", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		agency.GET("/realtors", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		agency.GET("/realtors/:id", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		agency.DELETE("/realtors/:id", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
	}

	// Realtor routes
	realtor := router.Group("/realtor")
	{
		realtor.POST("/creci/verify", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		realtor.POST("/creci/upload-url", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		realtor.POST("/invitation/accept", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		realtor.POST("/invitation/reject", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		realtor.GET("/agency", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		realtor.DELETE("/agency", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
	}
}

// RegisterListingRoutes registers all listing-related routes
func RegisterListingRoutes(router *gin.RouterGroup) {
	// TODO: Implement listing routes in Etapa 3
	// This is a placeholder for now
	
	listings := router.Group("/listings")
	{
		listings.GET("", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		listings.POST("", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		listings.GET("/:id", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		listings.PUT("/:id", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
		listings.DELETE("/:id", func(c *gin.Context) { c.JSON(501, gin.H{"error": "Not implemented yet"}) })
	}
}
