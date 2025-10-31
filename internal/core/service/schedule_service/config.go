package scheduleservices

import (
	"fmt"
	"strconv"
	"strings"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
)

const minutesPerDay = 24 * 60

// Config represents runtime options consumed by the schedule service.
// DefaultBlockRuleRanges holds the recurring rule windows created for new agendas.
// Minutes use a 0-1440 range where the end minute is exclusive.
type Config struct {
	DefaultBlockRuleRanges []RuleTimeRange
}

// DefaultConfig returns the built-in configuration when no overrides are provided.
func DefaultConfig() Config {
	return Config{
		DefaultBlockRuleRanges: []RuleTimeRange{
			{StartMinute: 0, EndMinute: 480},     // 00:00 - 08:00 exclusive (covers up to 07:59)
			{StartMinute: 1140, EndMinute: 1440}, // 19:00 - 24:00 exclusive (covers up to 23:59)
		},
	}
}

func (c Config) ensureDefaults() Config {
	if len(c.DefaultBlockRuleRanges) == 0 {
		return DefaultConfig()
	}
	cleaned := make([]RuleTimeRange, 0, len(c.DefaultBlockRuleRanges))
	for _, rng := range c.DefaultBlockRuleRanges {
		cleaned = append(cleaned, RuleTimeRange{
			StartMinute: rng.StartMinute,
			EndMinute:   rng.EndMinute,
		})
	}
	c.DefaultBlockRuleRanges = cleaned
	return c
}

// ConfigFromEnvironment converts the environment YAML representation into a service configuration.
func ConfigFromEnvironment(env *globalmodel.Environment) (Config, error) {
	if env == nil {
		return DefaultConfig(), nil
	}
	entries := env.Schedule.DefaultBlockRules
	if len(entries) == 0 {
		return DefaultConfig(), nil
	}

	ranges := make([]RuleTimeRange, 0, len(entries))
	for idx, entry := range entries {
		start, err := parseConfigTime(entry.Start)
		if err != nil {
			return Config{}, fmt.Errorf("schedule.default_block_rules[%d].start: %w", idx, err)
		}
		end, err := parseConfigTime(entry.End)
		if err != nil {
			return Config{}, fmt.Errorf("schedule.default_block_rules[%d].end: %w", idx, err)
		}
		if start >= end {
			return Config{}, fmt.Errorf("schedule.default_block_rules[%d]: start must be before end", idx)
		}
		ranges = append(ranges, RuleTimeRange{StartMinute: start, EndMinute: end})
	}

	if len(ranges) == 0 {
		return DefaultConfig(), nil
	}

	return Config{DefaultBlockRuleRanges: ranges}, nil
}

func parseConfigTime(raw string) (uint16, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return 0, fmt.Errorf("value must be in the format HH:MM")
	}
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("value must be in the format HH:MM")
	}
	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hour component")
	}
	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid minute component")
	}
	if hour < 0 || hour > 24 {
		return 0, fmt.Errorf("hour must be between 00 and 24")
	}
	if minute < 0 || minute >= 60 {
		return 0, fmt.Errorf("minutes must be between 00 and 59")
	}
	if hour == 24 && minute != 0 {
		return 0, fmt.Errorf("24 hour notation must be 24:00")
	}
	total := hour*60 + minute
	if total > minutesPerDay {
		return 0, fmt.Errorf("time must be less or equal to 24:00")
	}
	return uint16(total), nil
}
