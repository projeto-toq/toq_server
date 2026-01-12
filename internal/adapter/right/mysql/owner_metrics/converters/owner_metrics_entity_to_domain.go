package converters

import (
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/owner_metrics/entities"
	usermodel "github.com/projeto-toq/toq_server/internal/core/model/user_model"
)

// ToDomain converts a database entity into the domain aggregate.
func ToDomain(entity entities.OwnerMetricsEntity) usermodel.OwnerResponseMetrics {
	metrics := usermodel.NewOwnerResponseMetrics()
	metrics.SetOwnerID(entity.UserID)
	metrics.SetVisitAverageSeconds(entity.VisitAvgResponseSeconds)
	metrics.SetVisitResponsesTotal(entity.VisitTotalResponses)
	metrics.SetVisitLastResponseAt(entity.VisitLastResponseAt)
	metrics.SetProposalAverageSeconds(entity.ProposalAvgResponseSeconds)
	metrics.SetProposalResponsesTotal(entity.ProposalTotalResponses)
	metrics.SetProposalLastResponseAt(entity.ProposalLastResponseAt)
	return metrics
}
