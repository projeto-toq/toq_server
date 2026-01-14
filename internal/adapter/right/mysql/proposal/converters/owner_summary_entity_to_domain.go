package converters

import (
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/entities"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
)

// ToOwnerSummaryModel converts an OwnerSummaryEntity into the domain representation.
func ToOwnerSummaryModel(entity entities.OwnerSummaryEntity) proposalmodel.OwnerSummary {
	summary := proposalmodel.NewOwnerSummary()
	summary.SetID(entity.OwnerID)
	summary.SetFullName(entity.FullName)
	if entity.MemberSinceMonths.Valid {
		summary.SetMemberSinceMonths(int(entity.MemberSinceMonths.Int64))
	}
	summary.SetProposalAvgSeconds(entity.ProposalAvgSeconds)
	summary.SetVisitAvgSeconds(entity.VisitAvgSeconds)
	return summary
}
