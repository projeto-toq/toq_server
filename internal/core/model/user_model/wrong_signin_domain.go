package usermodel

import "time"

type wrongSigin struct {
	userid          int64
	failedAttempts  int64
	last_attempt_at time.Time
}

func (w *wrongSigin) GetUserID() int64 {
	return w.userid
}

func (w *wrongSigin) SetUserID(userid int64) {
	w.userid = userid
}

func (w *wrongSigin) GetFailedAttempts() int64 {
	return w.failedAttempts
}

func (w *wrongSigin) SetFailedAttempts(failedAttempts int64) {
	w.failedAttempts = failedAttempts
}

func (w *wrongSigin) GetLastAttemptAt() time.Time {
	return w.last_attempt_at
}

func (w *wrongSigin) SetLastAttemptAt(lastAttemptAt time.Time) {
	w.last_attempt_at = lastAttemptAt
}
