package mysqluseradapter

import (
	mysqluseradapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql"
)

type UserAdapter struct {
	db *mysqluseradapter.Database
}

func NewUserAdapter(db *mysqluseradapter.Database) *UserAdapter {
	return &UserAdapter{
		db: db,
	}
}
