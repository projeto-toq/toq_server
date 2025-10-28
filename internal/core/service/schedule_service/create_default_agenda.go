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

const (
	defaultMorningBlockEnd   = 480  // 08:00
	defaultEveningBlockStart = 1080 // 18:00
	minutesPerDay            = 1440
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
		logger.Error("schedule.create_default_agenda.missing_tx", "listing_id", input.ListingID)
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

	existing, err := s.scheduleRepo.GetAgendaByListingID(ctx, tx, input.ListingID)
	if err == nil && existing != nil {
		return existing, nil
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.create_default_agenda.get_agenda_error", "listing_id", input.ListingID, "err", err)
		return nil, utils.InternalError("")
	}

	agenda := schedulemodel.NewAgenda()
	agenda.SetListingID(input.ListingID)
	agenda.SetOwnerID(input.OwnerID)
	if loc, _ := utils.ResolveLocation("timezone", input.Timezone); loc != nil {
		agenda.SetTimezone(loc.String())
	} else {
		agenda.SetTimezone(strings.TrimSpace(input.Timezone))
	}

	id, err := s.scheduleRepo.InsertAgenda(ctx, tx, agenda)
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.create_default_agenda.insert_error", "listing_id", input.ListingID, "err", err)
		return nil, utils.InternalError("")
	}
	agenda.SetID(id)

	rules := buildDefaultBlockRules(agenda.ID())
	if len(rules) > 0 {
		if err := s.scheduleRepo.InsertRules(ctx, tx, rules); err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("schedule.create_default_agenda.rules_error", "listing_id", input.ListingID, "err", err)
			return nil, utils.InternalError("")
		}
	}

	return agenda, nil
}

func validateDefaultAgendaInput(input CreateDefaultAgendaInput) *utils.HTTPError {
	if input.ListingID <= 0 {
		return utils.ValidationError("listingId", "listingId must be greater than zero")
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

func buildDefaultBlockRules(agendaID uint64) []schedulemodel.AgendaRuleInterface {
	rules := make([]schedulemodel.AgendaRuleInterface, 0, 14)
	for weekday := time.Sunday; weekday <= time.Saturday; weekday++ {
		morning := schedulemodel.NewAgendaRule()
		morning.SetAgendaID(agendaID)
		morning.SetDayOfWeek(weekday)
		morning.SetStartMinutes(0)
		morning.SetEndMinutes(defaultMorningBlockEnd)
		morning.SetRuleType(schedulemodel.RuleTypeBlock)
		morning.SetActive(true)
		rules = append(rules, morning)

		evening := schedulemodel.NewAgendaRule()
		evening.SetAgendaID(agendaID)
		evening.SetDayOfWeek(weekday)
		evening.SetStartMinutes(defaultEveningBlockStart)
		evening.SetEndMinutes(minutesPerDay)
		evening.SetRuleType(schedulemodel.RuleTypeBlock)
		evening.SetActive(true)
		rules = append(rules, evening)
	}
	return rules
}
