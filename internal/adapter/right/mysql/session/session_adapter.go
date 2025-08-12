package sessionmysqladapter

import (
	mysqladapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql"
	sessionrepository "github.com/giulio-alfieri/toq_server/internal/core/port/right/repository/session_repository"
)

// Ensure implementation satisfies port interface
var _ sessionrepository.SessionRepoPortInterface = (*SessionAdapter)(nil)

type SessionAdapter struct {
	db *mysqladapter.Database
}

func NewSessionAdapter(db *mysqladapter.Database) sessionrepository.SessionRepoPortInterface {
	return &SessionAdapter{
		db: db,
	}
}
