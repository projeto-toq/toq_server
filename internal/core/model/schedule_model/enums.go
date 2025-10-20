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
