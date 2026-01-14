package converters

import (
	"database/sql"
	"encoding/base64"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	"github.com/projeto-toq/toq_server/internal/adapter/left/http/middlewares"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	permissionmodel "github.com/projeto-toq/toq_server/internal/core/model/permission_model"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
	proposalservice "github.com/projeto-toq/toq_server/internal/core/service/proposal_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

const (
	maxProposalTextLen = 5000
	defaultPage        = 1
	defaultPageSize    = 20
)

// ProposalActorFromContext extracts the authenticated actor metadata from middleware state.
func ProposalActorFromContext(c *gin.Context) (proposalservice.Actor, error) {
	if c == nil {
		return proposalservice.Actor{}, coreutils.AuthenticationError("")
	}
	userInfo, ok := middlewares.GetUserInfoFromContext(c)
	if !ok || userInfo.ID <= 0 {
		return proposalservice.Actor{}, coreutils.AuthenticationError("")
	}
	return proposalservice.Actor{
		UserID:   userInfo.ID,
		RoleSlug: userInfo.RoleSlug,
	}, nil
}

// CreateProposalDTOToInput validates payload fields and builds the service input for Create.
func CreateProposalDTOToInput(req dto.CreateProposalRequest, actor proposalservice.Actor) (proposalservice.CreateProposalInput, error) {
	if actor.UserID <= 0 {
		return proposalservice.CreateProposalInput{}, coreutils.AuthenticationError("")
	}
	if actor.RoleSlug != permissionmodel.RoleSlugRealtor {
		return proposalservice.CreateProposalInput{}, coreutils.AuthorizationError("Somente corretores podem enviar propostas")
	}
	if req.ListingIdentityID <= 0 {
		return proposalservice.CreateProposalInput{}, coreutils.ValidationError("listingIdentityId", "must be greater than zero")
	}
	text := strings.TrimSpace(req.ProposalText)
	if text == "" {
		return proposalservice.CreateProposalInput{}, coreutils.ValidationError("proposalText", "cannot be empty")
	}
	if utf8.RuneCountInString(text) > maxProposalTextLen {
		return proposalservice.CreateProposalInput{}, coreutils.ValidationError("proposalText", "exceeds maximum length")
	}

	docPayload, err := decodeDocumentUpload(req.Document)
	if err != nil {
		return proposalservice.CreateProposalInput{}, err
	}

	return proposalservice.CreateProposalInput{
		ListingIdentityID: req.ListingIdentityID,
		RealtorID:         actor.UserID,
		ProposalText:      text,
		Document:          docPayload,
	}, nil
}

// UpdateProposalDTOToInput validates payload fields and builds the service input for Update.
func UpdateProposalDTOToInput(req dto.UpdateProposalRequest, actor proposalservice.Actor) (proposalservice.UpdateProposalInput, error) {
	if actor.UserID <= 0 {
		return proposalservice.UpdateProposalInput{}, coreutils.AuthenticationError("")
	}
	if actor.RoleSlug != permissionmodel.RoleSlugRealtor {
		return proposalservice.UpdateProposalInput{}, coreutils.AuthorizationError("Somente corretores podem editar propostas")
	}
	if req.ProposalID <= 0 {
		return proposalservice.UpdateProposalInput{}, coreutils.ValidationError("proposalId", "must be greater than zero")
	}
	text := strings.TrimSpace(req.ProposalText)
	if text == "" {
		return proposalservice.UpdateProposalInput{}, coreutils.ValidationError("proposalText", "cannot be empty")
	}
	if utf8.RuneCountInString(text) > maxProposalTextLen {
		return proposalservice.UpdateProposalInput{}, coreutils.ValidationError("proposalText", "exceeds maximum length")
	}

	docPayload, err := decodeDocumentUpload(req.Document)
	if err != nil {
		return proposalservice.UpdateProposalInput{}, err
	}

	return proposalservice.UpdateProposalInput{
		ProposalID:   req.ProposalID,
		EditorID:     actor.UserID,
		ProposalText: text,
		Document:     docPayload,
	}, nil
}

// ProposalListFilterFromQuery transforms query params into a service list filter.
func ProposalListFilterFromQuery(query dto.ListProposalsQuery, actor proposalservice.Actor) (proposalservice.ListFilter, error) {
	if actor.UserID <= 0 {
		return proposalservice.ListFilter{}, coreutils.AuthenticationError("")
	}
	statuses, err := convertStatuses(query.Statuses)
	if err != nil {
		return proposalservice.ListFilter{}, err
	}
	var listingID *int64
	if query.ListingIdentityID > 0 {
		listingID = ptrInt64(query.ListingIdentityID)
	}

	page := query.Page
	if page <= 0 {
		page = defaultPage
	}
	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	return proposalservice.ListFilter{
		Actor:     actor,
		Statuses:  statuses,
		ListingID: listingID,
		Page:      page,
		PageSize:  pageSize,
	}, nil
}

// ProposalDomainToResponse maps a domain proposal into DTO representation.
func ProposalDomainToResponse(proposal proposalmodel.ProposalInterface) dto.ProposalResponse {
	if proposal == nil {
		return dto.ProposalResponse{}
	}

	response := dto.ProposalResponse{
		ID:                proposal.ID(),
		ListingIdentityID: proposal.ListingIdentityID(),
		Status:            proposal.Status().String(),
		ProposalText:      proposal.ProposalText(),
		DocumentsCount:    proposal.DocumentsCount(),
	}

	if proposal.RejectionReason().Valid {
		reason := proposal.RejectionReason().String
		response.RejectionReason = &reason
	}
	if ptr := timePtrFromNull(proposal.AcceptedAt()); ptr != nil {
		response.AcceptedAt = ptr
	}
	if ptr := timePtrFromNull(proposal.RejectedAt()); ptr != nil {
		response.RejectedAt = ptr
	}
	if ptr := timePtrFromNull(proposal.CancelledAt()); ptr != nil {
		response.CancelledAt = ptr
	}
	if ptr := timePtr(proposal.CreatedAt()); ptr != nil {
		response.CreatedAt = ptr
	}
	if ptr := timePtrFromNull(proposal.FirstOwnerActionAt()); ptr != nil {
		response.ReceivedAt = ptr
	}
	response.RespondedAt = firstNonNilTimePtr(
		timePtrFromNull(proposal.AcceptedAt()),
		timePtrFromNull(proposal.RejectedAt()),
		timePtrFromNull(proposal.CancelledAt()),
	)

	return response
}

// ProposalListToResponse builds a response for list endpoints.
func ProposalListToResponse(result proposalservice.ListResult) dto.ListProposalsResponse {
	items := make([]dto.ProposalResponse, 0, len(result.Items))
	for _, item := range result.Items {
		if item.Proposal == nil {
			continue
		}
		response := ProposalDomainToResponse(item.Proposal)
		response.Documents = proposalDocumentsToResponse(item.Documents)
		response.Realtor = proposalRealtorToResponse(item.Realtor)
		response.Owner = proposalOwnerToResponse(item.Owner)
		response.Listing = listingSummaryFromDomain(item.Listing)
		items = append(items, response)
	}

	return dto.ListProposalsResponse{
		Items: items,
		Total: result.Total,
	}
}

// ProposalDetailToResponse aggregates proposal info and documents.
func ProposalDetailToResponse(detail proposalservice.DetailResult) dto.ProposalDetailResponse {
	proposalDTO := ProposalDomainToResponse(detail.Proposal)
	proposalDTO.Realtor = proposalRealtorToResponse(detail.Realtor)
	proposalDTO.Owner = proposalOwnerToResponse(detail.Owner)
	proposalDTO.Documents = proposalDocumentsToResponse(detail.Documents)
	proposalDTO.Listing = listingSummaryFromDomain(detail.Listing)

	documents := proposalDocumentsToResponse(detail.Documents)

	return dto.ProposalDetailResponse{
		Proposal:  proposalDTO,
		Documents: documents,
		Realtor:   proposalDTO.Realtor,
		Owner:     proposalDTO.Owner,
	}
}

func listingSummaryFromDomain(listing listingmodel.ListingInterface) *dto.ListingSummaryDTO {
	if listing == nil {
		return nil
	}

	summary := dto.ListingSummaryDTO{
		ListingIdentityID: listing.IdentityID(),
		Title:             strings.TrimSpace(listing.Title()),
		Description:       strings.TrimSpace(listing.Description()),
		ZipCode:           listing.ZipCode(),
		Street:            listing.Street(),
		Number:            listing.Number(),
		Neighborhood:      listing.Neighborhood(),
		City:              listing.City(),
		State:             listing.State(),
		PropertyType:      BuildListingPropertyTypeDTO(listing.ListingType()),
	}

	if complement := strings.TrimSpace(listing.Complement()); complement != "" {
		summary.Complement = complement
	}

	return &summary
}

func decodeDocumentUpload(doc *dto.ProposalDocumentUpload) (*proposalservice.DocumentPayload, error) {
	if doc == nil {
		return nil, nil
	}
	fileName := strings.TrimSpace(doc.FileName)
	if fileName == "" {
		return nil, coreutils.ValidationError("document.fileName", "cannot be empty")
	}
	mime := strings.TrimSpace(doc.MimeType)
	if mime == "" {
		return nil, coreutils.ValidationError("document.mimeType", "cannot be empty")
	}
	payload := strings.TrimSpace(doc.Base64Payload)
	if payload == "" {
		return nil, coreutils.ValidationError("document.base64Payload", "cannot be empty")
	}
	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, coreutils.ValidationError("document.base64Payload", "must be valid base64")
	}
	return &proposalservice.DocumentPayload{
		FileName:  fileName,
		MimeType:  mime,
		Bytes:     decoded,
		SizeBytes: int64(len(decoded)),
	}, nil
}

