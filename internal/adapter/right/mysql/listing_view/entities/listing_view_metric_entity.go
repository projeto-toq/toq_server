package listingviewentity

import "time"

// ListingViewMetricEntity represents aggregated view counters for a listing identity.
//
// Schema expectations (managed by DBA):
//   - Table: listing_view_metrics (InnoDB, utf8mb4)
//   - Primary Key: listing_identity_id (INT/BIGINT) referencing listing_identities.id
//   - Columns:
//   - listing_identity_id (PK, NOT NULL)
//   - views BIGINT UNSIGNED NOT NULL DEFAULT 0
//   - last_view_at DATETIME NULL
//
// This entity is kept inside the MySQL adapter layer and should not leak to domain layers.
type ListingViewMetricEntity struct {
	ListingIdentityID int64
	Views             int64
	LastViewAt        *time.Time
}
