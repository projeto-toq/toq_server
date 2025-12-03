#!/bin/bash
set -e

export AWS_SHARED_CREDENTIALS_FILE=/codigos/go_code/toq_server/configs/aws_credentials
export AWS_PROFILE=admin
export AWS_REGION=us-east-1

echo "Deploying validate..."
aws lambda update-function-code --function-name listing-media-validate-staging --zip-file fileb://aws/lambdas/bin/validate.zip > /dev/null

echo "Deploying thumbnails..."
aws lambda update-function-code --function-name listing-media-thumbnails-staging --zip-file fileb://aws/lambdas/bin/thumbnails.zip > /dev/null

echo "Deploying zip..."
aws lambda update-function-code --function-name listing-media-zip-staging --zip-file fileb://aws/lambdas/bin/zip.zip > /dev/null

echo "Deploying consolidate..."
aws lambda update-function-code --function-name listing-media-consolidate-staging --zip-file fileb://aws/lambdas/bin/consolidate.zip > /dev/null

echo "Deploying callback..."
aws lambda update-function-code --function-name listing-media-callback-staging --zip-file fileb://aws/lambdas/bin/callback.zip > /dev/null

echo "Updating Step Functions definition..."
aws stepfunctions update-state-machine \
	--state-machine-arn arn:aws:states:us-east-1:058264253741:stateMachine:listing-media-processing-sm-staging \
	--definition file://aws/step_functions/media_processing_pipeline.json > /dev/null

aws stepfunctions update-state-machine \
	--state-machine-arn arn:aws:states:us-east-1:058264253741:stateMachine:listing-media-finalization-sm-staging \
	--definition file://aws/step_functions/media_finalization_pipeline.json > /dev/null

echo "Deployment complete."
