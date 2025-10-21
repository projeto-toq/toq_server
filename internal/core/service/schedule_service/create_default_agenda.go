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
	if input.ListingID <= 0 {
		return nil, utils.ValidationError("listingId", "listingId must be greater than zero")
	}
	if input.OwnerID <= 0 {
		return nil, utils.ValidationError("ownerId", "ownerId must be greater than zero")
	}
	if strings.TrimSpace(input.Timezone) == "" {
		return nil, utils.ValidationError("timezone", "timezone is required")
	}
	if input.ActorID <= 0 {
		return nil, utils.ValidationError("actorId", "actorId must be greater than zero")
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

	existing, err := s.scheduleRepo.GetAgendaByListingID(ctx, tx, input.ListingID)
	if err == nil && existing != nil {
		return nil, utils.ConflictError("Agenda already exists for listing")
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		utils.SetSpanError(ctx, err)
		logger.Error("schedule.create_default_agenda.get_agenda_error", "listing_id", input.ListingID, "err", err)
		return nil, utils.InternalError("")
	}

	agenda := schedulemodel.NewAgenda()
	agenda.SetListingID(input.ListingID)
	agenda.SetOwnerID(input.OwnerID)
	agenda.SetTimezone(strings.TrimSpace(input.Timezone))

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

	if cmErr := s.globalService.CommitTransaction(ctx, tx); cmErr != nil {
		utils.SetSpanError(ctx, cmErr)
		logger.Error("schedule.create_default_agenda.tx_commit_error", "err", cmErr)
		return nil, utils.InternalError("")
	}

	committed = true

	return agenda, nil
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
