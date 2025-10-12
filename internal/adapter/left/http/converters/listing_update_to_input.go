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
			owner := listingmodel.PropertyOwner(value)
			if err := validatePropertyOwner(owner); err != nil {
				return input, err
			}
			input.Owner = coreutils.NewOptionalValue(owner)
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
			delivered := listingmodel.PropertyDelivered(value)
			if err := validateDelivered(delivered); err != nil {
				return input, err
			}
			input.Delivered = coreutils.NewOptionalValue(delivered)
		}
	}

	if req.WhoLives.IsPresent() {
		if req.WhoLives.IsNull() {
			input.WhoLives = coreutils.NewOptionalNull[listingmodel.WhoLives]()
		} else if value, ok := req.WhoLives.Value(); ok {
			wl := listingmodel.WhoLives(value)
			if err := validateWhoLives(wl); err != nil {
				return input, err
			}
			input.WhoLives = coreutils.NewOptionalValue(wl)
		}
	}

	if req.Description.IsPresent() {
		if req.Description.IsNull() {
			input.Description = coreutils.NewOptionalNull[string]()
		} else if value, ok := req.Description.Value(); ok {
			input.Description = coreutils.NewOptionalValue(strings.TrimSpace(value))
		}
	}

	if req.Transaction.IsPresent() {
		if req.Transaction.IsNull() {
			input.Transaction = coreutils.NewOptionalNull[listingmodel.TransactionType]()
		} else if value, ok := req.Transaction.Value(); ok {
			transaction := listingmodel.TransactionType(value)
			if err := validateTransaction(transaction); err != nil {
				return input, err
			}
			input.Transaction = coreutils.NewOptionalValue(transaction)
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
			plan := listingmodel.InstallmentPlan(value)
			if err := validateInstallment(plan); err != nil {
				return input, err
			}
			input.Installment = coreutils.NewOptionalValue(plan)
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
				blockerEnum := listingmodel.FinancingBlocker(blockerValue)
				if err := validateFinancingBlocker(blockerEnum); err != nil {
					return input, err
				}
				blocker := listingmodel.NewFinancingBlocker()
				blocker.SetBlocker(blockerEnum)
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
				guaranteeEnum := listingmodel.GuaranteeType(guaranteeReq.Guarantee)
				if err := validateGuarantee(guaranteeEnum); err != nil {
					return input, err
				}
				guarantee := listingmodel.NewGuarantee()
				guarantee.SetPriority(guaranteeReq.Priority)
				guarantee.SetGuarantee(guaranteeEnum)
				guarantees = append(guarantees, guarantee)
			}
			input.Guarantees = coreutils.NewOptionalValue(guarantees)
		}
	}

	if req.Visit.IsPresent() {
		if req.Visit.IsNull() {
			input.Visit = coreutils.NewOptionalNull[listingmodel.VisitType]()
		} else if value, ok := req.Visit.Value(); ok {
			visit := listingmodel.VisitType(value)
			if err := validateVisit(visit); err != nil {
				return input, err
			}
			input.Visit = coreutils.NewOptionalValue(visit)
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
			accompanying := listingmodel.AccompanyingType(value)
			if err := validateAccompanying(accompanying); err != nil {
				return input, err
			}
			input.Accompanying = coreutils.NewOptionalValue(accompanying)
		}
	}

	return input, nil
}

func validatePropertyOwner(owner listingmodel.PropertyOwner) error {
	if owner == 0 {
		return nil
	}
	if owner < listingmodel.OwnerMyself || owner > listingmodel.OwnerSiblings {
		return fmt.Errorf("owner inválido")
	}
	return nil
}

func validateDelivered(delivered listingmodel.PropertyDelivered) error {
	if delivered == 0 {
		return nil
	}
	if delivered < listingmodel.DeliveredFurnishedDecorated || delivered > listingmodel.DeliverdAsPictured {
		return fmt.Errorf("delivered inválido")
	}
	return nil
}

func validateWhoLives(who listingmodel.WhoLives) error {
	if who == 0 {
		return nil
	}
	if who < listingmodel.LivesOwner || who > listingmodel.LivesVacant {
		return fmt.Errorf("whoLives inválido")
	}
	return nil
}

func validateTransaction(tx listingmodel.TransactionType) error {
	if tx == 0 {
		return nil
	}
	if tx < listingmodel.TransactionSale || tx > listingmodel.TransactionBoth {
		return fmt.Errorf("transaction inválido")
	}
	return nil
}

func validateInstallment(plan listingmodel.InstallmentPlan) error {
	if plan == 0 {
		return nil
	}
	if plan < listingmodel.PlanCash || plan > listingmodel.PlanLongTerm {
		return fmt.Errorf("installment inválido")
	}
	return nil
}

func validateVisit(visit listingmodel.VisitType) error {
	if visit == 0 {
		return nil
	}
	if visit < listingmodel.VisitClient || visit > listingmodel.VisitAll {
		return fmt.Errorf("visit inválido")
	}
	return nil
}

func validateAccompanying(accompanying listingmodel.AccompanyingType) error {
	if accompanying == 0 {
		return nil
	}
	if accompanying < listingmodel.AccompanyingOwner || accompanying > listingmodel.AccompanyingAlone {
		return fmt.Errorf("accompanying inválido")
	}
	return nil
}

func validateFinancingBlocker(blocker listingmodel.FinancingBlocker) error {
	if blocker == 0 {
		return nil
	}
	if blocker < listingmodel.BlockerPendingProbate || blocker > listingmodel.Blockerother {
		return fmt.Errorf("financingBlocker inválido")
	}
	return nil
}

func validateGuarantee(guarantee listingmodel.GuaranteeType) error {
	if guarantee == 0 {
		return nil
	}
	if guarantee < listingmodel.GuaranteeDeposit || guarantee > listingmodel.GuaranteeRentalBond {
		return fmt.Errorf("guarantee inválido")
	}
	return nil
}
