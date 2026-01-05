package converters

import (
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/entities"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
)

// ToProposalEntity converts a domain ProposalInterface into a ProposalEntity for persistence.
func ToProposalEntity(model proposalmodel.ProposalInterface) entities.ProposalEntity {
	return entities.ProposalEntity{
		ID:                model.ID(),
		ListingIdentityID: model.ListingIdentityID(),
		RealtorID:         model.RealtorID(),
		OwnerID:           model.OwnerID(),
		TransactionType:   string(model.TransactionType()),
		PaymentMethod:     string(model.PaymentMethod()),
		ProposedValue:     model.ProposedValue(),
		OriginalValue:     model.OriginalValue(),
		DownPayment:       model.DownPayment(),
		Installments:      model.Installments(),
		AcceptsExchange:   model.AcceptsExchange(),
		RentalMonths:      model.RentalMonths(),
		GuaranteeType:     model.GuaranteeType(),
		SecurityDeposit:   model.SecurityDeposit(),
		ClientName:        model.ClientName(),
		ClientPhone:       model.ClientPhone(),
		ProposalNotes:     model.ProposalNotes(),
		OwnerNotes:        model.OwnerNotes(),
		RejectionReason:   model.RejectionReason(),
		Status:            string(model.Status()),
		ExpiresAt:         model.ExpiresAt(),
		AcceptedAt:        model.AcceptedAt(),
		RejectedAt:        model.RejectedAt(),
		CancelledAt:       model.CancelledAt(),
		IsFavorite:        false,
		CreatedAt:         model.CreatedAt(),
		UpdatedAt:         model.UpdatedAt(),
		Deleted:           false,
	}
}

// ToProposalDocumentEntity converts a ProposalDocumentInterface into a ProposalDocumentEntity.
func ToProposalDocumentEntity(doc proposalmodel.ProposalDocumentInterface) entities.ProposalDocumentEntity {
	return entities.ProposalDocumentEntity{
		ID:            doc.ID(),
		ProposalID:    doc.ProposalID(),
		FileName:      doc.FileName(),
		FileType:      doc.FileType(),
		FileURL:       doc.FileURL(),
		FileSizeBytes: doc.FileSizeBytes(),
		UploadedAt:    doc.UploadedAt(),
	}
}
