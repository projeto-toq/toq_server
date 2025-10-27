package photosessionservices

// Config holds tunable parameters for slot generation.
type Config struct {
	SlotDurationMinutes int
	SlotsPerPeriod      int
	MorningStartHour    int
	AfternoonStartHour  int
}

func normalizeConfig(cfg Config) Config {
	if cfg.SlotDurationMinutes <= 0 {
		cfg.SlotDurationMinutes = 60
	}
	if cfg.SlotsPerPeriod <= 0 {
		cfg.SlotsPerPeriod = 4
	}
	if cfg.MorningStartHour <= 0 {
		cfg.MorningStartHour = 8
	}
	if cfg.AfternoonStartHour <= 0 {
		cfg.AfternoonStartHour = 14
	}
	return cfg
}
