package globalconverters

import (
	"context"

	globalentities "github.com/projeto-toq/toq_server/internal/adapter/right/mysql/global/entities"
	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
)

func AuditDomainToEntity(ctx context.Context, audit globalmodel.AuditInterface) (entity globalentities.AuditEntity) {
	entity = globalentities.AuditEntity{}

	entity.ID = audit.ID()
	entity.ExecutedAT = audit.ExecutedAt()
	entity.ExecutedBY = audit.ExecutedBy()
	entity.TableName = audit.TableName()
	entity.TableID = audit.TableID()
	entity.Action = audit.Action()

	return
}
