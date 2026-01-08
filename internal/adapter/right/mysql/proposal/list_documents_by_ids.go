package mysqlproposaladapter

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/converters"
	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/entities"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// ListDocumentsByProposalIDs returns documents grouped by proposal id in a single round trip.
func (a *ProposalAdapter) ListDocumentsByProposalIDs(
	ctx context.Context,
	tx *sql.Tx,
	proposalIDs []int64,
	includeBlob bool,
) (map[int64][]proposalmodel.ProposalDocumentInterface, error) {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return nil, err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	if len(proposalIDs) == 0 {
		return map[int64][]proposalmodel.ProposalDocumentInterface{}, nil
	}

	placeholders := make([]string, len(proposalIDs))
	args := make([]interface{}, len(proposalIDs))
	for i, id := range proposalIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	blobExpr := "file_blob"
	if !includeBlob {
		blobExpr = "NULL AS file_blob"
	}

	query := fmt.Sprintf(`SELECT id, proposal_id, file_name, mime_type, file_size_bytes, %s, uploaded_at
		FROM proposal_documents
		WHERE proposal_id IN (%s)
		ORDER BY proposal_id ASC, uploaded_at DESC`, blobExpr, strings.Join(placeholders, ","))

	rows, queryErr := a.QueryContext(ctx, tx, "list_proposal_documents_bulk", query, args...)
	if queryErr != nil {
		utils.SetSpanError(ctx, queryErr)
		logger.Error("mysql.proposal_document.list.bulk_query_error", "err", queryErr)
		return nil, fmt.Errorf("list proposal documents bulk: %w", queryErr)
	}
	defer rows.Close()

	result := make(map[int64][]proposalmodel.ProposalDocumentInterface, len(proposalIDs))

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
			logger.Error("mysql.proposal_document.list.bulk_scan_error", "err", scanErr)
			return nil, fmt.Errorf("scan proposal document: %w", scanErr)
		}
		result[entity.ProposalID] = append(result[entity.ProposalID], converters.ToProposalDocumentModel(entity, includeBlob))
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		utils.SetSpanError(ctx, rowsErr)
		logger.Error("mysql.proposal_document.list.bulk_rows_error", "err", rowsErr)
		return nil, fmt.Errorf("iterate proposal documents bulk: %w", rowsErr)
	}

	return result, nil
}
