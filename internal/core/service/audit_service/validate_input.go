package auditservice

import (
	"errors"
	"fmt"

	auditmodel "github.com/projeto-toq/toq_server/internal/core/model/audit_model"
)

func (s *auditService) validateInput(input auditmodel.RecordInput) error {
	if input.Target.ID == 0 {
		return errors.New("target_id is required")
	}
	if input.Target.Type == "" {
		return errors.New("target_type is required")
	}
	if input.Operation == "" {
		return errors.New("operation is required")
	}
	// Optional but validated when present
	if input.Target.Version != nil && *input.Target.Version < 0 {
		return fmt.Errorf("target_version must be positive")
	}
	return nil
}
