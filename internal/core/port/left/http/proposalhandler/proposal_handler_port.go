package proposalhandler

import "github.com/gin-gonic/gin"

// Handler defines all HTTP entrypoints for the proposal domain.
type Handler interface {
	CreateProposal(c *gin.Context)
	UpdateProposal(c *gin.Context)
	CancelProposal(c *gin.Context)
	AcceptProposal(c *gin.Context)
	RejectProposal(c *gin.Context)
	ListRealtorProposals(c *gin.Context)
	ListOwnerProposals(c *gin.Context)
	GetProposalDetail(c *gin.Context)
}
