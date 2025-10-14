package userservices

import (
	"database/sql"
	"strings"
)

func firstToken(text string) string {
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

func errorsIsNoRows(err error) bool {
	return errorsIs(err, sql.ErrNoRows)
}

func errorsIs(err, target error) bool { //nolint:revive
	type is interface{ Is(error) bool }
	if x, ok := err.(is); ok {
		return x.Is(target)
	}
	for e := err; e != nil; {
		if e == target {
			return true
		}
		type unw interface{ Unwrap() error }
		if u, ok := e.(unw); ok {
			e = u.Unwrap()
			continue
		}
		break
	}
	return false
}
