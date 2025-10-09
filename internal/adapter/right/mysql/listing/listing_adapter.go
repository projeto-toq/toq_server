package mysqllistingadapter

import (
	mysqllistingadapter "github.com/projeto-toq/toq_server/internal/adapter/right/mysql"
)

type ListingAdapter struct {
	db *mysqllistingadapter.Database
}

func NewListingAdapter(db *mysqllistingadapter.Database) *ListingAdapter {
	return &ListingAdapter{
		db: db,
	}
}
