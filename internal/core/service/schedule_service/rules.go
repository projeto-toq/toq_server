package scheduleservices

import (
	"context"
	"database/sql"
	"errors"
	"time"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *scheduleService) CreateRules(ctx context.Context, input CreateRuleInput) (RuleMutationResult, error) {
	if input.ListingIdentityID <= 0 {
		return RuleMutationResult{}, utils.ValidationError("listingIdentityId", "listingIdentityId must be greater than zero")
	}
	if input.OwnerID <= 0 {
		return RuleMutationResult{}, utils.ValidationError("ownerId", "ownerId must be greater than zero")
	}
	if input.ActorID <= 0 {
		return RuleMutationResult{}, utils.ValidationError("actorId", "actorId must be greater than zero")
	}
	if len(input.Weekdays) == 0 {
		return RuleMutationResult{}, utils.ValidationError("weekDays", "weekDays must contain at least one value")
	}
	if rngErr := validateRuleRangeMinutes(input.Range); rngErr != nil {
		return RuleMutationResult{}, rngErr
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return RuleMutationResult{}, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("schedule.rules.create.tx_start_error", "err", txErr, "listing_identity_id", input.ListingIdentityID)
		return RuleMutationResult{}, utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("schedule.rules.create.tx_rollback_error", "err", rbErr, "listing_identity_id", input.ListingIdentityID)
			}
		}
	}()

	agenda, err := s.scheduleRepo.GetAgendaByListingIdentityID(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return RuleMutationResult{}, utils.NotFoundError("Agenda")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.rules.create.get_agenda_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return RuleMutationResult{}, utils.InternalError("")
	}

	if agenda.OwnerID() != input.OwnerID {
		return RuleMutationResult{}, utils.AuthorizationError("Owner does not match listing agenda")
	}

	existingRules, err := s.scheduleRepo.ListRulesByAgenda(ctx, tx, agenda.ID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.rules.create.list_rules_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return RuleMutationResult{}, utils.InternalError("")
	}

	for _, weekday := range input.Weekdays {
		if ruleConflict(existingRules, weekday, input.Range, 0) {
			return RuleMutationResult{}, utils.ConflictError("Rule overlaps with an existing rule for the same weekday")
		}
	}

	newRules := make([]schedulemodel.AgendaRuleInterface, 0, len(input.Weekdays))
	for _, weekday := range input.Weekdays {
		rule := schedulemodel.NewAgendaRule()
		rule.SetAgendaID(agenda.ID())
		rule.SetDayOfWeek(weekday)
		rule.SetStartMinutes(input.Range.StartMinute)
		rule.SetEndMinutes(input.Range.EndMinute)
		rule.SetRuleType(schedulemodel.RuleTypeBlock)
		rule.SetActive(input.Active)
		newRules = append(newRules, rule)
	}

	if err := s.scheduleRepo.InsertRules(ctx, tx, newRules); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.rules.create.insert_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return RuleMutationResult{}, utils.InternalError("")
	}

	if cmErr := s.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("schedule.rules.create.tx_commit_error", "err", cmErr, "listing_identity_id", input.ListingIdentityID)
		return RuleMutationResult{}, utils.InternalError("")
	}

	committed = true

	return RuleMutationResult{ListingIdentityID: input.ListingIdentityID, Timezone: agenda.Timezone(), Rules: newRules}, nil
}

