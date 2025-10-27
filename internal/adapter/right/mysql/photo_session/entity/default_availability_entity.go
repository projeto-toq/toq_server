package entity

// DefaultAvailabilityEntity mirrors the photographer_default_availability table.
type DefaultAvailabilityEntity struct {
	ID                 uint64
	PhotographerUserID uint64
	Weekday            int
	Period             string
	StartHour          int
	SlotsPerPeriod     int
	SlotDurationMin    int
}
