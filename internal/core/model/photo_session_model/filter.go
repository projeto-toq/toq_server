package photosessionmodel

import "time"

// SlotListParams encapsulates Filtering/pagination parameters for querying slots.
type SlotListParams struct {
	From          *time.Time
	To            *time.Time
	Period        *SlotPeriod
	Limit         int
	Offset        int
	SortColumn    string
	SortDirection string
}
