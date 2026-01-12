package converters

import (
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/entities"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
)

// ToRealtorSummaryModel converts a RealtorSummaryEntity into the domain representation.
func ToRealtorSummaryModel(entity entities.RealtorSummaryEntity) proposalmodel.RealtorSummary {
	summary := proposalmodel.NewRealtorSummary()
	summary.SetID(entity.RealtorID)
	summary.SetName(entity.FullName)
	summary.SetNickname(entity.NickName)
	if entity.UsageMonths.Valid {
		summary.SetUsageMonths(int(entity.UsageMonths.Int64))
	}
	if entity.ProposalsCount.Valid {
		summary.SetProposalsCreated(entity.ProposalsCount.Int64)
	}
	if entity.AcceptedProposals.Valid {
		summary.SetAcceptedProposals(entity.AcceptedProposals.Int64)
	}
	return summary
}
