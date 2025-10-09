package authhandlers

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	httputils "github.com/projeto-toq/toq_server/internal/adapter/left/http/utils"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
	"github.com/projeto-toq/toq_server/internal/core/utils/hmacauth"
	validators "github.com/projeto-toq/toq_server/internal/core/utils/validators"
)

const (
	dateLayoutISO = "2006-01-02"
)

// ValidateCPF valida um CPF utilizando o serviço externo disponível via userService
//
// @Summary     Validate CPF
// @Description Validates a CPF using Receita Federal integration. Requires signed payload.
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       request body dto.ValidateCPFRequest true "Signed CPF validation payload"
// @Success     200 {object} dto.ValidationResultResponse
// @Failure     400 {object} dto.ErrorResponse "Validation error"
// @Failure     401 {object} dto.ErrorResponse "Invalid signature or expired request"
// @Failure     422 {object} dto.ErrorResponse "Semantic validation error"
// @Failure     429 {object} dto.ErrorResponse "Rate limited"
// @Failure     500 {object} dto.ErrorResponse "Internal server error"
// @Router      /auth/validate/cpf [post]
func (ah *AuthHandler) ValidateCPF(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	logger := coreutils.LoggerFromContext(ctx)

	rawBody, err := ah.readRequestBody(c)
	if err != nil {
		logger.Error("auth.validate_cpf.read_body_error", "err", err)
		httperrors.SendHTTPErrorObj(c, coreutils.InternalError("Failed to process request body."))
		return
	}

	var request dto.ValidateCPFRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	if !ah.validateRequestSignature(c, ctx, rawBody, request.Timestamp, request.HMAC, request.NationalID, "cpf") {
		return
	}

	bornAt, err := time.Parse(dateLayoutISO, request.BornAt)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, coreutils.ValidationError("bornAt", "Invalid date format. Expected YYYY-MM-DD."))
		return
	}

	if err := ah.userService.ValidateCPF(ctx, request.NationalID, bornAt); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ValidationResultResponse{Valid: true})
}

// ValidateCNPJ valida um CNPJ utilizando o serviço externo já integrado
//
// @Summary     Validate CNPJ
// @Description Validates a CNPJ using Receita Federal integration. Requires signed payload.
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       request body dto.ValidateCNPJRequest true "Signed CNPJ validation payload"
// @Success     200 {object} dto.ValidationResultResponse
// @Failure     400 {object} dto.ErrorResponse "Validation error"
// @Failure     401 {object} dto.ErrorResponse "Invalid signature or expired request"
// @Failure     422 {object} dto.ErrorResponse "Semantic validation error"
// @Failure     429 {object} dto.ErrorResponse "Rate limited"
// @Failure     500 {object} dto.ErrorResponse "Internal server error"
// @Router      /auth/validate/cnpj [post]
func (ah *AuthHandler) ValidateCNPJ(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	logger := coreutils.LoggerFromContext(ctx)

	rawBody, err := ah.readRequestBody(c)
	if err != nil {
		logger.Error("auth.validate_cnpj.read_body_error", "err", err)
		httperrors.SendHTTPErrorObj(c, coreutils.InternalError("Failed to process request body."))
		return
	}

	var request dto.ValidateCNPJRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	if !ah.validateRequestSignature(c, ctx, rawBody, request.Timestamp, request.HMAC, request.NationalID, "cnpj") {
		return
	}

	if err := ah.userService.ValidateCNPJ(ctx, request.NationalID); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ValidationResultResponse{Valid: true})
}

