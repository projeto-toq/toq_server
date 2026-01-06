package proposalhandlers

import proposalservice "github.com/projeto-toq/toq_server/internal/core/service/proposal_service"

// ProposalHandler wires DTO conversion, authentication context and the service port.
type ProposalHandler struct {
	proposalService proposalservice.Service
}

// NewProposalHandler builds a handler with its dependencies injected by the factory.
func NewProposalHandler(service proposalservice.Service) *ProposalHandler {
	return &ProposalHandler{proposalService: service}
}
