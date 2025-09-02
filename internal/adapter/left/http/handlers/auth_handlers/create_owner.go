package authhandlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/giulio-alfieri/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/http_errors"
	httputils "github.com/giulio-alfieri/toq_server/internal/adapter/left/http/utils"
	permissionmodel "github.com/giulio-alfieri/toq_server/internal/core/model/permission_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

// CreateOwner handles owner account creation (public endpoint)
//
//	@Summary		Create owner account
//	@Description	Create a new owner account with user information
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateOwnerRequest	true	"Owner creation data (include optional deviceToken for push notifications)"
//	@Success		201		{object}	dto.CreateOwnerResponse
//	@Failure		400		{object}	dto.ErrorResponse	"Invalid request format"
//	@Failure		409		{object}	dto.ErrorResponse	"User already exists"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal server error"
//	@Failure		422		{object}	dto.ErrorResponse	"Validation failed"
//	@Router			/auth/owner [post]
func (ah *AuthHandler) CreateOwner(c *gin.Context) {
	ctx, spanEnd, err := utils.GenerateTracer(c.Request.Context())
	if err != nil {
		httperrors.SendHTTPErrorObj(c, utils.NewHTTPError(http.StatusInternalServerError, "Failed to generate tracer"))
		return
	}
	defer spanEnd()

	// Parse request
	var request dto.CreateOwnerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	// Create user model from DTO
	user, err := ah.createUserFromDTO(request.Owner, permissionmodel.RoleSlugOwner)
	if err != nil {
		// If parsing failed (e.g., date), return 422 with field-level detail
		// Wrap generic errors as validation failure for client clarity
		httperrors.SendHTTPErrorObj(c, utils.NewHTTPError(http.StatusUnprocessableEntity, "Validation failed", map[string]any{
			"errors": []map[string]string{{"field": "owner.bornAt", "message": "invalid date, expected YYYY-MM-DD"}},
		}))
		return
	}

	// Extract request context for security logging and session metadata
	reqContext := utils.ExtractRequestContext(c)

	// Call service: cria a conta e autentica via SignIn padrão
	tokens, err := ah.userService.CreateOwner(ctx, user, request.Owner.Password, request.DeviceToken, reqContext.IPAddress, reqContext.UserAgent)
	if err != nil {
		if derr, ok := err.(utils.DomainError); ok {
			httperrors.SendHTTPErrorObj(c, derr)
			return
		}
		// Fallback: conflito genérico
		httperrors.SendHTTPErrorObj(c, utils.NewHTTPError(http.StatusConflict, "Failed to create owner"))
		return
	}

	// Success response
	c.JSON(http.StatusCreated, dto.CreateOwnerResponse{
		Tokens: dto.TokensResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	})
}

// createUserFromDTO converts DTO to User model
func (ah *AuthHandler) createUserFromDTO(dtoUser dto.UserCreateRequest, role permissionmodel.RoleSlug) (usermodel.UserInterface, error) {
	user := usermodel.NewUser()

	// Parse dates
	bornAt, err := time.Parse("2006-01-02", dtoUser.BornAt)
	if err != nil {
		return nil, err
	}

	var creciValidity time.Time
	if dtoUser.CreciValidity != "" {
		creciValidity, err = time.Parse("2006-01-02", dtoUser.CreciValidity)
		if err != nil {
			return nil, err
		}
	}

	// Set user data
	user.SetFullName(dtoUser.FullName)
	user.SetNickName(dtoUser.NickName)
	user.SetNationalID(dtoUser.NationalID)
	user.SetCreciNumber(dtoUser.CreciNumber)
	user.SetCreciState(dtoUser.CreciState)
	if !creciValidity.IsZero() {
		user.SetCreciValidity(creciValidity)
	}
	user.SetBornAt(bornAt)
	user.SetPhoneNumber(dtoUser.PhoneNumber)
	user.SetEmail(dtoUser.Email)
	user.SetZipCode(dtoUser.ZipCode)
	user.SetStreet(dtoUser.Street)
	user.SetNumber(dtoUser.Number)
	user.SetComplement(dtoUser.Complement)
	user.SetNeighborhood(dtoUser.Neighborhood)
	user.SetCity(dtoUser.City)
	user.SetState(dtoUser.State)
	user.SetPassword(dtoUser.Password)

	// Note: Active role will be set by the service layer, not the handler

	user.SetOptStatus(false) // Default opt-in status
	user.SetDeleted(false)

	return user, nil
}
