package mysqlvisitadapter

import "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/visit/entity"

type rowScanner interface {
	Scan(dest ...any) error
}

func scanVisitEntity(scanner rowScanner) (entity.VisitEntity, error) {
	var visit entity.VisitEntity
	if err := scanner.Scan(
		&visit.ID,
		&visit.ListingID,
		&visit.OwnerID,
		&visit.RealtorID,
		&visit.ScheduledStart,
		&visit.ScheduledEnd,
		&visit.Status,
		&visit.CancelReason,
		&visit.Notes,
		&visit.CreatedBy,
		&visit.UpdatedBy,
	); err != nil {
		return entity.VisitEntity{}, err
	}

	return visit, nil
}
