package complexservices

import (
	"context"
	"database/sql"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (cs *complexService) GetOptions(ctx context.Context, zipCode string, number string) (propertyTypes globalmodel.PropertyType, err error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return propertyTypes, err
	}
	defer spanEnd()

	tx, err := cs.gsi.StartTransaction(ctx)
	if err != nil {
		return
	}

	propertyTypes, err = cs.getOptions(ctx, tx, zipCode, number)
	if err != nil {
		cs.gsi.RollbackTransaction(ctx, tx)
		return
	}

	err = cs.gsi.CommitTransaction(ctx, tx)
	if err != nil {
		cs.gsi.RollbackTransaction(ctx, tx)
		return
	}

	return
}

func (cs *complexService) getOptions(ctx context.Context, tx *sql.Tx, zipCode string, number string) (propertyTypes globalmodel.PropertyType, err error) {
	callhorizontal := false
	complex, err := cs.complexRepository.GetVerticalComplex(ctx, tx, zipCode, number)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			callhorizontal = true
		} else {
			return
		}
	}
	if callhorizontal {
		complex, err = cs.complexRepository.GetHorizontalComplex(ctx, tx, zipCode)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return propertyTypes, status.Error(codes.InvalidArgument, "Area not covered yet")
			} else {
				return
			}
		}
	}

	propertyTypes = complex.GetPropertyType()
	return

}
