package converters

import (
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/adapter/left/http/dto"
	listingmodel "github.com/projeto-toq/toq_server/internal/core/model/listing_model"
	listingservices "github.com/projeto-toq/toq_server/internal/core/service/listing_service"
	coreutils "github.com/projeto-toq/toq_server/internal/core/utils"
)

// UpdateListingRequestToInput converte o DTO de update em um input para o service,
// mantendo a distinção entre campos omitidos, vazios e nulos.
func UpdateListingRequestToInput(req dto.UpdateListingRequest) (listingservices.UpdateListingInput, error) {
	input := listingservices.UpdateListingInput{}

	if req.Owner.IsPresent() {
		if req.Owner.IsNull() {
			input.Owner = coreutils.NewOptionalNull[listingservices.CatalogSelection]()
		} else if value, ok := req.Owner.Value(); ok {
			slug := strings.TrimSpace(value)
			if slug == "" {
				return input, fmt.Errorf("owner slug deve ser informado")
			}
			input.Owner = coreutils.NewOptionalValue(listingservices.NewCatalogSelectionFromSlug(slug))
		}
	}

	if req.Features.IsPresent() {
		if req.Features.IsNull() {
			input.Features = coreutils.NewOptionalNull[[]listingmodel.FeatureInterface]()
		} else if featuresReq, ok := req.Features.Value(); ok {
			features := make([]listingmodel.FeatureInterface, 0, len(featuresReq))
			for _, item := range featuresReq {
				if item.FeatureID <= 0 {
					return input, fmt.Errorf("featureId deve ser maior que zero")
				}
				feature := listingmodel.NewFeature()
				feature.SetFeatureID(item.FeatureID)
				feature.SetQuantity(item.Quantity)
				features = append(features, feature)
			}
			input.Features = coreutils.NewOptionalValue(features)
		}
	}

	if req.LandSize.IsPresent() {
		if req.LandSize.IsNull() {
			input.LandSize = coreutils.NewOptionalNull[float64]()
		} else if value, ok := req.LandSize.Value(); ok {
			input.LandSize = coreutils.NewOptionalValue(value)
		}
	}

	if req.Corner.IsPresent() {
		if req.Corner.IsNull() {
			input.Corner = coreutils.NewOptionalNull[bool]()
		} else if value, ok := req.Corner.Value(); ok {
			input.Corner = coreutils.NewOptionalValue(value)
		}
	}

	if req.NonBuildable.IsPresent() {
		if req.NonBuildable.IsNull() {
			input.NonBuildable = coreutils.NewOptionalNull[float64]()
		} else if value, ok := req.NonBuildable.Value(); ok {
			input.NonBuildable = coreutils.NewOptionalValue(value)
		}
	}

	if req.Buildable.IsPresent() {
		if req.Buildable.IsNull() {
			input.Buildable = coreutils.NewOptionalNull[float64]()
		} else if value, ok := req.Buildable.Value(); ok {
			input.Buildable = coreutils.NewOptionalValue(value)
		}
	}

	if req.Delivered.IsPresent() {
		if req.Delivered.IsNull() {
			input.Delivered = coreutils.NewOptionalNull[listingservices.CatalogSelection]()
		} else if value, ok := req.Delivered.Value(); ok {
			slug := strings.TrimSpace(value)
			if slug == "" {
				return input, fmt.Errorf("delivered slug deve ser informado")
			}
			input.Delivered = coreutils.NewOptionalValue(listingservices.NewCatalogSelectionFromSlug(slug))
		}
	}

	if req.WhoLives.IsPresent() {
		if req.WhoLives.IsNull() {
			input.WhoLives = coreutils.NewOptionalNull[listingservices.CatalogSelection]()
		} else if value, ok := req.WhoLives.Value(); ok {
			slug := strings.TrimSpace(value)
			if slug == "" {
				return input, fmt.Errorf("whoLives slug deve ser informado")
			}
			input.WhoLives = coreutils.NewOptionalValue(listingservices.NewCatalogSelectionFromSlug(slug))
		}
	}

	if req.Complex.IsPresent() {
		if req.Complex.IsNull() {
			input.Complex = coreutils.NewOptionalNull[string]()
		} else if value, ok := req.Complex.Value(); ok {
			input.Complex = coreutils.NewOptionalValue(strings.TrimSpace(value))
		}
	}

	if req.Title.IsPresent() {
		if req.Title.IsNull() {
			input.Title = coreutils.NewOptionalNull[string]()
		} else if value, ok := req.Title.Value(); ok {
			trimmed := strings.TrimSpace(value)
			input.Title = coreutils.NewOptionalValue(trimmed)
		}
	}

	if req.Description.IsPresent() {
		if req.Description.IsNull() {
			input.Description = coreutils.NewOptionalNull[string]()
		} else if value, ok := req.Description.Value(); ok {
			input.Description = coreutils.NewOptionalValue(value)
		}
	}

	if req.Transaction.IsPresent() {
		if req.Transaction.IsNull() {
			input.Transaction = coreutils.NewOptionalNull[listingservices.CatalogSelection]()
		} else if value, ok := req.Transaction.Value(); ok {
			slug := strings.TrimSpace(value)
			if slug == "" {
				return input, fmt.Errorf("transaction slug deve ser informado")
			}
			input.Transaction = coreutils.NewOptionalValue(listingservices.NewCatalogSelectionFromSlug(slug))
		}
	}

	if req.SellNet.IsPresent() {
		if req.SellNet.IsNull() {
			input.SellNet = coreutils.NewOptionalNull[float64]()
		} else if value, ok := req.SellNet.Value(); ok {
			input.SellNet = coreutils.NewOptionalValue(value)
		}
	}

	if req.RentNet.IsPresent() {
		if req.RentNet.IsNull() {
			input.RentNet = coreutils.NewOptionalNull[float64]()
		} else if value, ok := req.RentNet.Value(); ok {
			input.RentNet = coreutils.NewOptionalValue(value)
		}
	}

	if req.Condominium.IsPresent() {
		if req.Condominium.IsNull() {
			input.Condominium = coreutils.NewOptionalNull[float64]()
		} else if value, ok := req.Condominium.Value(); ok {
			input.Condominium = coreutils.NewOptionalValue(value)
		}
	}

	if req.AnnualTax.IsPresent() {
		if req.AnnualTax.IsNull() {
			input.AnnualTax = coreutils.NewOptionalNull[float64]()
		} else if value, ok := req.AnnualTax.Value(); ok {
			input.AnnualTax = coreutils.NewOptionalValue(value)
		}
	}

	if req.MonthlyTax.IsPresent() {
		if req.MonthlyTax.IsNull() {
			input.MonthlyTax = coreutils.NewOptionalNull[float64]()
		} else if value, ok := req.MonthlyTax.Value(); ok {
			input.MonthlyTax = coreutils.NewOptionalValue(value)
		}
	}

	if req.AnnualGroundRent.IsPresent() {
		if req.AnnualGroundRent.IsNull() {
			input.AnnualGroundRent = coreutils.NewOptionalNull[float64]()
		} else if value, ok := req.AnnualGroundRent.Value(); ok {
			input.AnnualGroundRent = coreutils.NewOptionalValue(value)
		}
	}

	if req.MonthlyGroundRent.IsPresent() {
		if req.MonthlyGroundRent.IsNull() {
			input.MonthlyGroundRent = coreutils.NewOptionalNull[float64]()
		} else if value, ok := req.MonthlyGroundRent.Value(); ok {
			input.MonthlyGroundRent = coreutils.NewOptionalValue(value)
		}
	}

	if req.Exchange.IsPresent() {
		if req.Exchange.IsNull() {
			input.Exchange = coreutils.NewOptionalNull[bool]()
		} else if value, ok := req.Exchange.Value(); ok {
			input.Exchange = coreutils.NewOptionalValue(value)
		}
	}

	if req.ExchangePercentual.IsPresent() {
		if req.ExchangePercentual.IsNull() {
			input.ExchangePercentual = coreutils.NewOptionalNull[float64]()
		} else if value, ok := req.ExchangePercentual.Value(); ok {
			input.ExchangePercentual = coreutils.NewOptionalValue(value)
		}
	}

	if req.ExchangePlaces.IsPresent() {
		if req.ExchangePlaces.IsNull() {
			input.ExchangePlaces = coreutils.NewOptionalNull[[]listingmodel.ExchangePlaceInterface]()
		} else if placesReq, ok := req.ExchangePlaces.Value(); ok {
			places := make([]listingmodel.ExchangePlaceInterface, 0, len(placesReq))
			for _, placeReq := range placesReq {
				place := listingmodel.NewExchangePlace()
				place.SetNeighborhood(strings.TrimSpace(placeReq.Neighborhood))
				place.SetCity(strings.TrimSpace(placeReq.City))
				place.SetState(strings.TrimSpace(placeReq.State))
				places = append(places, place)
			}
			input.ExchangePlaces = coreutils.NewOptionalValue(places)
		}
	}

	if req.Installment.IsPresent() {
		if req.Installment.IsNull() {
			input.Installment = coreutils.NewOptionalNull[listingservices.CatalogSelection]()
		} else if value, ok := req.Installment.Value(); ok {
			slug := strings.TrimSpace(value)
			if slug == "" {
				return input, fmt.Errorf("installment slug deve ser informado")
			}
			input.Installment = coreutils.NewOptionalValue(listingservices.NewCatalogSelectionFromSlug(slug))
		}
	}

	if req.Financing.IsPresent() {
		if req.Financing.IsNull() {
			input.Financing = coreutils.NewOptionalNull[bool]()
		} else if value, ok := req.Financing.Value(); ok {
			input.Financing = coreutils.NewOptionalValue(value)
		}
	}

	if req.FinancingBlockers.IsPresent() {
		if req.FinancingBlockers.IsNull() {
			input.FinancingBlockers = coreutils.NewOptionalNull[[]listingservices.CatalogSelection]()
		} else if blockersReq, ok := req.FinancingBlockers.Value(); ok {
			blockers := make([]listingservices.CatalogSelection, 0, len(blockersReq))
			for _, blockerSlug := range blockersReq {
				slug := strings.TrimSpace(blockerSlug)
				if slug == "" {
					return input, fmt.Errorf("financingBlocker slug deve ser informado")
				}
				blockers = append(blockers, listingservices.NewCatalogSelectionFromSlug(slug))
			}
			input.FinancingBlockers = coreutils.NewOptionalValue(blockers)
		}
	}

	if req.Guarantees.IsPresent() {
		if req.Guarantees.IsNull() {
			input.Guarantees = coreutils.NewOptionalNull[[]listingservices.GuaranteeUpdate]()
		} else if guaranteesReq, ok := req.Guarantees.Value(); ok {
			guarantees := make([]listingservices.GuaranteeUpdate, 0, len(guaranteesReq))
			for _, guaranteeReq := range guaranteesReq {
				slug := strings.TrimSpace(guaranteeReq.Guarantee)
				if slug == "" {
					return input, fmt.Errorf("guarantee slug deve ser informado")
				}
				guarantees = append(guarantees, listingservices.GuaranteeUpdate{
					Priority:  guaranteeReq.Priority,
					Selection: listingservices.NewCatalogSelectionFromSlug(slug),
				})
			}
			input.Guarantees = coreutils.NewOptionalValue(guarantees)
		}
	}

	if req.Visit.IsPresent() {
		if req.Visit.IsNull() {
			input.Visit = coreutils.NewOptionalNull[listingservices.CatalogSelection]()
		} else if value, ok := req.Visit.Value(); ok {
			slug := strings.TrimSpace(value)
			if slug == "" {
				return input, fmt.Errorf("visit slug deve ser informado")
			}
			input.Visit = coreutils.NewOptionalValue(listingservices.NewCatalogSelectionFromSlug(slug))
		}
	}

	if req.TenantName.IsPresent() {
		if req.TenantName.IsNull() {
			input.TenantName = coreutils.NewOptionalNull[string]()
		} else if value, ok := req.TenantName.Value(); ok {
			input.TenantName = coreutils.NewOptionalValue(strings.TrimSpace(value))
		}
	}

	if req.TenantEmail.IsPresent() {
		if req.TenantEmail.IsNull() {
			input.TenantEmail = coreutils.NewOptionalNull[string]()
		} else if value, ok := req.TenantEmail.Value(); ok {
			input.TenantEmail = coreutils.NewOptionalValue(strings.TrimSpace(value))
		}
	}

	if req.TenantPhone.IsPresent() {
		if req.TenantPhone.IsNull() {
			input.TenantPhone = coreutils.NewOptionalNull[string]()
		} else if value, ok := req.TenantPhone.Value(); ok {
			input.TenantPhone = coreutils.NewOptionalValue(strings.TrimSpace(value))
		}
	}

	if req.Accompanying.IsPresent() {
		if req.Accompanying.IsNull() {
			input.Accompanying = coreutils.NewOptionalNull[listingservices.CatalogSelection]()
		} else if value, ok := req.Accompanying.Value(); ok {
			slug := strings.TrimSpace(value)
			if slug == "" {
				return input, fmt.Errorf("accompanying slug deve ser informado")
			}
			input.Accompanying = coreutils.NewOptionalValue(listingservices.NewCatalogSelectionFromSlug(slug))
		}
	}

	if req.CompletionForecast.IsPresent() {
		if req.CompletionForecast.IsNull() {
			input.CompletionForecast = coreutils.NewOptionalNull[string]()
		} else if value, ok := req.CompletionForecast.Value(); ok {
			input.CompletionForecast = coreutils.NewOptionalValue(strings.TrimSpace(value))
		}
	}

	if req.LandBlock.IsPresent() {
		if req.LandBlock.IsNull() {
			input.LandBlock = coreutils.NewOptionalNull[string]()
		} else if value, ok := req.LandBlock.Value(); ok {
			input.LandBlock = coreutils.NewOptionalValue(strings.TrimSpace(value))
		}
	}

	if req.LandLot.IsPresent() {
		if req.LandLot.IsNull() {
			input.LandLot = coreutils.NewOptionalNull[string]()
		} else if value, ok := req.LandLot.Value(); ok {
			input.LandLot = coreutils.NewOptionalValue(strings.TrimSpace(value))
		}
	}

	if req.LandFront.IsPresent() {
		if req.LandFront.IsNull() {
			input.LandFront = coreutils.NewOptionalNull[float64]()
		} else if value, ok := req.LandFront.Value(); ok {
			input.LandFront = coreutils.NewOptionalValue(value)
		}
	}

	if req.LandSide.IsPresent() {
		if req.LandSide.IsNull() {
			input.LandSide = coreutils.NewOptionalNull[float64]()
		} else if value, ok := req.LandSide.Value(); ok {
			input.LandSide = coreutils.NewOptionalValue(value)
		}
	}

	if req.LandBack.IsPresent() {
		if req.LandBack.IsNull() {
			input.LandBack = coreutils.NewOptionalNull[float64]()
		} else if value, ok := req.LandBack.Value(); ok {
			input.LandBack = coreutils.NewOptionalValue(value)
		}
	}

	if req.LandTerrainType.IsPresent() {
		if req.LandTerrainType.IsNull() {
			input.LandTerrainType = coreutils.NewOptionalNull[listingservices.CatalogSelection]()
		} else if value, ok := req.LandTerrainType.Value(); ok {
			slug := strings.TrimSpace(value)
			if slug == "" {
				return input, fmt.Errorf("landTerrainType slug deve ser informado")
			}
			input.LandTerrainType = coreutils.NewOptionalValue(listingservices.NewCatalogSelectionFromSlug(slug))
		}
	}

	if req.HasKmz.IsPresent() {
		if req.HasKmz.IsNull() {
			input.HasKmz = coreutils.NewOptionalNull[bool]()
		} else if value, ok := req.HasKmz.Value(); ok {
			input.HasKmz = coreutils.NewOptionalValue(value)
		}
	}

	if req.KmzFile.IsPresent() {
		if req.KmzFile.IsNull() {
			input.KmzFile = coreutils.NewOptionalNull[string]()
		} else if value, ok := req.KmzFile.Value(); ok {
			input.KmzFile = coreutils.NewOptionalValue(strings.TrimSpace(value))
		}
	}

	if req.BuildingFloors.IsPresent() {
		if req.BuildingFloors.IsNull() {
			input.BuildingFloors = coreutils.NewOptionalNull[int16]()
		} else if value, ok := req.BuildingFloors.Value(); ok {
			input.BuildingFloors = coreutils.NewOptionalValue(value)
		}
	}

	if req.UnitTower.IsPresent() {
		if req.UnitTower.IsNull() {
			input.UnitTower = coreutils.NewOptionalNull[string]()
		} else if value, ok := req.UnitTower.Value(); ok {
			input.UnitTower = coreutils.NewOptionalValue(strings.TrimSpace(value))
		}
	}

	if req.UnitFloor.IsPresent() {
		if req.UnitFloor.IsNull() {
			input.UnitFloor = coreutils.NewOptionalNull[int16]()
		} else if value, ok := req.UnitFloor.Value(); ok {
			input.UnitFloor = coreutils.NewOptionalValue(value)
		}
	}

	if req.UnitNumber.IsPresent() {
		if req.UnitNumber.IsNull() {
			input.UnitNumber = coreutils.NewOptionalNull[string]()
		} else if value, ok := req.UnitNumber.Value(); ok {
			input.UnitNumber = coreutils.NewOptionalValue(strings.TrimSpace(value))
		}
	}

	if req.WarehouseManufacturingArea.IsPresent() {
		if req.WarehouseManufacturingArea.IsNull() {
			input.WarehouseManufacturingArea = coreutils.NewOptionalNull[float64]()
		} else if value, ok := req.WarehouseManufacturingArea.Value(); ok {
			input.WarehouseManufacturingArea = coreutils.NewOptionalValue(value)
		}
	}

	if req.WarehouseSector.IsPresent() {
		if req.WarehouseSector.IsNull() {
			input.WarehouseSector = coreutils.NewOptionalNull[listingservices.CatalogSelection]()
		} else if value, ok := req.WarehouseSector.Value(); ok {
			slug := strings.TrimSpace(value)
			if slug == "" {
				return input, fmt.Errorf("warehouseSector slug deve ser informado")
			}
			input.WarehouseSector = coreutils.NewOptionalValue(listingservices.NewCatalogSelectionFromSlug(slug))
		}
	}

	if req.WarehouseHasPrimaryCabin.IsPresent() {
		if req.WarehouseHasPrimaryCabin.IsNull() {
			input.WarehouseHasPrimaryCabin = coreutils.NewOptionalNull[bool]()
		} else if value, ok := req.WarehouseHasPrimaryCabin.Value(); ok {
			input.WarehouseHasPrimaryCabin = coreutils.NewOptionalValue(value)
		}
	}

	if req.WarehouseCabinKva.IsPresent() {
		if req.WarehouseCabinKva.IsNull() {
			input.WarehouseCabinKva = coreutils.NewOptionalNull[float64]()
		} else if value, ok := req.WarehouseCabinKva.Value(); ok {
			input.WarehouseCabinKva = coreutils.NewOptionalValue(value)
		}
	}

	if req.WarehouseGroundFloor.IsPresent() {
		if req.WarehouseGroundFloor.IsNull() {
			input.WarehouseGroundFloor = coreutils.NewOptionalNull[float64]()
		} else if value, ok := req.WarehouseGroundFloor.Value(); ok {
			input.WarehouseGroundFloor = coreutils.NewOptionalValue(value)
		}
	}

	if req.WarehouseFloorResistance.IsPresent() {
		if req.WarehouseFloorResistance.IsNull() {
			input.WarehouseFloorResistance = coreutils.NewOptionalNull[float64]()
		} else if value, ok := req.WarehouseFloorResistance.Value(); ok {
			input.WarehouseFloorResistance = coreutils.NewOptionalValue(value)
		}
	}

	if req.WarehouseZoning.IsPresent() {
		if req.WarehouseZoning.IsNull() {
			input.WarehouseZoning = coreutils.NewOptionalNull[string]()
		} else if value, ok := req.WarehouseZoning.Value(); ok {
			input.WarehouseZoning = coreutils.NewOptionalValue(strings.TrimSpace(value))
		}
	}

	if req.WarehouseHasOfficeArea.IsPresent() {
		if req.WarehouseHasOfficeArea.IsNull() {
			input.WarehouseHasOfficeArea = coreutils.NewOptionalNull[bool]()
		} else if value, ok := req.WarehouseHasOfficeArea.Value(); ok {
			input.WarehouseHasOfficeArea = coreutils.NewOptionalValue(value)
		}
	}

	if req.WarehouseOfficeArea.IsPresent() {
		if req.WarehouseOfficeArea.IsNull() {
			input.WarehouseOfficeArea = coreutils.NewOptionalNull[float64]()
		} else if value, ok := req.WarehouseOfficeArea.Value(); ok {
			input.WarehouseOfficeArea = coreutils.NewOptionalValue(value)
		}
	}

	if req.WarehouseAdditionalFloors.IsPresent() {
		if req.WarehouseAdditionalFloors.IsNull() {
			input.WarehouseAdditionalFloors = coreutils.NewOptionalNull[[]listingmodel.WarehouseAdditionalFloorInterface]()
		} else if floorsReq, ok := req.WarehouseAdditionalFloors.Value(); ok {
			floors := make([]listingmodel.WarehouseAdditionalFloorInterface, 0, len(floorsReq))
			for _, floorReq := range floorsReq {
				floor := listingmodel.NewWarehouseAdditionalFloor()
				floor.SetFloorName(strings.TrimSpace(floorReq.FloorName))
				floor.SetFloorOrder(floorReq.FloorOrder)
				floor.SetFloorHeight(floorReq.FloorHeight)
				floors = append(floors, floor)
			}
			input.WarehouseAdditionalFloors = coreutils.NewOptionalValue(floors)
		}
	}

	if req.StoreHasMezzanine.IsPresent() {
		if req.StoreHasMezzanine.IsNull() {
			input.StoreHasMezzanine = coreutils.NewOptionalNull[bool]()
		} else if value, ok := req.StoreHasMezzanine.Value(); ok {
			input.StoreHasMezzanine = coreutils.NewOptionalValue(value)
		}
	}

	if req.StoreMezzanineArea.IsPresent() {
		if req.StoreMezzanineArea.IsNull() {
			input.StoreMezzanineArea = coreutils.NewOptionalNull[float64]()
		} else if value, ok := req.StoreMezzanineArea.Value(); ok {
			input.StoreMezzanineArea = coreutils.NewOptionalValue(value)
		}
	}

	return input, nil
}
