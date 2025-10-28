package userservices

import (
	"strings"
	"time"
)

// Config aggregates runtime settings consumed by the user service.
type Config struct {
	SystemUserResetPasswordURL        string
	PhotographerTimezone              string
	PhotographerAgendaHorizonMonths   int
	PhotographerAgendaRefreshInterval time.Duration
}

func normalizeConfig(cfg Config) Config {
	cfg.SystemUserResetPasswordURL = strings.TrimSpace(cfg.SystemUserResetPasswordURL)
	if cfg.SystemUserResetPasswordURL == "" {
		cfg.SystemUserResetPasswordURL = "https://gca.dev.br/app/#/password/request"
	}
	cfg.PhotographerTimezone = strings.TrimSpace(cfg.PhotographerTimezone)
	if cfg.PhotographerTimezone == "" {
		cfg.PhotographerTimezone = "America/Sao_Paulo"
	}
	if cfg.PhotographerAgendaHorizonMonths <= 0 {
		cfg.PhotographerAgendaHorizonMonths = 3
	}
	if cfg.PhotographerAgendaRefreshInterval <= 0 {
		cfg.PhotographerAgendaRefreshInterval = 24 * time.Hour
	}
	return cfg
}
