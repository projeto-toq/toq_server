package scheduleservices

import (
	"context"
	"database/sql"
	"errors"
	"time"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *scheduleService) ListBlockRules(ctx context.Context, filter schedulemodel.BlockRulesFilter) (schedulemodel.RuleListResult, error) {
	if filter.OwnerID <= 0 {
		return schedulemodel.RuleListResult{}, utils.ValidationError("ownerId", "ownerId must be greater than zero")
	}
	if filter.ListingIdentityID <= 0 {
		return schedulemodel.RuleListResult{}, utils.ValidationError("listingIdentityId", "listingIdentityId must be greater than zero")
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return schedulemodel.RuleListResult{}, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartReadOnlyTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("schedule.list_block_rules.tx_start_error", "err", txErr, "listing_identity_id", filter.ListingIdentityID)
		return schedulemodel.RuleListResult{}, utils.InternalError("")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("schedule.list_block_rules.tx_rollback_error", "err", rbErr, "listing_identity_id", filter.ListingIdentityID)
		}
	}()

	agenda, err := s.scheduleRepo.GetAgendaByListingIdentityID(ctx, tx, filter.ListingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return schedulemodel.RuleListResult{}, utils.NotFoundError("Agenda")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.list_block_rules.get_agenda_error", "err", err, "listing_identity_id", filter.ListingIdentityID)
		return schedulemodel.RuleListResult{}, utils.InternalError("")
	}

	if agenda.OwnerID() != filter.OwnerID {
		return schedulemodel.RuleListResult{}, utils.AuthorizationError("Owner does not match listing agenda")
	}

	repoFilter := schedulemodel.BlockRulesFilter{
		OwnerID:           filter.OwnerID,
		ListingIdentityID: filter.ListingIdentityID,
		Weekdays:          uniqueWeekdays(filter.Weekdays),
	}

	rules, err := s.scheduleRepo.ListBlockRules(ctx, tx, repoFilter)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.list_block_rules.repo_error", "err", err, "listing_identity_id", filter.ListingIdentityID)
		return schedulemodel.RuleListResult{}, utils.InternalError("")
	}

	return schedulemodel.RuleListResult{
		ListingIdentityID: filter.ListingIdentityID,
		Timezone:          agenda.Timezone(),
		Rules:             rules,
	}, nil
}

func uniqueWeekdays(values []time.Weekday) []time.Weekday {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[time.Weekday]struct{}, len(values))
	result := make([]time.Weekday, 0, len(values))
	for _, weekday := range values {
		if _, ok := seen[weekday]; ok {
			continue
		}
		seen[weekday] = struct{}{}
		result = append(result, weekday)
	}
	return result
}
