package converters

import (
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/entities"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
)

// ToProposalModel converts a ProposalEntity into the domain ProposalInterface.
func ToProposalModel(entity entities.ProposalEntity) proposalmodel.ProposalInterface {
	model := proposalmodel.NewProposal()
	model.SetID(entity.ID)
	model.SetListingIdentityID(entity.ListingIdentityID)
	model.SetRealtorID(entity.RealtorID)
	model.SetOwnerID(entity.OwnerID)
	model.SetProposalText(entity.ProposalText.String)
	model.SetRejectionReason(entity.RejectionReason)
	model.SetStatus(proposalmodel.Status(entity.Status))
	model.SetAcceptedAt(entity.AcceptedAt)
	model.SetRejectedAt(entity.RejectedAt)
	model.SetCancelledAt(entity.CancelledAt)
	if entity.DocumentsCount.Valid {
		model.SetDocumentsCount(int(entity.DocumentsCount.Int64))
	}
	return model
}

// ToProposalDocumentModel converts a ProposalDocumentEntity into a ProposalDocumentInterface.
func ToProposalDocumentModel(entity entities.ProposalDocumentEntity, includeBlob bool) proposalmodel.ProposalDocumentInterface {
	doc := proposalmodel.NewProposalDocument()
	doc.SetID(entity.ID)
	doc.SetProposalID(entity.ProposalID)
	doc.SetFileName(entity.FileName)
	doc.SetMimeType(entity.MimeType)
	doc.SetFileSizeBytes(entity.FileSizeBytes)
	if includeBlob {
		doc.SetFileData(entity.FileBlob)
	}
	doc.SetUploadedAt(entity.UploadedAt)
	return doc
}
