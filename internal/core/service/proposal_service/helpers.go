package proposalservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/projeto-toq/toq_server/internal/core/derrors"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
	ownermetricsrepository "github.com/projeto-toq/toq_server/internal/core/port/right/repository/owner_metrics_repository"
	globalservice "github.com/projeto-toq/toq_server/internal/core/service/global_service"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

const (
	pdfMimeType         = "application/pdf"
	realtorPhotoVariant = "small"
)

func (s *proposalService) validateCreateInput(input CreateProposalInput) error {
	if input.ListingIdentityID <= 0 {
		return derrors.Validation("listingIdentityId must be greater than zero", map[string]any{"listingIdentityId": "required"})
	}
	if input.RealtorID <= 0 {
		return derrors.Validation("realtorId must be greater than zero", map[string]any{"realtorId": "required"})
	}
	if strings.TrimSpace(input.ProposalText) == "" {
		return derrors.Validation("proposal text is required", map[string]any{"proposalText": "required"})
	}
	return nil
}

func (s *proposalService) validateUpdateInput(input UpdateProposalInput) error {
	if input.ProposalID <= 0 {
		return derrors.Validation("proposalId must be greater than zero", map[string]any{"proposalId": "required"})
	}
	if input.EditorID <= 0 {
		return derrors.Validation("editorId must be greater than zero", map[string]any{"editorId": "required"})
	}
	if strings.TrimSpace(input.ProposalText) == "" {
		return derrors.Validation("proposal text is required", map[string]any{"proposalText": "required"})
	}
	return nil
}

func (s *proposalService) validateDocumentPayload(payload *DocumentPayload) error {
	if payload == nil {
		return nil
	}

	fileName := strings.TrimSpace(payload.FileName)
	if fileName == "" {
		return derrors.Validation("fileName is required", map[string]any{"document.fileName": "required"})
	}
	if len(payload.Bytes) == 0 {
		return derrors.Validation("document bytes are required", map[string]any{"document": "empty"})
	}

	mime := strings.ToLower(strings.TrimSpace(payload.MimeType))
	if mime == "" {
		mime = pdfMimeType
	}
	if mime != pdfMimeType {
		return derrors.Validation("document must be a PDF", map[string]any{"document.mimeType": mime})
	}

	size := payload.SizeBytes
	if size <= 0 {
		size = int64(len(payload.Bytes))
	}
	if size <= 0 {
		return derrors.Validation("document size is invalid", map[string]any{"document": "size"})
	}
	if size > s.maxDocBytes {
		return derrors.Validation(
			"document exceeds maximum size",
			map[string]any{"maxBytes": s.maxDocBytes},
		)
	}

	return nil
}

func (s *proposalService) createProposalDocument(ctx context.Context, tx *sql.Tx, proposalID int64, payload *DocumentPayload) error {
	if payload == nil {
		return nil
	}

	doc := proposalmodel.NewProposalDocument()
	doc.SetProposalID(proposalID)
	doc.SetFileName(strings.TrimSpace(payload.FileName))
	doc.SetMimeType(strings.TrimSpace(payload.MimeType))
	if doc.MimeType() == "" {
		doc.SetMimeType(pdfMimeType)
	}
	size := payload.SizeBytes
	if size <= 0 {
		size = int64(len(payload.Bytes))
	}
	doc.SetFileSizeBytes(size)
	doc.SetFileData(payload.Bytes)

	if err := s.proposalRepo.CreateDocument(ctx, tx, doc); err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("proposal.document.create_error", "err", err, "proposal_id", proposalID)
		return derrors.Infra("failed to store proposal document", err)
	}

	return nil
}

func (s *proposalService) rollbackOnError(ctx context.Context, tx *sql.Tx, opErr *error) {
	if tx == nil || opErr == nil || *opErr == nil {
		return
	}
	if rbErr := s.globalSvc.RollbackTransaction(ctx, tx); rbErr != nil {
		utils.SetSpanError(ctx, rbErr)
		utils.LoggerFromContext(ctx).Error("proposal.tx.rollback_error", "err", rbErr)
	}
}

