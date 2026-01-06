package mysqlproposaladapter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/projeto-toq/toq_server/internal/adapter/right/mysql/proposal/converters"
	proposalmodel "github.com/projeto-toq/toq_server/internal/core/model/proposal_model"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// CreateDocument inserts a proposal document row using the provided transaction.
func (a *ProposalAdapter) CreateDocument(ctx context.Context, tx *sql.Tx, document proposalmodel.ProposalDocumentInterface) error {
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return err
	}
	defer spanEnd()

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)

	entity := converters.ToProposalDocumentEntity(document)
	if entity.UploadedAt.IsZero() {
		entity.UploadedAt = time.Now().UTC()
		document.SetUploadedAt(entity.UploadedAt)
	}

	query := `INSERT INTO proposal_documents (
		proposal_id,
		file_name,
		mime_type,
		file_size_bytes,
		file_blob,
		uploaded_at
	) VALUES (?,?,?,?,?,?)`

	result, execErr := a.ExecContext(ctx, tx, "insert_proposal_document", query,
		entity.ProposalID,
		entity.FileName,
		entity.MimeType,
		entity.FileSizeBytes,
		entity.FileBlob,
		entity.UploadedAt,
	)
	if execErr != nil {
		utils.SetSpanError(ctx, execErr)
		logger.Error("mysql.proposal_document.create.exec_error", "proposal_id", entity.ProposalID, "err", execErr)
		return fmt.Errorf("create proposal document: %w", execErr)
	}

	id, idErr := result.LastInsertId()
	if idErr != nil {
		utils.SetSpanError(ctx, idErr)
		logger.Error("mysql.proposal_document.create.last_insert_id_error", "proposal_id", entity.ProposalID, "err", idErr)
		return fmt.Errorf("proposal document last insert id: %w", idErr)
	}

	document.SetID(id)
	return nil
}
