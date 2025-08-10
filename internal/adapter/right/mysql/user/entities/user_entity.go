package userentity

import (
	"database/sql"
	"time"
)

type UserEntity struct {
	ID                int64
	FullName          string
	NickName          sql.NullString
	NationalID        string
	CreciNumber       sql.NullString
	CreciState        sql.NullString
	CreciValidity     sql.NullTime
	BornAT            time.Time
	PhoneNumber       string
	Email             string
	ZipCode           string
	Street            string
	Number            string
	Complement        sql.NullString
	Neighborhood      string
	City              string
	State             string
	Photo             sql.NullString
	Password          string
	OptStatus         bool
	LastActivityAT    time.Time
	Deleted           bool
	LastSignInAttempt sql.NullTime
}
