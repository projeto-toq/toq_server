package userhandlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// UpdateProfile handles updating user profile
func (uh *UserHandler) UpdateProfile(c *gin.Context) {
	ctx, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "TRACER_ERROR", "Failed to generate tracer")
		return
	}
	defer spanEnd()

	// Get user info from context (set by auth middleware)
	infos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfos := infos.(usermodel.UserInfos)

	// Parse request body using DTO
	var request dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Get current user data first
	currentUser, err := uh.userService.GetProfile(ctx, userInfos.ID)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Update fields from request
	if request.User.NickName != "" {
		currentUser.SetNickName(request.User.NickName)
	}
	if request.User.BornAt != "" {
		bornAt, err := time.Parse("2006-01-02", request.User.BornAt)
		if err != nil {
			httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_DATE", "Invalid born_at date format")
			return
		}
		currentUser.SetBornAt(bornAt)
	}
	if request.User.ZipCode != "" {
		currentUser.SetZipCode(request.User.ZipCode)
	}
	if request.User.Street != "" {
		currentUser.SetStreet(request.User.Street)
	}
	if request.User.Number != "" {
		currentUser.SetNumber(request.User.Number)
	}
	if request.User.Complement != "" {
		currentUser.SetComplement(request.User.Complement)
	}
	if request.User.Neighborhood != "" {
		currentUser.SetNeighborhood(request.User.Neighborhood)
	}
	if request.User.City != "" {
		currentUser.SetCity(request.User.City)
	}
	if request.User.State != "" {
		currentUser.SetState(request.User.State)
	}

	// Call service to update profile
	if err := uh.userService.UpdateProfile(ctx, currentUser); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Success response using DTO
	response := dto.UpdateProfileResponse{
		Message: "Profile updated successfully",
	}
	c.JSON(http.StatusOK, response)
}