func (s *proposalService) mapListingError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return derrors.NotFound("listing identity not found")
	}
	return derrors.Infra("listing repository failure", err)
}

func (s *proposalService) mapProposalError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return derrors.NotFound("proposal not found")
	}
	return derrors.Infra("proposal repository failure", err)
}

func (s *proposalService) actorCanViewProposal(actor Actor, proposal proposalmodel.ProposalInterface) bool {
	if actor.UserID <= 0 {
		return false
	}
	return actor.UserID == proposal.RealtorID() || actor.UserID == proposal.OwnerID()
}

func (s *proposalService) listProposals(ctx context.Context, scope proposalmodel.ActorScope, filter ListFilter) (ListResult, error) {
	ctx, spanEnd, tracerErr := utils.GenerateTracer(ctx)
	if tracerErr != nil {
		return ListResult{}, derrors.Infra("failed to start tracer", tracerErr)
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	repoFilter, err := s.buildRepositoryFilter(scope, filter)
	if err != nil {
		return ListResult{}, err
	}

	tx, txErr := s.globalSvc.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("proposal.list.tx_start_error", "err", txErr, "scope", scope)
		return ListResult{}, derrors.Infra("failed to start read-only transaction", txErr)
	}
	defer func() {
		if rbErr := s.globalSvc.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("proposal.list.tx_rollback_error", "err", rbErr)
		}
	}()

	repoResult, err := s.proposalRepo.ListProposals(ctx, tx, repoFilter)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("proposal.list.repo_error", "err", err, "scope", scope)
		return ListResult{}, derrors.Infra("failed to list proposals", err)
	}

	proposalIDs := make([]int64, 0, len(repoResult.Items))
	realtorIDs := make([]int64, 0, len(repoResult.Items))
	realtorSet := make(map[int64]struct{}, len(repoResult.Items))
	ownerIDs := make([]int64, 0, len(repoResult.Items))
	ownerSet := make(map[int64]struct{}, len(repoResult.Items))
	listingIDs := make([]int64, 0, len(repoResult.Items))
	listingSet := make(map[int64]struct{}, len(repoResult.Items))
	for _, proposal := range repoResult.Items {
		if proposal == nil {
			continue
		}
		proposalIDs = append(proposalIDs, proposal.ID())
		realtorID := proposal.RealtorID()
		if realtorID > 0 {
			if _, exists := realtorSet[realtorID]; !exists {
				realtorSet[realtorID] = struct{}{}
				realtorIDs = append(realtorIDs, realtorID)
			}
		}
		ownerID := proposal.OwnerID()
		if ownerID > 0 {
			if _, exists := ownerSet[ownerID]; !exists {
				ownerSet[ownerID] = struct{}{}
				ownerIDs = append(ownerIDs, ownerID)
			}
		}
		listingID := proposal.ListingIdentityID()
		if listingID > 0 {
			if _, exists := listingSet[listingID]; !exists {
				listingSet[listingID] = struct{}{}
				listingIDs = append(listingIDs, listingID)
			}
		}
	}

	documentsByProposal, err := s.loadDocumentsByProposals(ctx, tx, proposalIDs)
	if err != nil {
		return ListResult{}, err
	}

	realtorSummaries, err := s.loadRealtorSummaryMap(ctx, tx, realtorIDs)
	if err != nil {
		return ListResult{}, err
	}
	s.enrichRealtorSummaries(ctx, realtorSummaries)

	ownerSummaries, err := s.loadOwnerSummaryMap(ctx, tx, ownerIDs)
	if err != nil {
		return ListResult{}, err
	}
	s.enrichOwnerSummaries(ctx, ownerSummaries)

	listingsByIdentity, err := s.loadListingsByIdentity(ctx, tx, listingIDs)
	if err != nil {
		return ListResult{}, err
	}

	items := make([]ListItem, 0, len(repoResult.Items))
	for _, proposal := range repoResult.Items {
		if proposal == nil {
			continue
		}
		items = append(items, ListItem{
			Proposal:  proposal,
			Documents: documentsByProposal[proposal.ID()],
			Realtor:   s.getRealtorSummaryOrDefault(proposal.RealtorID(), realtorSummaries),
			Owner:     s.getOwnerSummaryOrDefault(proposal.OwnerID(), ownerSummaries),
			Listing:   listingsByIdentity[proposal.ListingIdentityID()],
		})
	}

	return ListResult{Items: items, Total: repoResult.Total}, nil
}

