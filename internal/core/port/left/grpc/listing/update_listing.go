package grpclistingport

import (
	"context"

	"github.com/giulio-alfieri/toq_server/internal/adapter/left/grpc/pb"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	listingmodel "github.com/giulio-alfieri/toq_server/internal/core/model/listing_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (lr *ListingHandler) UpdateListing(ctx context.Context, in *pb.UpdateListingRequest) (out *pb.UpdateListingResponse, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return
	}
	defer spanEnd()

	infos := ctx.Value(globalmodel.TokenKey).(usermodel.UserInfos)

	//create auxiliares models
	features := lr.FeaturesPBToModel(ctx, in.Listing.Features, in.Listing.Id)
	places := lr.ExchangePlacesPBToModel(ctx, in.Listing.ExchangePlaces, in.Listing.Id)
	blockers := lr.FinancingBlockersPBToModel(ctx, in.Listing.FinancingBlockers, in.Listing.Id)
	guarantees := lr.GuaranteesPBToModel(ctx, in.Listing.Guarantees, in.Listing.Id)

	listing := listingmodel.NewListing()

	listing.SetID(in.Listing.Id)
	listing.SetUserID(infos.ID)
	listing.SetOwner(listingmodel.PropertyOwner(in.Listing.Owner))
	listing.SetFeatures(features)
	listing.SetLandSize(float64(in.Listing.Landsize))
	listing.SetCorner(in.Listing.Corner)
	listing.SetNonBuildable(float64(in.Listing.NonBuildable))
	listing.SetBuildable(float64(in.Listing.Buildable))
	listing.SetDelivered(listingmodel.PropertyDelivered(in.Listing.Delivered))
	listing.SetWhoLives(listingmodel.WhoLives(in.Listing.WhoLives))
	listing.SetDescription(in.Listing.Description)
	listing.SetTransaction(listingmodel.TransactionType(in.Listing.Transaction))
	listing.SetSellNet(float64(in.Listing.SellNet))
	listing.SetRentNet(float64(in.Listing.RentNet))
	listing.SetCondominium(float64(in.Listing.Condominium))
	listing.SetAnnualTax(float64(in.Listing.AnnualTax))
	listing.SetAnnualGroundRent(float64(in.Listing.AnnualGroundRent))
	listing.SetExchange(in.Listing.Exchange)
	listing.SetExchangePercentual(float64(in.Listing.ExchangePercentual))
	listing.SetExchangePlaces(places)
	listing.SetInstallment(listingmodel.InstallmentPlan(in.Listing.Installment))
	listing.SetFinancing(in.Listing.Financing)
	listing.SetFinancingBlockers(blockers)
	listing.SetGuarantees(guarantees)
	listing.SetVisit(listingmodel.VisitType(in.Listing.Visit))
	listing.SetTenantName(in.Listing.TenantName)
	listing.SetTenantEmail(in.Listing.TenantEmail)
	listing.SetTenantPhone(in.Listing.TenantPhone)
	listing.SetAccompanying(listingmodel.AccompanyingType(in.Listing.Accompanying))

	err = lr.service.UpdateListing(ctx, listing)
	if err != nil {
		return

	}

	return
}
