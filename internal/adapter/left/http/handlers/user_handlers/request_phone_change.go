package userhandlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils/validators"
)

func (uh *UserHandler) RequestPhoneChange(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user information from context (set by middleware)
	userInfos, exists := c.Get(string(globalmodel.TokenKey))
	if !exists {
		httperrors.SendHTTPError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	userInfo := userInfos.(usermodel.UserInfos)

	// Parse request body
	var request dto.RequestPhoneChangeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	// Validate and clean phone number
	newPhone := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(request.NewPhoneNumber, " ", ""), "-", ""), "(", ""), ")", ""), ".", "")
	if err := validators.ValidateE164(newPhone); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_PHONE", "Invalid phone number format")
		return
	}

	// Call service to request phone change
	if err := uh.userService.RequestPhoneChange(ctx, userInfo.ID, newPhone); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	// Prepare response
	response := dto.RequestPhoneChangeResponse{
		Message: "Phone change request sent successfully",
	}

	c.JSON(http.StatusOK, response)
}
