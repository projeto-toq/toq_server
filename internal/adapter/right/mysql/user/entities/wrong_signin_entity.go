package userentity

import "time"

type WrongSignInEntity struct {
	UserID         int64
	FailedAttempts uint8
	LastAttemptAT  time.Time
}
