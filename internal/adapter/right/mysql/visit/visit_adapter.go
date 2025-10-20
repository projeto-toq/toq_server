package mysqlvisitadapter

import mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"

// VisitAdapter provides persistence access for listing visits.
type VisitAdapter struct {
	db *mysqladapter.Database
}

// NewVisitAdapter creates a new adapter.
func NewVisitAdapter(db *mysqladapter.Database) *VisitAdapter {
	return &VisitAdapter{db: db}
}
