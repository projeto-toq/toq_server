package mysqlproposaladapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/entities"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListDocuments returns proposal documents optionally loading the binary blob.
func (a *ProposalAdapter) ListDocuments(ctx context.Context, tx *sql.Tx, proposalID int64, includeBlob bool) ([]proposalmodel.ProposalDocumentInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	blobExpr := "file_blob"
	if !includeBlob {
		blobExpr = "NULL AS file_blob"
	}

	query := fmt.Sprintf(`SELECT id, proposal_id, file_name, mime_type, file_size_bytes, %s, uploaded_at
	FROM proposal_documents
	WHERE proposal_id = ?
	ORDER BY uploaded_at DESC`, blobExpr)

	rows, queryErr := a.QueryContext(ctx, tx, "list_proposal_documents", query, proposalID)
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
			&entity.MimeType,
			&entity.FileSizeBytes,
			&entity.FileBlob,
			&entity.UploadedAt,
		); scanErr != nil {
			utils.SetSpanError(ctx, scanErr)
			logger.Error("mysql.proposal_document.list.scan_error", "proposal_id", proposalID, "err", scanErr)
			return nil, fmt.Errorf("scan proposal document: %w", scanErr)
		}
		documents = append(documents, converters.ToProposalDocumentModel(entity, includeBlob))
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.proposal_document.list.rows_error", "proposal_id", proposalID, "err", rowsErr)
		return nil, fmt.Errorf("iterate proposal documents: %w", rowsErr)
	}

	return documents, nil
}
