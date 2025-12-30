package stepfunctions

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
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

var _ workflow.WorkflowPortInterface = (*StepFunctionsAdapter)(nil)
