package userhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
)

// VerifyCreciDocuments confirms that the three CRECI documents exist in storage and moves the user to PendingManual.

// Ensure Swag can resolve dto.ErrorResponse referenced in annotations without affecting runtime
type _ = dto.ErrorResponse

//		@Summary      Verify CRECI documents presence
//		@Description  Check if selfie.jpg, front.jpg and back.jpg exist in the user's S3 folder; if all exist, set status to PendingManual
//		@Tags         Realtor
//		@Accept       json
//		@Produce      json
//		@Success      200  {object}  map[string]string  "verification accepted"
//	 @Failure      401  {object}  dto.ErrorResponse  "Unauthorized"
//	 @Failure      403  {object}  dto.ErrorResponse  "Forbidden"
//	 @Failure      422  {object}  dto.ErrorResponse  "Missing required documents"
//	 @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
//		@Router       /realtor/creci/verify [post]
//		@Security     BearerAuth
func (uh *UserHandler) VerifyCreciDocuments(c *gin.Context) {
	ctx := c.Request.Context()
	if err := uh.userService.VerifyCreciDocuments(ctx); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "verification accepted"})
}
