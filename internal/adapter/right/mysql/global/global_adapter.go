package mysqlglobaladapter

import (
	mysqlglobaladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
)

type GlobalAdapter struct {
	db *mysqlglobaladapter.Database
}

func NewGlobalAdapter(db *mysqlglobaladapter.Database) *GlobalAdapter {
	return &GlobalAdapter{
		db: db,
	}
}
