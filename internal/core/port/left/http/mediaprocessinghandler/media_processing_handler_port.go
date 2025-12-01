package mediaprocessinghandler

import "github.com/gin-gonic/gin"

// MediaProcessingHandlerPort define o contrato para operações HTTP de mídia.
type MediaProcessingHandlerPort interface {
	// Upload & Processing Flow
	RequestUploadURLs(c *gin.Context)
	ProcessMedia(c *gin.Context)

	// Retrieval
	ListMedia(c *gin.Context)
	GenerateDownloadURLs(c *gin.Context)

	// Management
	UpdateMedia(c *gin.Context)
	DeleteMedia(c *gin.Context)
	CompleteMedia(c *gin.Context) // Finalização manual/zip

	// Callbacks (Internal/Webhook)
	HandleProcessingCallback(c *gin.Context)
}