func (s *proposalService) loadListingsByIdentity(ctx context.Context, tx *sql.Tx, listingIDs []int64) (map[int64]listingmodel.ListingInterface, error) {
	result := make(map[int64]listingmodel.ListingInterface, len(listingIDs))

	if len(listingIDs) == 0 {
		return result, nil
	}

	logger := utils.LoggerFromContext(ctx)

	for _, listingID := range listingIDs {
		if listingID <= 0 {
			continue
		}

		listing, err := s.listingRepo.GetActiveListingVersion(ctx, tx, listingID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			utils.SetSpanError(ctx, err)
			logger.Error("proposal.list.listing_error", "err", err, "listing_identity_id", listingID)
			return nil, derrors.Infra("failed to load listing for proposal", err)
		}

		if listing != nil {
			result[listingID] = listing
		}
	}

	return result, nil
}

func (s *proposalService) buildRepositoryFilter(scope proposalmodel.ActorScope, filter ListFilter) (proposalmodel.ListFilter, error) {
	if filter.Actor.UserID <= 0 {
		return proposalmodel.ListFilter{}, derrors.Auth("actor metadata missing")
	}

	repoFilter := proposalmodel.ListFilter{
		ActorScope: scope,
		ActorID:    filter.Actor.UserID,
		ListingID:  filter.ListingID,
		Statuses:   filter.Statuses,
		Page:       filter.Page,
		Limit:      filter.PageSize,
	}

	repoFilter.Page = normalizePage(repoFilter.Page)
	repoFilter.Limit = normalizePageSize(repoFilter.Limit)

	return repoFilter, nil
}

func normalizePage(page int) int {
	if page < 1 {
		return 1
	}
	return page
}

func normalizePageSize(limit int) int {
	switch {
	case limit <= 0:
		return 20
	case limit > 100:
		return 100
	default:
		return limit
	}
}

func (s *proposalService) notifyOwnerNewProposal(ctx context.Context, proposal proposalmodel.ProposalInterface, listingCode uint32) {
	subject := "Nova proposta recebida"
	target := fmt.Sprintf("anúncio %d", proposal.ListingIdentityID())
	if listingCode > 0 {
		target = fmt.Sprintf("anúncio %d", listingCode)
	}
	body := fmt.Sprintf("Você recebeu uma nova proposta para o %s.", target)
	data := map[string]string{
		"event":             "proposal_created",
		"proposalId":        strconv.FormatInt(proposal.ID(), 10),
		"listingIdentityId": strconv.FormatInt(proposal.ListingIdentityID(), 10),
	}
	s.notifyUserDevices(ctx, proposal.OwnerID(), subject, body, data)
}

func (s *proposalService) notifyProposalStatusChange(ctx context.Context, proposal proposalmodel.ProposalInterface, userID int64, event, body string) {
	subject := map[string]string{
		"proposal_cancelled": "Proposta cancelada",
		"proposal_accepted":  "Proposta aceita",
		"proposal_rejected":  "Proposta rejeitada",
	}[event]
	if subject == "" {
		subject = "Atualização de proposta"
	}
	data := map[string]string{
		"event":             event,
		"proposalId":        strconv.FormatInt(proposal.ID(), 10),
		"listingIdentityId": strconv.FormatInt(proposal.ListingIdentityID(), 10),
		"status":            proposal.Status().String(),
	}
	s.notifyUserDevices(ctx, userID, subject, body, data)
}

func (s *proposalService) loadDocumentsByProposals(ctx context.Context, tx *sql.Tx, proposalIDs []int64) (map[int64][]proposalmodel.ProposalDocumentInterface, error) {
	if len(proposalIDs) == 0 {
		return map[int64][]proposalmodel.ProposalDocumentInterface{}, nil
	}

	documents, err := s.proposalRepo.ListDocumentsByProposalIDs(ctx, tx, proposalIDs, true)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("proposal.list.documents_error", "err", err)
		return nil, derrors.Infra("failed to load proposal documents", err)
	}
	if documents == nil {
		documents = make(map[int64][]proposalmodel.ProposalDocumentInterface)
	}
	return documents, nil
}

