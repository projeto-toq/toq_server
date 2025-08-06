package mysqlcomplexadapter

import (
	mysqlcomplexadapter "github.com/giulio-alfieri/toq_server/internal/adapter/right/mysql"
)

type ComplexAdapter struct {
	db *mysqlcomplexadapter.Database
}

func NewComplexAdapter(db *mysqlcomplexadapter.Database) *ComplexAdapter {
	return &ComplexAdapter{
		db: db,
	}
}
