package photosessionhandlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/core/derrors"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
)

// UpdateSessionStatus handles the request to accept or reject a photo session.
//
//	@Summary   Accept or reject a photo session
//	@Tags      Photographer
//	@Accept    json
//	@Produce   json
//	@Param     sessionId path      int                true "Session ID"
//	@Param     request   body      dto.UpdateSessionStatusRequest true "Status update request"
//	@Success   204       {object}  nil                "No Content"
//	@Failure   400       {object}  dto.ErrorResponse  "Invalid payload"
//	@Failure   401       {object}  dto.ErrorResponse  "Unauthorized"
//	@Failure   403       {object}  dto.ErrorResponse  "Forbidden"
//	@Failure   404       {object}  dto.ErrorResponse  "Session not found"
//	@Failure   409       {object}  dto.ErrorResponse  "Session not in pending state"
//	@Failure   422       {object}  dto.ErrorResponse  "Invalid status value"
//	@Failure   500       {object}  dto.ErrorResponse  "Internal error"
//	@Router    /photographer/sessions/{sessionId}/status [post]
//	@Security  BearerAuth
func (h *PhotoSessionHandler) UpdateSessionStatus(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := h.globalService.GetUserIDFromContext(ctx)
	if err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}

	sessionIDStr := c.Param("sessionId")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 64)
	if err != nil {
		http_errors.SendHTTPErrorObj(c, derrors.BadRequest("invalid sessionId"))
		return
	}

	var req dto.UpdateSessionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		http_errors.SendHTTPErrorObj(c, http_errors.ConvertBindError(err))
		return
	}

	status := strings.TrimSpace(req.Status)
	if status == "" {
		http_errors.SendHTTPErrorObj(c, derrors.Validation("status is required", map[string]string{"status": "required"}))
		return
	}

	input := photosessionservices.UpdateSessionStatusInput{
		SessionID:      sessionID,
		PhotographerID: uint64(userID),
		Status:         status,
	}

	if err := h.service.UpdateSessionStatus(ctx, input); err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