// ValidateCEP consulta um CEP e retorna os dados completos do endereço
//
// @Summary     Validate CEP
// @Description Retrieves CEP information using the configured provider. Requires signed payload.
// @Tags        Authentication
// @Accept      json
// @Produce     json
// @Param       request body dto.ValidateCEPRequest true "Signed CEP validation payload"
// @Success     200 {object} dto.CEPValidationResponse
// @Failure     400 {object} dto.ErrorResponse "Validation error"
// @Failure     401 {object} dto.ErrorResponse "Invalid signature or expired request"
// @Failure     422 {object} dto.ErrorResponse "Semantic validation error"
// @Failure     429 {object} dto.ErrorResponse "Rate limited"
// @Failure     500 {object} dto.ErrorResponse "Internal server error"
// @Router      /auth/validate/cep [post]
func (ah *AuthHandler) ValidateCEP(c *gin.Context) {
	ctx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)
	logger := coreutils.LoggerFromContext(ctx)

	rawBody, err := ah.readRequestBody(c)
	if err != nil {
		logger.Error("auth.validate_cep.read_body_error", "err", err)
		httperrors.SendHTTPErrorObj(c, coreutils.InternalError("Failed to process request body."))
		return
	}

	var request dto.ValidateCEPRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPErrorObj(c, httputils.MapBindingError(err))
		return
	}

	if !ah.validateRequestSignature(c, ctx, rawBody, request.Timestamp, request.HMAC, request.PostalCode, "cep") {
		return
	}

	cepInfo, err := ah.globalService.GetCEP(ctx, request.PostalCode)
	if err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	response := dto.CEPValidationResponse{
		Valid:        true,
		PostalCode:   cepInfo.GetCep(),
		Street:       cepInfo.GetStreet(),
		Complement:   cepInfo.GetComplement(),
		Neighborhood: cepInfo.GetNeighborhood(),
		City:         cepInfo.GetCity(),
		State:        cepInfo.GetState(),
	}

	c.JSON(http.StatusOK, response)
}

func (ah *AuthHandler) readRequestBody(c *gin.Context) ([]byte, error) {
	body, err := c.GetRawData()
	if err != nil {
		return nil, err
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	return body, nil
}

func (ah *AuthHandler) validateRequestSignature(c *gin.Context, ctx context.Context, rawBody []byte, timestamp int64, signature string, identifier string, validationType string) bool {
	if ah.hmacValidator == nil {
		logger := coreutils.LoggerFromContext(ctx)
		logger.Error("auth.validate.signature.validator_missing")
		httperrors.SendHTTPErrorObj(c, coreutils.InternalError("Signature validator not configured."))
		return false
	}

	path := c.FullPath()
	if path == "" {
		path = c.Request.URL.Path
	}

	err := ah.hmacValidator.ValidateSignature(c.Request.Method, path, timestamp, rawBody, signature)
	if err == nil {
		return true
	}

	reqCtx := coreutils.ExtractRequestContext(c)
	logger := coreutils.LoggerFromContext(ctx)
	logger.Warn("auth.validate.signature_failed",
		"type", validationType,
		"path", path,
		"ip", reqCtx.IPAddress,
		"user_agent", reqCtx.UserAgent,
		"identifier", maskIdentifier(identifier),
		"err", err,
	)

	httperrors.SendHTTPErrorObj(c, mapSignatureError(err))
	return false
}

func mapSignatureError(err error) error {
	switch {
	case errors.Is(err, hmacauth.ErrTimestampMissing), errors.Is(err, hmacauth.ErrTimestampInvalid):
		return coreutils.ValidationError("timestamp", "Invalid or missing request timestamp.")
	case errors.Is(err, hmacauth.ErrTimestampOutsideSkew):
		return coreutils.NewHTTPError(http.StatusUnauthorized, "Request timestamp expired.")
	case errors.Is(err, hmacauth.ErrSignatureRequired):
		return coreutils.ValidationError("hmac", "Request signature is required.")
	case errors.Is(err, hmacauth.ErrSignatureInvalid):
		return coreutils.ValidationError("hmac", "Invalid request signature format.")
	case errors.Is(err, hmacauth.ErrSignatureMismatch):
		return coreutils.NewHTTPError(http.StatusUnauthorized, "Invalid request signature.")
	default:
		return coreutils.InternalError("Failed to validate request signature.")
	}
}

func maskIdentifier(value string) string {
	digits := validators.OnlyDigits(value)
	length := len(digits)
	if length == 0 {
		return "***"
	}

	if length <= 4 {
		return strings.Repeat("*", length)
	}

	prefix := 3
	suffix := 2
	if length < 6 {
		prefix = 1
		suffix = 1
	}

	masked := strings.Repeat("*", length-prefix-suffix)
	return digits[:prefix] + masked + digits[length-suffix:]
}
