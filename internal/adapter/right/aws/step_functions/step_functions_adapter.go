package stepfunctions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	"github.com/projeto-toq/toq_server/internal/core/port/right/workflow"
	"github.com/projeto-toq/toq_server/internal/core/utils"
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
	ctx = utils.ContextWithLogger(ctx)
	ctx, spanEnd, err := utils.GenerateTracer(ctx)
	if err != nil {
		return "", derrors.Infra("failed to create tracer", err)
	}
	defer spanEnd()

	logger := utils.LoggerFromContext(ctx)

	payloadBytes, err := json.Marshal(input)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return "", fmt.Errorf("failed to marshal input: %w", err)
	}
	payload := string(payloadBytes)

	out, err := a.client.StartExecution(ctx, &sfn.StartExecutionInput{
		StateMachineArn: aws.String(a.finalizationARN),
		Input:           aws.String(payload),
		Name:            aws.String(fmt.Sprintf("finalization-%d-%d", input.ListingIdentityID, input.JobID)),
	})
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("adapter.stepfunctions.start_execution_error", "err", err, "listing_identity_id", input.ListingIdentityID, "job_id", input.JobID)
		return "", err
	}

	logger.Info("adapter.stepfunctions.start_execution_success", "execution_arn", *out.ExecutionArn, "listing_identity_id", input.ListingIdentityID, "job_id", input.JobID)
	return *out.ExecutionArn, nil
}

var _ workflow.WorkflowPortInterface = (*StepFunctionsAdapter)(nil)
