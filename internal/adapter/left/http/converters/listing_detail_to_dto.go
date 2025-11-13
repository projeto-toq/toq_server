package converters

import (
	"strings"

	dto "github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
)

// ListingDetailToDTO converte o retorno do service para o DTO exposto pelo handler.
func ListingDetailToDTO(detail listingservices.ListingDetailOutput) dto.ListingDetailResponse {
	listing := detail.Listing
	title := strings.TrimSpace(listing.Title())

	resp := dto.ListingDetailResponse{
		ID:                 listing.ID(),
		ListingIdentityID:  listing.IdentityID(),
		ListingUUID:        listing.UUID(),
		ActiveVersionID:    listing.ActiveVersionID(),
		UserID:             listing.UserID(),
		Code:               listing.Code(),
		Version:            listing.Version(),
		Status:             listing.Status().String(),
		ZipCode:            listing.ZipCode(),
		Street:             listing.Street(),
		Number:             listing.Number(),
		Complement:         listing.Complement(),
		Neighborhood:       listing.Neighborhood(),
		City:               listing.City(),
		State:              listing.State(),
		Title:              title,
		LandSize:           listing.LandSize(),
		Corner:             listing.Corner(),
		NonBuildable:       listing.NonBuildable(),
		Buildable:          listing.Buildable(),
		Description:        listing.Description(),
		SellNet:            listing.SellNet(),
		RentNet:            listing.RentNet(),
		Condominium:        listing.Condominium(),
		AnnualTax:          listing.AnnualTax(),
		MonthlyTax:         listing.MonthlyTax(),
		AnnualGroundRent:   listing.AnnualGroundRent(),
		MonthlyGroundRent:  listing.MonthlyGroundRent(),
		Exchange:           listing.Exchange(),
		ExchangePercentual: listing.ExchangePercentual(),
		Financing:          listing.Financing(),
		TenantName:         listing.TenantName(),
		TenantEmail:        listing.TenantEmail(),
		TenantPhone:        listing.TenantPhone(),
		Deleted:            listing.Deleted(),
		PhotoSessionID:     detail.PhotoSessionID,
	}

	if option, ok := listingmodel.PropertyTypeOptionFromBit(listing.ListingType()); ok {
		resp.PropertyType = &dto.ListingPropertyTypeResponse{
			Code:        option.Code,
			Label:       option.Label,
			PropertyBit: uint16(option.PropertyBit),
		}
	} else {
		resp.PropertyType = &dto.ListingPropertyTypeResponse{PropertyBit: uint16(listing.ListingType())}
	}

	resp.Owner = catalogDetailToPointer(detail.Owner, uint8(listing.Owner()))
	resp.Delivered = catalogDetailToPointer(detail.Delivered, uint8(listing.Delivered()))
	resp.WhoLives = catalogDetailToPointer(detail.WhoLives, uint8(listing.WhoLives()))
	resp.Transaction = catalogDetailToPointer(detail.Transaction, uint8(listing.Transaction()))
	resp.Installment = catalogDetailToPointer(detail.Installment, uint8(listing.Installment()))
	resp.Visit = catalogDetailToPointer(detail.Visit, uint8(listing.Visit()))
	resp.Accompanying = catalogDetailToPointer(detail.Accompanying, uint8(listing.Accompanying()))

	resp.Features = make([]dto.ListingFeatureResponse, 0, len(detail.Features))
	for _, feature := range detail.Features {
		resp.Features = append(resp.Features, dto.ListingFeatureResponse{
			Feature:     feature.Feature,
			Description: strings.TrimSpace(feature.Description),
			Quantity:    feature.Quantity,
		})
	}

	exchangePlaces := listing.ExchangePlaces()
	resp.ExchangePlaces = make([]dto.ListingExchangePlaceResponse, 0, len(exchangePlaces))
	for _, place := range exchangePlaces {
		resp.ExchangePlaces = append(resp.ExchangePlaces, dto.ListingExchangePlaceResponse{
			ID:               place.ID(),
			ListingID:        place.ListingID(),
			ListingVersionID: place.ListingVersionID(),
			Neighborhood:     strings.TrimSpace(place.Neighborhood()),
			City:             strings.TrimSpace(place.City()),
			State:            strings.TrimSpace(place.State()),
		})
	}

	resp.FinancingBlockers = make([]dto.ListingFinancingBlockerResponse, 0, len(detail.FinancingBlockers))
	for _, blocker := range detail.FinancingBlockers {
		resp.FinancingBlockers = append(resp.FinancingBlockers, dto.ListingFinancingBlockerResponse{
			ID:               blocker.Item.ID(),
			ListingID:        blocker.Item.ListingID(),
			ListingVersionID: blocker.Item.ListingVersionID(),
			Blocker:          catalogDetailToDTOWithFallback(blocker.Catalog, uint8(blocker.Item.Blocker())),
		})
	}

	resp.Guarantees = make([]dto.ListingGuaranteeResponse, 0, len(detail.Guarantees))
	for _, guarantee := range detail.Guarantees {
		resp.Guarantees = append(resp.Guarantees, dto.ListingGuaranteeResponse{
			ID:               guarantee.Item.ID(),
			ListingID:        guarantee.Item.ListingID(),
			ListingVersionID: guarantee.Item.ListingVersionID(),
			Priority:         guarantee.Item.Priority(),
			Guarantee:        catalogDetailToDTOWithFallback(guarantee.Catalog, uint8(guarantee.Item.Guarantee())),
		})
	}

	if draftVersion, ok := listing.DraftVersion(); ok && draftVersion != nil {
		if draftID := draftVersion.ID(); draftID > 0 {
			resp.DraftVersionID = &draftID
		}
	}

	if resp.ActiveVersionID == 0 {
		resp.ActiveVersionID = listing.ID()
	}

	return resp
}

func catalogDetailToPointer(detail *listingservices.CatalogValueDetail, fallback uint8) *dto.CatalogItemResponse {
	if detail != nil {
		return &dto.CatalogItemResponse{
			NumericValue: detail.NumericValue,
			Slug:         detail.Slug,
			Label:        detail.Label,
		}
	}

	if fallback == 0 {
		return nil
	}

	return &dto.CatalogItemResponse{NumericValue: fallback}
}

func catalogDetailToDTOWithFallback(detail *listingservices.CatalogValueDetail, fallback uint8) dto.CatalogItemResponse {
	if detail != nil {
		return dto.CatalogItemResponse{
			NumericValue: detail.NumericValue,
			Slug:         detail.Slug,
			Label:        detail.Label,
		}
	}

	return dto.CatalogItemResponse{
		NumericValue: fallback,
	}
}
