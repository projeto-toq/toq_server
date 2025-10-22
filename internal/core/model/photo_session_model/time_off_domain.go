package photosessionmodel

import "time"

type photographerTimeOff struct {
	id                 uint64
	photographerUserID uint64
	startDate          time.Time
	endDate            time.Time
	reason             *string
}

func (t *photographerTimeOff) ID() uint64 { return t.id }

func (t *photographerTimeOff) SetID(value uint64) { t.id = value }

func (t *photographerTimeOff) PhotographerUserID() uint64 { return t.photographerUserID }

func (t *photographerTimeOff) SetPhotographerUserID(value uint64) { t.photographerUserID = value }

func (t *photographerTimeOff) StartDate() time.Time { return t.startDate }

func (t *photographerTimeOff) SetStartDate(value time.Time) { t.startDate = value }

func (t *photographerTimeOff) EndDate() time.Time { return t.endDate }

func (t *photographerTimeOff) SetEndDate(value time.Time) { t.endDate = value }

func (t *photographerTimeOff) Reason() *string { return t.reason }

func (t *photographerTimeOff) SetReason(value *string) { t.reason = value }
