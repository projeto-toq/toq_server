#!/bin/bash
set -e

export AWS_SHARED_CREDENTIALS_FILE=/codigos/go_code/toq_server/configs/aws_credentials
export AWS_PROFILE=admin
export AWS_REGION=us-east-1

ROOT_DIR=$(pwd)
BIN_DIR="$ROOT_DIR/aws/lambdas/bin"

shopt -s nullglob
zip_files=("$BIN_DIR"/*.zip)
shopt -u nullglob

if [ ${#zip_files[@]} -eq 0 ]; then
	echo "❌ Nenhum artefato encontrado em $BIN_DIR. Execute scripts/build_lambdas.sh antes."
	exit 1
fi

for zip_file in "${zip_files[@]}"; do
	lambda=$(basename "$zip_file" .zip)
	function_name="listing-media-${lambda}-staging"

	echo "Deploying $lambda..."
	aws lambda update-function-code \
		--function-name "$function_name" \
		--zip-file "fileb://$zip_file" > /dev/null
done

echo "Updating Step Functions definition..."
aws stepfunctions update-state-machine \
	--state-machine-arn arn:aws:states:us-east-1:058264253741:stateMachine:listing-media-processing-sm-staging \
	--definition file://aws/step_functions/media_processing_pipeline.json > /dev/null

aws stepfunctions update-state-machine \
	--state-machine-arn arn:aws:states:us-east-1:058264253741:stateMachine:listing-media-finalization-sm-staging \
	--definition file://aws/step_functions/media_finalization_pipeline.json > /dev/null

echo "Deployment complete."
