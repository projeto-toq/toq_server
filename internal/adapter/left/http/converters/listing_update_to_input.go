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
	var listingID int64
	if req.ID.IsPresent() && !req.ID.IsNull() {
		if id, ok := req.ID.Value(); ok {
			listingID = id
		}
	}

	input := listingservices.UpdateListingInput{ID: listingID}

	if req.Owner.IsPresent() {
		if req.Owner.IsNull() {
			input.Owner = coreutils.NewOptionalNull[listingmodel.PropertyOwner]()
		} else if value, ok := req.Owner.Value(); ok {
			input.Owner = coreutils.NewOptionalValue(listingmodel.PropertyOwner(value))
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
			input.Delivered = coreutils.NewOptionalNull[listingmodel.PropertyDelivered]()
		} else if value, ok := req.Delivered.Value(); ok {
			input.Delivered = coreutils.NewOptionalValue(listingmodel.PropertyDelivered(value))
		}
	}

	if req.WhoLives.IsPresent() {
		if req.WhoLives.IsNull() {
			input.WhoLives = coreutils.NewOptionalNull[listingmodel.WhoLives]()
		} else if value, ok := req.WhoLives.Value(); ok {
			input.WhoLives = coreutils.NewOptionalValue(listingmodel.WhoLives(value))
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
			input.Transaction = coreutils.NewOptionalNull[listingmodel.TransactionType]()
		} else if value, ok := req.Transaction.Value(); ok {
			input.Transaction = coreutils.NewOptionalValue(listingmodel.TransactionType(value))
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

	if req.AnnualGroundRent.IsPresent() {
		if req.AnnualGroundRent.IsNull() {
			input.AnnualGroundRent = coreutils.NewOptionalNull[float64]()
		} else if value, ok := req.AnnualGroundRent.Value(); ok {
			input.AnnualGroundRent = coreutils.NewOptionalValue(value)
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
			input.Installment = coreutils.NewOptionalNull[listingmodel.InstallmentPlan]()
		} else if value, ok := req.Installment.Value(); ok {
			input.Installment = coreutils.NewOptionalValue(listingmodel.InstallmentPlan(value))
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
			input.FinancingBlockers = coreutils.NewOptionalNull[[]listingmodel.FinancingBlockerInterface]()
		} else if blockersReq, ok := req.FinancingBlockers.Value(); ok {
			blockers := make([]listingmodel.FinancingBlockerInterface, 0, len(blockersReq))
			for _, blockerValue := range blockersReq {
				if blockerValue <= 0 {
					return input, fmt.Errorf("financingBlocker deve ser maior que zero")
				}
				blocker := listingmodel.NewFinancingBlocker()
				blocker.SetBlocker(listingmodel.FinancingBlocker(blockerValue))
				blockers = append(blockers, blocker)
			}
			input.FinancingBlockers = coreutils.NewOptionalValue(blockers)
		}
	}

	if req.Guarantees.IsPresent() {
		if req.Guarantees.IsNull() {
			input.Guarantees = coreutils.NewOptionalNull[[]listingmodel.GuaranteeInterface]()
		} else if guaranteesReq, ok := req.Guarantees.Value(); ok {
			guarantees := make([]listingmodel.GuaranteeInterface, 0, len(guaranteesReq))
			for _, guaranteeReq := range guaranteesReq {
				if guaranteeReq.Guarantee <= 0 {
					return input, fmt.Errorf("guarantee deve ser maior que zero")
				}
				guarantee := listingmodel.NewGuarantee()
				guarantee.SetPriority(guaranteeReq.Priority)
				guarantee.SetGuarantee(listingmodel.GuaranteeType(guaranteeReq.Guarantee))
				guarantees = append(guarantees, guarantee)
			}
			input.Guarantees = coreutils.NewOptionalValue(guarantees)
		}
	}

	if req.Visit.IsPresent() {
		if req.Visit.IsNull() {
			input.Visit = coreutils.NewOptionalNull[listingmodel.VisitType]()
		} else if value, ok := req.Visit.Value(); ok {
			input.Visit = coreutils.NewOptionalValue(listingmodel.VisitType(value))
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
			input.Accompanying = coreutils.NewOptionalNull[listingmodel.AccompanyingType]()
		} else if value, ok := req.Accompanying.Value(); ok {
			input.Accompanying = coreutils.NewOptionalValue(listingmodel.AccompanyingType(value))
		}
	}

	return input, nil
}
