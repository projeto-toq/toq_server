package photosessionhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	photosessionservices "github.com/projeto-toq/toq_server/internal/core/service/photo_session_service"
)

// PhotoSessionHandler handles HTTP requests for photographer agenda management.
type PhotoSessionHandler struct {
	service       photosessionservices.PhotoSessionServiceInterface
	globalService globalservice.GlobalServiceInterface
}

// NewPhotoSessionHandler creates a new handler with its dependencies.
func NewPhotoSessionHandler(service photosessionservices.PhotoSessionServiceInterface, globalService globalservice.GlobalServiceInterface) *PhotoSessionHandler {
	return &PhotoSessionHandler{
		service:       service,
		globalService: globalService,
	}
}

// CreateTimeOff handles the creation of a new time-off entry for a photographer.
// @Summary      Create Photographer Time-Off
// @Description  Registers a new time-off period for the authenticated photographer, blocking their agenda.
// @Tags         Photo Session
// @Accept       json
// @Produce      json
// @Param        input body photosessionservices.TimeOffInput true "Time-Off Input"
// @Success      201 {object} object{message=string,timeOffId=int}
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /photographer/agenda/time-off [post]
func (h *PhotoSessionHandler) CreateTimeOff(c *gin.Context) {
	var input photosessionservices.TimeOffInput
	if err := c.ShouldBindJSON(&input); err != nil {
		http_errors.SendHTTPError(c, http.StatusBadRequest, "invalid_json", "Invalid JSON body")
		return
	}

	userID, err := h.globalService.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}
	input.PhotographerID = uint64(userID)

	id, dErr := h.service.CreateTimeOff(c.Request.Context(), input)
	if dErr != nil {
		http_errors.SendHTTPErrorObj(c, dErr)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Time-off created successfully", "timeOffId": id})
}

// DeleteTimeOff handles the deletion of a time-off entry.
// @Summary      Delete Photographer Time-Off
// @Description  Removes an existing time-off period for the authenticated photographer, making slots available again.
// @Tags         Photo Session
// @Accept       json
// @Produce      json
// @Param        input body photosessionservices.DeleteTimeOffInput true "Delete Time-Off Input"
// @Success      200 {object} object{message=string}
// @Failure      400 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /photographer/agenda/time-off [delete]
func (h *PhotoSessionHandler) DeleteTimeOff(c *gin.Context) {
	var input photosessionservices.DeleteTimeOffInput
	if err := c.ShouldBindJSON(&input); err != nil {
		http_errors.SendHTTPError(c, http.StatusBadRequest, "invalid_json", "Invalid JSON body")
		return
	}

	userID, err := h.globalService.GetUserIDFromContext(c.Request.Context())
	if err != nil {
		http_errors.SendHTTPErrorObj(c, err)
		return
	}
	input.PhotographerID = uint64(userID)

	dErr := h.service.DeleteTimeOff(c.Request.Context(), input)
	if dErr != nil {
		http_errors.SendHTTPErrorObj(c, dErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Time-off deleted successfully"})
}
