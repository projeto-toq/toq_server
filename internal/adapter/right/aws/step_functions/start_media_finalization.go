package stepfunctions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/smithy-go"
	derrors "github.com/projeto-toq/toq_server/internal/core/derrors"
	mediaprocessingmodel "github.com/projeto-toq/toq_server/internal/core/model/media_processing_model"
	"github.com/projeto-toq/toq_server/internal/core/port/right/workflow"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

// StartMediaFinalization triggers the AWS Step Functions workflow responsible for building the ZIP bundle.
//
// The method serializes the MediaFinalizationInput payload, starts the execution and
// maps AccessDenied responses to workflow.ErrFinalizationAccessDenied so services can
// convert IAM misconfigurations into actionable domain errors.
func (a *StepFunctionsAdapter) StartMediaFinalization(
	ctx context.Context,
	input mediaprocessingmodel.MediaFinalizationInput,
) (string, error) {
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
		return "", fmt.Errorf("marshal finalization payload: %w", err)
	}

	out, err := a.client.StartExecution(ctx, &sfn.StartExecutionInput{
		StateMachineArn: aws.String(a.finalizationARN),
		Input:           aws.String(string(payloadBytes)),
		Name:            aws.String(fmt.Sprintf("finalization-%d-%d", input.ListingIdentityID, input.JobID)),
	})
	if err != nil {
		utils.SetSpanError(ctx, err)
		logger.Error("adapter.stepfunctions.start_execution_error",
			"err", err,
			"listing_identity_id", input.ListingIdentityID,
			"job_id", input.JobID,
		)

		if isAccessDenied(err) {
			return "", workflow.ErrFinalizationAccessDenied
		}

		return "", fmt.Errorf("start finalization workflow: %w", err)
	}

	executionARN := aws.ToString(out.ExecutionArn)
	logger.Info("adapter.stepfunctions.start_execution_success",
		"execution_arn", executionARN,
		"listing_identity_id", input.ListingIdentityID,
		"job_id", input.JobID,
	)

	return executionARN, nil
}

func isAccessDenied(err error) bool {
	var apiErr smithy.APIError
	if !errors.As(err, &apiErr) {
		return false
	}

	return apiErr.ErrorCode() == "AccessDeniedException"
}
