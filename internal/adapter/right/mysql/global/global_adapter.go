package mysqlglobaladapter

import (
	mysqlglobaladapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql"
)

type GlobalAdapter struct {
	db *mysqlglobaladapter.Database
}

func NewGlobalAdapter(db *mysqlglobaladapter.Database) *GlobalAdapter {
	return &GlobalAdapter{
		db: db,
	}
}
