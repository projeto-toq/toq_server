package photosessionservices

// Config holds tunable parameters for slot generation.
type Config struct {
	SlotDurationMinutes int
	SlotsPerPeriod      int
	MorningStartHour    int
	AfternoonStartHour  int
	BusinessStartHour   int
	BusinessEndHour     int
	AgendaHorizonMonths int
}

func normalizeConfig(cfg Config) Config {
	if cfg.SlotDurationMinutes <= 0 {
		cfg.SlotDurationMinutes = 240
	}
	if cfg.SlotsPerPeriod <= 0 {
		cfg.SlotsPerPeriod = 1
	}
	if cfg.MorningStartHour <= 0 {
		cfg.MorningStartHour = 8
	}
	if cfg.AfternoonStartHour <= 0 {
		cfg.AfternoonStartHour = 14
	}
	if cfg.BusinessStartHour <= 0 {
		cfg.BusinessStartHour = defaultWorkdayStartHour
	}
	if cfg.BusinessEndHour <= cfg.BusinessStartHour {
		cfg.BusinessEndHour = defaultWorkdayEndHour
	}
	if cfg.AgendaHorizonMonths <= 0 {
		cfg.AgendaHorizonMonths = defaultHorizonMonths
	}
	return cfg
}
