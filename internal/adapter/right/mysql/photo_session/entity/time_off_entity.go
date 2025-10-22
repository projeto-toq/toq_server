package entity

import "time"

// TimeOffEntity maps photographer_time_off rows.
type TimeOffEntity struct {
	ID                 uint64
	PhotographerUserID uint64
	StartDate          time.Time
	EndDate            time.Time
	Reason             *string
}
