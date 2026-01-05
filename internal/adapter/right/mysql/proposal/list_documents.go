package mysqlproposaladapter

import (
	"context"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/entities"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListDocuments returns all documents linked to a proposal ordered by uploaded_at desc.
func (a *ProposalAdapter) ListDocuments(ctx context.Context, proposalID int64) ([]proposalmodel.ProposalDocumentInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	query := `SELECT id, proposal_id, file_name, file_type, file_url, file_size_bytes, uploaded_at
	FROM proposal_documents
	WHERE proposal_id = ?
	ORDER BY uploaded_at DESC`

	rows, queryErr := a.QueryContext(ctx, nil, "list_proposal_documents", query, proposalID)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.proposal_document.list.query_error", "proposal_id", proposalID, "err", queryErr)
		return nil, fmt.Errorf("list proposal documents: %w", queryErr)
	}
	defer rows.Close()

	documents := make([]proposalmodel.ProposalDocumentInterface, 0)
	for rows.Next() {
		entity := entities.ProposalDocumentEntity{}
		if scanErr := rows.Scan(
			&entity.ID,
			&entity.ProposalID,
			&entity.FileName,
			&entity.FileType,
			&entity.FileURL,
			&entity.FileSizeBytes,
			&entity.UploadedAt,
		); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.proposal_document.list.scan_error", "proposal_id", proposalID, "err", scanErr)
			return nil, fmt.Errorf("scan proposal document: %w", scanErr)
		}
		documents = append(documents, converters.ToProposalDocumentModel(entity))
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.proposal_document.list.rows_error", "proposal_id", proposalID, "err", rowsErr)
		return nil, fmt.Errorf("iterate proposal documents: %w", rowsErr)
	}

	return documents, nil
}
