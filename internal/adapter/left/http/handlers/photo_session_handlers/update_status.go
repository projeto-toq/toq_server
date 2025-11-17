package photosessionhandlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/core/derrors"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
)

// UpdateSessionStatus handles the request to update photo session status.
//
//	@Summary   Update photo session status (approve/reject/complete)
//	@Description Updates a photo session booking status with support for approval/rejection and completion workflows.
//	             **Approval/Rejection Workflow** (only when manual approval enabled):
//	             - ACCEPTED: Photographer accepts the session (requires require_photographer_approval=true)
//	             - REJECTED: Photographer declines the session (requires require_photographer_approval=true)
//	             **Completion Workflow** (always available):
//	             - DONE: Photographer marks session as completed after performing it
//	             When automatic approval is enabled (require_photographer_approval=false), ACCEPTED/REJECTED
//	             transitions are blocked (sessions are auto-approved at reservation time), but DONE is always available.
//	@Tags      Photographer
//	@Accept    json
//	@Produce   json
//	@Param     request   body      dto.UpdateSessionStatusRequest true "Status update request (ACCEPTED/REJECTED/DONE)"
//	@Success   204       {object}  nil                "Status successfully updated"
//	@Failure   400       {object}  dto.ErrorResponse  "Invalid payload or approval disabled for ACCEPTED/REJECTED transitions"
//	@Failure   401       {object}  dto.ErrorResponse  "Unauthorized"
//	@Failure   403       {object}  dto.ErrorResponse  "Forbidden - session does not belong to photographer"
//	@Failure   404       {object}  dto.ErrorResponse  "Session not found"
//	@Failure   409       {object}  dto.ErrorResponse  "Session not in valid state for requested transition"
//	@Failure   422       {object}  dto.ErrorResponse  "Invalid status value (must be ACCEPTED, REJECTED, or DONE)"
//	@Failure   500       {object}  dto.ErrorResponse  "Internal error"
//	@Router    /photographer/sessions/status [post]
//	@Security  BearerAuth
func (h *PhotoSessionHandler) UpdateSessionStatus(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := h.globalService.GetUserIDFromContext(ctx)
	if err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}

	var req dto.UpdateSessionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		http_errors.SendHTTPErrorObj(c, http_errors.ConvertBindError(err))
		return
	}

	if req.SessionID == 0 {
		http_errors.SendHTTPErrorObj(c, derrors.Validation("sessionId is required", map[string]string{"sessionId": "required"}))
		return
	}

	status := strings.TrimSpace(req.Status)
	if status == "" {
		http_errors.SendHTTPErrorObj(c, derrors.Validation("status is required", map[string]string{"status": "required"}))
		return
	}

	input := photosessionservices.UpdateSessionStatusInput{
		SessionID:      req.SessionID,
		PhotographerID: uint64(userID),
		Status:         status,
	}

	if err := h.service.UpdateSessionStatus(ctx, input); err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