func (s *scheduleService) UpdateRule(ctx context.Context, input UpdateRuleInput) (schedulemodel.AgendaRuleInterface, error) {
	if input.RuleID == 0 {
		return nil, utils.ValidationError("ruleId", "ruleId must be greater than zero")
	}
	if input.ListingIdentityID <= 0 {
		return nil, utils.ValidationError("listingIdentityId", "listingIdentityId must be greater than zero")
	}
	if input.OwnerID <= 0 {
		return nil, utils.ValidationError("ownerId", "ownerId must be greater than zero")
	}
	if input.ActorID <= 0 {
		return nil, utils.ValidationError("actorId", "actorId must be greater than zero")
	}
	if rngErr := validateRuleRangeMinutes(input.Range); rngErr != nil {
		return nil, rngErr
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("schedule.rules.update.tx_start_error", "err", txErr, "listing_identity_id", input.ListingIdentityID)
		return nil, utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("schedule.rules.update.tx_rollback_error", "err", rbErr, "listing_identity_id", input.ListingIdentityID)
			}
		}
	}()

	agenda, err := s.scheduleRepo.GetAgendaByListingIdentityID(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NotFoundError("Agenda")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.rules.update.get_agenda_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return nil, utils.InternalError("")
	}

	if agenda.OwnerID() != input.OwnerID {
		return nil, utils.AuthorizationError("Owner does not match listing agenda")
	}

	rule, err := s.scheduleRepo.GetRuleByID(ctx, tx, input.RuleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NotFoundError("Agenda rule")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.rules.update.get_rule_error", "err", err, "rule_id", input.RuleID)
		return nil, utils.InternalError("")
	}

	if rule.AgendaID() != agenda.ID() {
		return nil, utils.AuthorizationError("Rule does not belong to provided listing")
	}

	existingRules, err := s.scheduleRepo.ListRulesByAgenda(ctx, tx, agenda.ID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.rules.update.list_rules_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return nil, utils.InternalError("")
	}

	if ruleConflict(existingRules, input.Weekday, input.Range, input.RuleID) {
		return nil, utils.ConflictError("Rule overlaps with an existing rule for the same weekday")
	}

	rule.SetDayOfWeek(input.Weekday)
	rule.SetStartMinutes(input.Range.StartMinute)
	rule.SetEndMinutes(input.Range.EndMinute)
	rule.SetRuleType(schedulemodel.RuleTypeBlock)
	rule.SetActive(input.Active)

	if err := s.scheduleRepo.UpdateRule(ctx, tx, rule); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.rules.update.exec_error", "err", err, "rule_id", input.RuleID)
		return nil, utils.InternalError("")
	}

	if cmErr := s.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("schedule.rules.update.tx_commit_error", "err", cmErr, "rule_id", input.RuleID)
		return nil, utils.InternalError("")
	}

	committed = true

	return rule, nil
}

func (s *scheduleService) DeleteRule(ctx context.Context, input DeleteRuleInput) error {
	if input.RuleID == 0 {
		return utils.ValidationError("ruleId", "ruleId must be greater than zero")
	}
	if input.ListingIdentityID <= 0 {
		return utils.ValidationError("listingIdentityId", "listingIdentityId must be greater than zero")
	}
	if input.OwnerID <= 0 {
		return utils.ValidationError("ownerId", "ownerId must be greater than zero")
	}
	if input.ActorID <= 0 {
		return utils.ValidationError("actorId", "actorId must be greater than zero")
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	tx, txErr := s.globalService.StartTransaction(ctx)
	if txErr != nil {
		utils.SetSpanError(ctx, txErr)
		logger.Error("schedule.rules.delete.tx_start_error", "err", txErr, "listing_identity_id", input.ListingIdentityID)
		return utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("schedule.rules.delete.tx_rollback_error", "err", rbErr, "listing_identity_id", input.ListingIdentityID)
			}
		}
	}()

	agenda, err := s.scheduleRepo.GetAgendaByListingIdentityID(ctx, tx, input.ListingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("Agenda")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.rules.delete.get_agenda_error", "err", err, "listing_identity_id", input.ListingIdentityID)
		return utils.InternalError("")
	}

	if agenda.OwnerID() != input.OwnerID {
		return utils.AuthorizationError("Owner does not match listing agenda")
	}

	rule, err := s.scheduleRepo.GetRuleByID(ctx, tx, input.RuleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundError("Agenda rule")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.rules.delete.get_rule_error", "err", err, "rule_id", input.RuleID)
		return utils.InternalError("")
	}

	if rule.AgendaID() != agenda.ID() {
		return utils.AuthorizationError("Rule does not belong to provided listing")
	}

	if err := s.scheduleRepo.DeleteRule(ctx, tx, input.RuleID); err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.rules.delete.exec_error", "err", err, "rule_id", input.RuleID)
		return utils.InternalError("")
	}

	if cmErr := s.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("schedule.rules.delete.tx_commit_error", "err", cmErr, "rule_id", input.RuleID)
		return utils.InternalError("")
	}

	committed = true
	return nil
}

