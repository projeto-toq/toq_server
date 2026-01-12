package mysqlvisitadapter

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/visit/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/visit/entities"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
)

// scanVisitWithListingRow hydrates a VisitWithListing struct with listing and participant snapshots.
func scanVisitWithListingRow(scanner rowScanner) (listingmodel.VisitWithListing, error) {
	var visitEntity entities.VisitEntity
	var listingID sql.NullInt64
	var listingVersion sql.NullInt64
	var listingZip sql.NullString
	var listingStreet sql.NullString
	var listingNumber sql.NullString
	var listingComplement sql.NullString
	var listingNeighborhood sql.NullString
	var listingCity sql.NullString
	var listingState sql.NullString
	var listingTitle sql.NullString
	var listingDescription sql.NullString
	var ownerFullName sql.NullString
	var ownerCreatedAt sql.NullTime
	var ownerAvgResponse sql.NullInt64
	var realtorFullName sql.NullString
	var realtorCreatedAt sql.NullTime
	var realtorTotalVisits sql.NullInt64

	if err := scanner.Scan(
		&visitEntity.ID,
		&visitEntity.ListingIdentityID,
		&visitEntity.ListingVersion,
		&visitEntity.RequesterUserID,
		&visitEntity.OwnerUserID,
		&visitEntity.ScheduledStart,
		&visitEntity.ScheduledEnd,
		&visitEntity.Status,
		&visitEntity.Source,
		&visitEntity.Notes,
		&visitEntity.RejectionReason,
		&visitEntity.FirstOwnerActionAt,
		&visitEntity.RequestedAt,
		&listingID,
		&listingVersion,
		&listingZip,
		&listingStreet,
		&listingNumber,
		&listingComplement,
		&listingNeighborhood,
		&listingCity,
		&listingState,
		&listingTitle,
		&listingDescription,
		&ownerFullName,
		&ownerCreatedAt,
		&ownerAvgResponse,
		&realtorFullName,
		&realtorCreatedAt,
		&realtorTotalVisits,
	); err != nil {
		return listingmodel.VisitWithListing{}, err
	}

	if !listingID.Valid || !listingVersion.Valid {
		return listingmodel.VisitWithListing{}, fmt.Errorf("listing version not found for listing_identity_id=%d", visitEntity.ListingIdentityID)
	}

	requiredSnapshotFields := []struct {
		valid bool
		name  string
	}{
		{listingZip.Valid, "zip_code"},
		{listingStreet.Valid, "street"},
		{listingNeighborhood.Valid, "neighborhood"},
		{listingCity.Valid, "city"},
		{listingState.Valid, "state"},
	}

	for _, field := range requiredSnapshotFields {
		if field.valid {
			continue
		}
		return listingmodel.VisitWithListing{}, fmt.Errorf("listing snapshot missing %s for listing_identity_id=%d", field.name, visitEntity.ListingIdentityID)
	}

	visitModel := converters.ToVisitModel(visitEntity)
	listingModel := listingmodel.NewListing()
	listingModel.SetID(listingID.Int64)
	listingModel.SetListingIdentityID(visitEntity.ListingIdentityID)
	listingModel.SetVersion(uint8(listingVersion.Int64))
	listingModel.SetZipCode(listingZip.String)
	listingModel.SetStreet(listingStreet.String)
	if listingNumber.Valid {
		listingModel.SetNumber(listingNumber.String)
	}
	if listingComplement.Valid {
		trimmed := strings.TrimSpace(listingComplement.String)
		if trimmed != "" {
			listingModel.SetComplement(trimmed)
		}
	}
	listingModel.SetNeighborhood(listingNeighborhood.String)
	listingModel.SetCity(listingCity.String)
	listingModel.SetState(listingState.String)
	if listingTitle.Valid {
		listingModel.SetTitle(listingTitle.String)
	} else {
		listingModel.UnsetTitle()
	}
	if listingDescription.Valid {
		listingModel.SetDescription(listingDescription.String)
	} else {
		listingModel.UnsetDescription()
	}

	ownerSnapshot := listingmodel.VisitParticipantSnapshot{
		UserID:             visitEntity.OwnerUserID,
		FullName:           strings.TrimSpace(ownerFullName.String),
		AvgResponseSeconds: ownerAvgResponse,
	}
	if ownerCreatedAt.Valid {
		ownerSnapshot.CreatedAt = ownerCreatedAt.Time
	}

	realtorSnapshot := listingmodel.VisitParticipantSnapshot{
		UserID:      visitEntity.RequesterUserID,
		FullName:    strings.TrimSpace(realtorFullName.String),
		TotalVisits: realtorTotalVisits,
	}
	if realtorCreatedAt.Valid {
		realtorSnapshot.CreatedAt = realtorCreatedAt.Time
	}

	return listingmodel.VisitWithListing{
		Visit:   visitModel,
		Listing: listingModel,
		Owner:   ownerSnapshot,
		Realtor: realtorSnapshot,
	}, nil
}
