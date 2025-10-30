package schedulemodel

// RuleType describes the effect applied by a recurring rule.
type RuleType string

const (
	// RuleTypeBlock marks the interval as unavailable until overridden by a specific entry.
	RuleTypeBlock RuleType = "BLOCK"
	// RuleTypeFree marks the interval as free and informative.
	RuleTypeFree RuleType = "FREE"
)

// EntryType describes the source that created a specific agenda entry.
type EntryType string

const (
	EntryTypeBlock          EntryType = "BLOCK"
	EntryTypeTemporaryBlock EntryType = "TEMP_BLOCK"
	EntryTypeVisitPending   EntryType = "VISIT_PENDING"
	EntryTypeVisitConfirmed EntryType = "VISIT_CONFIRMED"
	EntryTypePhotoSession   EntryType = "PHOTO_SESSION"
	EntryTypeHolidayInfo    EntryType = "HOLIDAY_INFO"
)

// TimelineSource identifies the origin of a timeline window.
type TimelineSource string

const (
	// TimelineSourceEntry indicates the window comes from a persisted agenda entry.
	TimelineSourceEntry TimelineSource = "ENTRY"
	// TimelineSourceRule indicates the window was generated from a recurring rule.
	TimelineSourceRule TimelineSource = "RULE"
)