func convertStatuses(values []string) ([]proposalmodel.Status, error) {
	if len(values) == 0 {
		return nil, nil
	}
	statuses := make([]proposalmodel.Status, 0, len(values))
	for _, value := range values {
		status := proposalmodel.Status(strings.TrimSpace(strings.ToLower(value)))
		switch status {
		case proposalmodel.StatusPending,
			proposalmodel.StatusAccepted,
			proposalmodel.StatusRefused,
			proposalmodel.StatusCancelled:
			statuses = append(statuses, status)
		default:
			return nil, coreutils.ValidationError("status", "unsupported value")
		}
	}
	return statuses, nil
}

func proposalDocumentsToResponse(documents []proposalmodel.ProposalDocumentInterface) []dto.ProposalDocumentResponse {
	if len(documents) == 0 {
		return nil
	}
	result := make([]dto.ProposalDocumentResponse, 0, len(documents))
	for _, doc := range documents {
		if doc == nil {
			continue
		}
		payload := ""
		if data := doc.FileData(); len(data) > 0 {
			payload = base64.StdEncoding.EncodeToString(data)
		}
		result = append(result, dto.ProposalDocumentResponse{
			ID:            doc.ID(),
			FileName:      doc.FileName(),
			MimeType:      doc.MimeType(),
			FileSizeBytes: doc.FileSizeBytes(),
			Base64Payload: payload,
		})
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func proposalRealtorToResponse(summary proposalmodel.RealtorSummary) dto.ProposalRealtorResponse {
	if summary == nil {
		return dto.ProposalRealtorResponse{}
	}
	nickname := summary.Nickname()
	value := ""
	if nickname.Valid {
		value = nickname.String
	}
	return dto.ProposalRealtorResponse{
		Name:              summary.Name(),
		Nickname:          value,
		AccountAgeMonths:  summary.UsageMonths(),
		ProposalsCreated:  summary.ProposalsCreated(),
		AcceptedProposals: summary.AcceptedProposals(),
		PhotoURL:          summary.PhotoURL(),
	}
}

func proposalOwnerToResponse(summary proposalmodel.OwnerSummary) dto.ProposalOwnerResponse {
	if summary == nil {
		return dto.ProposalOwnerResponse{}
	}

	var proposalAvg *int64
	if summary.ProposalAvgSeconds().Valid {
		v := summary.ProposalAvgSeconds().Int64
		proposalAvg = &v
	}

	var visitAvg *int64
	if summary.VisitAvgSeconds().Valid {
		v := summary.VisitAvgSeconds().Int64
		visitAvg = &v
	}

	return dto.ProposalOwnerResponse{
		ID:                 summary.ID(),
		FullName:           summary.FullName(),
		MemberSinceMonths:  summary.MemberSinceMonths(),
		PhotoURL:           summary.PhotoURL(),
		ProposalAvgSeconds: proposalAvg,
		VisitAvgSeconds:    visitAvg,
	}
}

func timePtrFromNull(v sql.NullTime) *time.Time {
	if !v.Valid {
		return nil
	}
	ts := v.Time
	return &ts
}

func timePtr(ts time.Time) *time.Time {
	if ts.IsZero() {
		return nil
	}
	value := ts
	return &value
}

func firstNonNilTimePtr(values ...*time.Time) *time.Time {
	for _, v := range values {
		if v != nil {
			return v
		}
	}
	return nil
}

func ptrInt64(value int64) *int64 {
	return &value
}
