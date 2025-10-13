package globalservice

import (
	"context"
	"database/sql"

	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

func (gs *globalService) ListCatalogValues(ctx context.Context, tx *sql.Tx, category string, includeInactive bool) ([]listingmodel.CatalogValueInterface, error) {
	return gs.globalRepo.ListCatalogValues(ctx, tx, category, includeInactive)
}

func (gs *globalService) GetCatalogValueByID(ctx context.Context, tx *sql.Tx, category string, id uint8) (listingmodel.CatalogValueInterface, error) {
	return gs.globalRepo.GetCatalogValueByID(ctx, tx, category, id)
}

func (gs *globalService) GetCatalogValueBySlug(ctx context.Context, tx *sql.Tx, category, slug string) (listingmodel.CatalogValueInterface, error) {
	return gs.globalRepo.GetCatalogValueBySlug(ctx, tx, category, slug)
}

func (gs *globalService) GetNextCatalogValueID(ctx context.Context, tx *sql.Tx, category string) (uint8, error) {
	return gs.globalRepo.GetNextCatalogValueID(ctx, tx, category)
}

func (gs *globalService) CreateCatalogValue(ctx context.Context, tx *sql.Tx, value listingmodel.CatalogValueInterface) error {
	return gs.globalRepo.CreateCatalogValue(ctx, tx, value)
}

func (gs *globalService) UpdateCatalogValue(ctx context.Context, tx *sql.Tx, value listingmodel.CatalogValueInterface) error {
	return gs.globalRepo.UpdateCatalogValue(ctx, tx, value)
}

func (gs *globalService) SoftDeleteCatalogValue(ctx context.Context, tx *sql.Tx, category string, id uint8) error {
	return gs.globalRepo.SoftDeleteCatalogValue(ctx, tx, category, id)
}
