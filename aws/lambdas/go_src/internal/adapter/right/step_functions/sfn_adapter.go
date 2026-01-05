package stepfunctions

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/projeto-toq/toq_server/aws/lambdas/go_src/internal/core/port"
)

type SfnAdapter struct {
	client *sfn.Client
}

func NewSfnAdapter(client *sfn.Client) port.WorkflowPort {
	return &SfnAdapter{
		client: client,
	}
}

func (a *SfnAdapter) StartExecution(ctx context.Context, stateMachineArn string, input string) (string, error) {
	out, err := a.client.StartExecution(ctx, &sfn.StartExecutionInput{
		StateMachineArn: aws.String(stateMachineArn),
		Input:           aws.String(input),
	})
	if err != nil {
		return "", fmt.Errorf("failed to start execution: %w", err)
	}
	return aws.ToString(out.ExecutionArn), nil
}
