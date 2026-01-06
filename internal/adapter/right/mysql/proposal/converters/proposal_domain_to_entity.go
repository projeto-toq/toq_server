package converters

import (
	"database/sql"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/entities"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
)

// ToProposalEntity converts domain proposals into persistence entities.
func ToProposalEntity(model proposalmodel.ProposalInterface) entities.ProposalEntity {
	return entities.ProposalEntity{
		ID:                model.ID(),
		ListingIdentityID: model.ListingIdentityID(),
		RealtorID:         model.RealtorID(),
		OwnerID:           model.OwnerID(),
		ProposalText: sql.NullString{
			String: model.ProposalText(),
			Valid:  model.ProposalText() != "",
		},
		RejectionReason: model.RejectionReason(),
		Status:          string(model.Status()),
		AcceptedAt:      model.AcceptedAt(),
		RejectedAt:      model.RejectedAt(),
		CancelledAt:     model.CancelledAt(),
	}
}

// ToProposalDocumentEntity converts domain documents into persistence entities.
func ToProposalDocumentEntity(doc proposalmodel.ProposalDocumentInterface) entities.ProposalDocumentEntity {
	return entities.ProposalDocumentEntity{
		ID:            doc.ID(),
		ProposalID:    doc.ProposalID(),
		FileName:      doc.FileName(),
		MimeType:      doc.MimeType(),
		FileSizeBytes: doc.FileSizeBytes(),
		FileBlob:      doc.FileData(),
		UploadedAt:    doc.UploadedAt(),
	}
}
