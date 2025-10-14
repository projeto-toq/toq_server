package mysqlphotosessionadapter

import mysqladapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"

// PhotoSessionAdapter provides DB access to photographer slots and bookings.
type PhotoSessionAdapter struct {
	db *mysqladapter.Database
}

// NewPhotoSessionAdapter builds a new adapter instance.
func NewPhotoSessionAdapter(db *mysqladapter.Database) *PhotoSessionAdapter {
	return &PhotoSessionAdapter{db: db}
}
