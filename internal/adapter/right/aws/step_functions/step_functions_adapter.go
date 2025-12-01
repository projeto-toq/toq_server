package stepfunctions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	"github.com/projeto-toq/toq_server/internal/core/port/right/workflow"
)

type StepFunctionsAdapter struct {
	client          *sfn.Client
	finalizationARN string
}

func NewStepFunctionsAdapter(cfg aws.Config, finalizationARN string) *StepFunctionsAdapter {
	return &StepFunctionsAdapter{
		client:          sfn.NewFromConfig(cfg),
		finalizationARN: finalizationARN,
	}
}

func (a *StepFunctionsAdapter) StartMediaFinalization(ctx context.Context, input mediaprocessingmodel.MediaFinalizationInput) (string, error) {
	payloadBytes, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("failed to marshal input: %w", err)
	}
	payload := string(payloadBytes)

	out, err := a.client.StartExecution(ctx, &sfn.StartExecutionInput{
		StateMachineArn: aws.String(a.finalizationARN),
		Input:           aws.String(payload),
		Name:            aws.String(fmt.Sprintf("finalization-%d-%d", input.ListingID, input.JobID)),
	})
	if err != nil {
		return "", err
	}

	return *out.ExecutionArn, nil
}

var _ workflow.WorkflowPortInterface = (*StepFunctionsAdapter)(nil)