func (s *proposalService) loadRealtorSummaryMap(ctx context.Context, tx *sql.Tx, realtorIDs []int64) (map[int64]proposalmodel.RealtorSummary, error) {
	if len(realtorIDs) == 0 {
		return map[int64]proposalmodel.RealtorSummary{}, nil
	}

	summaries, err := s.proposalRepo.ListRealtorSummaries(ctx, tx, realtorIDs)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("proposal.realtor.summary_error", "err", err)
		return nil, derrors.Infra("failed to load realtor summaries", err)
	}

	result := make(map[int64]proposalmodel.RealtorSummary, len(summaries))
	for _, summary := range summaries {
		if summary == nil {
			continue
		}
		result[summary.ID()] = summary
	}
	return result, nil
}

func (s *proposalService) loadOwnerSummaryMap(ctx context.Context, tx *sql.Tx, ownerIDs []int64) (map[int64]proposalmodel.OwnerSummary, error) {
	if len(ownerIDs) == 0 {
		return map[int64]proposalmodel.OwnerSummary{}, nil
	}

	summaries, err := s.proposalRepo.ListOwnerSummaries(ctx, tx, ownerIDs)
	if err != nil {
		utils.SetSpanError(ctx, err)
		utils.LoggerFromContext(ctx).Error("proposal.owner.summary_error", "err", err)
		return nil, derrors.Infra("failed to load owner summaries", err)
	}

	result := make(map[int64]proposalmodel.OwnerSummary, len(summaries))
	for _, summary := range summaries {
		if summary == nil {
			continue
		}
		result[summary.ID()] = summary
	}
	return result, nil
}

func (s *proposalService) loadSingleOwnerSummary(ctx context.Context, tx *sql.Tx, ownerID int64) (proposalmodel.OwnerSummary, error) {
	if ownerID <= 0 {
		return nil, nil
	}
	result, err := s.loadOwnerSummaryMap(ctx, tx, []int64{ownerID})
	if err != nil {
		return nil, err
	}
	return s.getOwnerSummaryOrDefault(ownerID, result), nil
}

// enrichOwnerSummaries hydrates owner summaries with signed photo URLs when possible.
func (s *proposalService) enrichOwnerSummaries(ctx context.Context, cache map[int64]proposalmodel.OwnerSummary) {
	if len(cache) == 0 || s.userService == nil {
		return
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	for ownerID, summary := range cache {
		if summary == nil || ownerID <= 0 {
			continue
		}
		if summary.PhotoURL() != "" {
			continue
		}
		photoURL, err := s.generateRealtorPhotoURL(ctx, ownerID)
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Debug("proposal.owner.photo_url_error", "owner_id", ownerID, "err", err)
			continue
		}
		summary.SetPhotoURL(photoURL)
	}
}

func (s *proposalService) getOwnerSummaryOrDefault(ownerID int64, cache map[int64]proposalmodel.OwnerSummary) proposalmodel.OwnerSummary {
	if cache != nil {
		if summary, ok := cache[ownerID]; ok && summary != nil {
			return summary
		}
	}
	if ownerID <= 0 {
		return nil
	}
	return s.blankOwnerSummary(ownerID)
}

func (s *proposalService) blankOwnerSummary(ownerID int64) proposalmodel.OwnerSummary {
	summary := proposalmodel.NewOwnerSummary()
	summary.SetID(ownerID)
	return summary
}

func (s *proposalService) loadSingleRealtorSummary(ctx context.Context, tx *sql.Tx, realtorID int64) (proposalmodel.RealtorSummary, error) {
	if realtorID <= 0 {
		return nil, nil
	}
	result, err := s.loadRealtorSummaryMap(ctx, tx, []int64{realtorID})
	if err != nil {
		return nil, err
	}
	return s.getRealtorSummaryOrDefault(realtorID, result), nil
}

