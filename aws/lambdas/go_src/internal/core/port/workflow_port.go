package port

import "context"

// WorkflowPort defines the interface for workflow orchestration (e.g., Step Functions)
type WorkflowPort interface {
	// StartExecution starts a new execution of the state machine
	StartExecution(ctx context.Context, stateMachineArn string, input string) error
}
