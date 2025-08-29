package userhandlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/validators"
)

func (uh *UserHandler) RequestEmailChange(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user information from context (set by middleware)
	userInfos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		utils.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfo := userInfos.(usermodel.UserInfos)

	// Parse request body
	var request dto.RequestEmailChangeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Validate and clean email
	email := strings.TrimSpace(request.NewEmail)
	if err := validators.ValidateEmail(email); err != nil {
		utils.SendHTTPError(c, http.StatusBadRequest, "INVALID_EMAIL", "Invalid email format")
		return
	}

	// Call service to request email change
	if err := uh.userService.RequestEmailChange(ctx, userInfo.ID, email); err != nil {
		utils.SendHTTPError(c, http.StatusInternalServerError, "REQUEST_EMAIL_CHANGE_FAILED", "Failed to request email change")
		return
	}

	// Prepare response
	response := dto.RequestEmailChangeResponse{
		Message: "Email change request sent successfully",
	}

	c.JSON(http.StatusOK, response)
}
