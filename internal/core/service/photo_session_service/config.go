package photosessionservices

// Config holds tunable parameters for slot generation.
type Config struct {
	SlotDurationMinutes         int
	SlotsPerPeriod              int
	MorningStartHour            int
	AfternoonStartHour          int
	BusinessStartHour           int
	BusinessEndHour             int
	AgendaHorizonMonths         int
	RequirePhotographerApproval bool
}
