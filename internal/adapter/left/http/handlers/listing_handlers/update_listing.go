package listinghandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	httperrors "github.com/projeto-toq/toq_server/internal/adapter/left/http/http_errors"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateListing atualiza as informações de um anúncio existente.
//
//	@Summary	Atualiza um anúncio
//	@Description	Permite atualização parcial dos campos de um anúncio em rascunho. Campos omitidos permanecem inalterados; campos presentes (inclusive null/vazio) sobrescrevem o valor atual.
//	@Tags		Listings
//	@Accept		json
//	@Produce	json
//	@Param		request	body	dto.UpdateListingRequest	true	"Dados para atualização (ID obrigatório no corpo)"
//	@Success	200	{object}	dto.UpdateListingResponse
//	@Failure	400	{object}	dto.ErrorResponse	"Payload inválido"
//	@Failure	401	{object}	dto.ErrorResponse	"Não autorizado"
//	@Failure	403	{object}	dto.ErrorResponse	"Proibido"
//	@Failure	404	{object}	dto.ErrorResponse	"Não encontrado"
//	@Failure	409	{object}	dto.ErrorResponse	"Conflito"
//	@Failure	500	{object}	dto.ErrorResponse	"Erro interno"
//	@Router		/listings [put]
//	@Security	BearerAuth
func (lh *ListingHandler) UpdateListing(c *gin.Context) {
	baseCtx := coreutils.EnrichContextWithRequestInfo(c.Request.Context(), c)

	if _, ok := middlewares.GetUserInfoFromContext(c); !ok {
		httperrors.SendHTTPError(c, http.StatusInternalServerError, "INTERNAL_CONTEXT_MISSING", "User context not found")
		return
	}

	var request dto.UpdateListingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request format")
		return
	}

	if !request.ID.IsPresent() || request.ID.IsNull() {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "MISSING_ID", "Listing ID must be provided in the request body")
		return
	}

	listingID, ok := request.ID.Value()
	if !ok || listingID <= 0 {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_ID", "Listing ID is invalid")
		return
	}

	input, err := converters.UpdateListingRequestToInput(request)
	if err != nil {
		httperrors.SendHTTPError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	input.ID = listingID

	if err := lh.listingService.UpdateListing(baseCtx, input); err != nil {
		httperrors.SendHTTPErrorObj(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.UpdateListingResponse{
		Success: true,
		Message: "Listing updated",
	})
}