// enrichRealtorSummaries hydrates realtor summaries with signed photo URLs when possible.
func (s *proposalService) enrichRealtorSummaries(ctx context.Context, cache map[int64]proposalmodel.RealtorSummary) {
	if len(cache) == 0 || s.userService == nil {
		return
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	for realtorID, summary := range cache {
		if summary == nil || realtorID <= 0 {
			continue
		}
		if summary.PhotoURL() != "" {
			continue
		}
		photoURL, err := s.generateRealtorPhotoURL(ctx, realtorID)
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Debug("proposal.realtor.photo_url_error", "realtor_id", realtorID, "err", err)
			continue
		}
		summary.SetPhotoURL(photoURL)
	}
}

// generateRealtorPhotoURL impersonates the realtor user and requests a signed download URL.
func (s *proposalService) generateRealtorPhotoURL(ctx context.Context, userID int64) (string, error) {
	if s.userService == nil || userID <= 0 {
		return "", nil
	}

	impersonatedCtx := utils.SetUserInContext(ctx, usermodel.UserInfos{ID: userID})
	return s.userService.GetPhotoDownloadURL(impersonatedCtx, realtorPhotoVariant)
}

func (s *proposalService) getRealtorSummaryOrDefault(realtorID int64, cache map[int64]proposalmodel.RealtorSummary) proposalmodel.RealtorSummary {
	if cache != nil {
		if summary, ok := cache[realtorID]; ok && summary != nil {
			return summary
		}
	}
	if realtorID <= 0 {
		return nil
	}
	return s.blankRealtorSummary(realtorID)
}

func (s *proposalService) blankRealtorSummary(realtorID int64) proposalmodel.RealtorSummary {
	summary := proposalmodel.NewRealtorSummary()
	summary.SetID(realtorID)
	return summary
}

func (s *proposalService) notifyUserDevices(ctx context.Context, userID int64, subject, body string, data map[string]string) {
	if s.notifier == nil || s.globalSvc == nil || userID <= 0 {
		return
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tokens, err := s.globalSvc.ListDeviceTokensByUserIDIfOptedIn(ctx, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("proposal.notifications.tokens_error", "err", err, "user_id", userID)
		return
	}
	if len(tokens) == 0 {
		return
	}

	for _, token := range tokens {
		if token == "" {
			continue
		}
		payload := globalservice.NotificationRequest{
			Type:    globalservice.NotificationTypeFCM,
			Subject: subject,
			Body:    body,
			Token:   token,
			Data:    cloneData(data),
		}
		if err := s.notifier.SendNotification(ctx, payload); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("proposal.notifications.send_error", "err", err, "user_id", userID)
		}
	}
}

func (s *proposalService) recordOwnerProposalResponse(ctx context.Context, tx *sql.Tx, proposal proposalmodel.ProposalInterface, actionTime time.Time) error {
	if proposal == nil || proposal.OwnerID() <= 0 {
		return nil
	}
	if s.ownerMetrics == nil {
		return derrors.Infra("owner metrics repository unavailable", fmt.Errorf("owner metrics repository is nil"))
	}
	if current := proposal.FirstOwnerActionAt(); current.Valid {
		return nil
	}

	proposal.SetFirstOwnerActionAt(sql.NullTime{Valid: true, Time: actionTime})
	createdAt := proposal.CreatedAt()
	if createdAt.IsZero() {
		createdAt = actionTime
	}
	delta := actionTime.Sub(createdAt)
	if delta < 0 {
		delta = 0
	}
	const maxDelta = 24 * time.Hour * 365
	if delta > maxDelta {
		delta = maxDelta
	}
	input := ownermetricsrepository.ProposalResponseInput{
		OwnerID:      proposal.OwnerID(),
		DeltaSeconds: int64(delta / time.Second),
		RespondedAt:  actionTime,
	}
	if err := s.ownerMetrics.UpsertProposalResponse(ctx, tx, input); err != nil {
		return derrors.Infra("failed to persist owner proposal metrics", err)
	}
	return nil
}

func cloneData(data map[string]string) map[string]string {
	if len(data) == 0 {
		return nil
	}
	cloned := make(map[string]string, len(data))
	for key, value := range data {
		cloned[key] = value
	}
	return cloned
}