func (s *scheduleService) ListRules(ctx context.Context, listingIdentityID, ownerID int64) (schedulemodel.RuleListResult, error) {
	if listingIdentityID <= 0 {
		return schedulemodel.RuleListResult{}, utils.ValidationError("listingIdentityId", "listingIdentityId must be greater than zero")
	}
	if ownerID <= 0 {
		return schedulemodel.RuleListResult{}, utils.ValidationError("ownerId", "ownerId must be greater than zero")
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
		logger.Error("schedule.rules.list.tx_start_error", "err", txErr, "listing_identity_id", listingIdentityID)
		return schedulemodel.RuleListResult{}, utils.InternalError("")
	}
	defer func() {
		if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
			utils.SetSpanError(ctx, rbErr)
			logger.Error("schedule.rules.list.tx_rollback_error", "err", rbErr, "listing_identity_id", listingIdentityID)
		}
	}()

	agenda, err := s.scheduleRepo.GetAgendaByListingIdentityID(ctx, tx, listingIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return schedulemodel.RuleListResult{}, utils.NotFoundError("Agenda")
		}
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.rules.list.get_agenda_error", "err", err, "listing_identity_id", listingIdentityID)
		return schedulemodel.RuleListResult{}, utils.InternalError("")
	}

	if agenda.OwnerID() != ownerID {
		return schedulemodel.RuleListResult{}, utils.AuthorizationError("Owner does not match listing agenda")
	}

	rules, err := s.scheduleRepo.ListRulesByAgenda(ctx, tx, agenda.ID())
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.rules.list.repo_error", "err", err, "listing_identity_id", listingIdentityID)
		return schedulemodel.RuleListResult{}, utils.InternalError("")
	}

	return schedulemodel.RuleListResult{ListingIdentityID: listingIdentityID, Timezone: agenda.Timezone(), Rules: rules}, nil
}

func validateRuleRangeMinutes(rng RuleTimeRange) *utils.HTTPError {
	if rng.StartMinute >= rng.EndMinute {
		return utils.ValidationError("range", "rangeStart must be before rangeEnd")
	}
	if rng.StartMinute >= minutesPerDay {
		return utils.ValidationError("range", "rangeStart must be less than 24:00")
	}
	if rng.EndMinute > minutesPerDay {
		return utils.ValidationError("range", "rangeEnd must be less or equal to 24:00")
	}
	return nil
}

func ruleConflict(existing []schedulemodel.AgendaRuleInterface, weekday time.Weekday, rng RuleTimeRange, ignoreID uint64) bool {
	for _, rule := range existing {
		if rule == nil {
			continue
		}
		if !rule.IsActive() {
			continue
		}
		if rule.DayOfWeek() != weekday {
			continue
		}
		if rule.ID() == ignoreID {
			continue
		}
		if intervalsOverlapMinutes(rule.StartMinutes(), rule.EndMinutes(), rng.StartMinute, rng.EndMinute) {
			return true
		}
	}
	return false
}

func intervalsOverlapMinutes(aStart, aEnd, bStart, bEnd uint16) bool {
	return uint32(aEnd) > uint32(bStart) && uint32(aStart) < uint32(bEnd)
}

func buildRuleWindow(base time.Time, rng RuleTimeRange, loc *time.Location) timeRange {
	start := time.Date(base.Year(), base.Month(), base.Day(), 0, 0, 0, 0, loc).Add(time.Duration(rng.StartMinute) * time.Minute)
	end := time.Date(base.Year(), base.Month(), base.Day(), 0, 0, 0, 0, loc).Add(time.Duration(rng.EndMinute) * time.Minute)
	return timeRange{start: start, end: end}
}
