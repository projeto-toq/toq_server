package scheduleservices

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	schedulemodel "github.com/projeto-toq/toq_server/internal/core/model/schedule_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *scheduleService) CreateDefaultAgenda(ctx context.Context, input CreateDefaultAgendaInput) (schedulemodel.AgendaInterface, error) {
	if err := validateDefaultAgendaInput(input); err != nil {
		return nil, err
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
		logger.Error("schedule.create_default_agenda.tx_start_error", "err", txErr)
		return nil, utils.InternalError("")
	}

	committed := false
	defer func() {
		if !committed {
			if rbErr := s.globalService.RollbackTransaction(ctx, tx); rbErr != nil {
				utils.SetSpanError(ctx, rbErr)
				logger.Error("schedule.create_default_agenda.tx_rollback_error", "err", rbErr)
			}
		}
	}()

	agenda, opErr := s.createDefaultAgendaTx(ctx, tx, input)
	if opErr != nil {
		return nil, opErr
	}

	if cmErr := s.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("schedule.create_default_agenda.tx_commit_error", "err", cmErr)
		return nil, utils.InternalError("")
	}

	committed = true

	return agenda, nil
}

func (s *scheduleService) CreateDefaultAgendaWithTx(ctx context.Context, tx *sql.Tx, input CreateDefaultAgendaInput) (schedulemodel.AgendaInterface, error) {
	if tx == nil {
		ctx = utils.ContextWithLogger(ctx)
		logger := utils.LoggerFromContext(ctx)
		logger.Error("schedule.create_default_agenda.missing_tx", "listing_identity_id", input.ListingIdentityID)
		return nil, utils.InternalError("")
	}
	if err := validateDefaultAgendaInput(input); err != nil {
		return nil, err
	}

	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, utils.InternalError("")
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)

	return s.createDefaultAgendaTx(ctx, tx, input)
}

func (s *scheduleService) createDefaultAgendaTx(ctx context.Context, tx *sql.Tx, input CreateDefaultAgendaInput) (schedulemodel.AgendaInterface, error) {
	logger := utils.LoggerFromContext(ctx)

	existing, err := s.scheduleRepo.GetAgendaByListingIdentityID(ctx, tx, input.ListingIdentityID)
	if err == nil && existing != nil {
		return existing, nil
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.create_default_agenda.get_agenda_error", "listing_identity_id", input.ListingIdentityID, "err", err)
		return nil, utils.InternalError("")
	}

	agenda := schedulemodel.NewAgenda()
	agenda.SetListingIdentityID(input.ListingIdentityID)
	agenda.SetOwnerID(input.OwnerID)
	if loc, _ := utils.ResolveLocation("timezone", input.Timezone); loc != nil {
		agenda.SetTimezone(loc.String())
	} else {
		agenda.SetTimezone(strings.TrimSpace(input.Timezone))
	}

	id, err := s.scheduleRepo.InsertAgenda(ctx, tx, agenda)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.create_default_agenda.insert_error", "listing_identity_id", input.ListingIdentityID, "err", err)
		return nil, utils.InternalError("")
	}
	agenda.SetID(id)

	rules := s.buildDefaultBlockRules(agenda.ID())
	if len(rules) > 0 {
		if err := s.scheduleRepo.InsertRules(ctx, tx, rules); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("schedule.create_default_agenda.rules_error", "listing_identity_id", input.ListingIdentityID, "err", err)
			return nil, utils.InternalError("")
		}
	}

	return agenda, nil
}

func validateDefaultAgendaInput(input CreateDefaultAgendaInput) *utils.HTTPError {
	if input.ListingIdentityID <= 0 {
		return utils.ValidationError("listingIdentityId", "listingIdentityId must be greater than zero")
	}
	if input.OwnerID <= 0 {
		return utils.ValidationError("ownerId", "ownerId must be greater than zero")
	}
	if strings.TrimSpace(input.Timezone) == "" {
		return utils.ValidationError("timezone", "timezone is required")
	}
	if _, err := utils.ResolveLocation("timezone", input.Timezone); err != nil {
		return err
	}
	if input.ActorID <= 0 {
		return utils.ValidationError("actorId", "actorId must be greater than zero")
	}
	return nil
}

func (s *scheduleService) buildDefaultBlockRules(agendaID uint64) []schedulemodel.AgendaRuleInterface {
	ranges := s.defaultBlockRuleRanges
	if len(ranges) == 0 {
		ranges = DefaultConfig().DefaultBlockRuleRanges
	}
	rules := make([]schedulemodel.AgendaRuleInterface, 0, len(ranges)*7)
	for weekday := time.Sunday; weekday <= time.Saturday; weekday++ {
		for _, rng := range ranges {
			rule := schedulemodel.NewAgendaRule()
			rule.SetAgendaID(agendaID)
			rule.SetDayOfWeek(weekday)
			rule.SetStartMinutes(rng.StartMinute)
			rule.SetEndMinutes(rng.EndMinute)
			rule.SetRuleType(schedulemodel.RuleTypeBlock)
			rule.SetActive(true)
			rules = append(rules, rule)
		}
	}
	return rules
}
