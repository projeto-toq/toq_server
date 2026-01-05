#!/bin/bash
set -e

export AWS_SHARED_CREDENTIALS_FILE=/codigos/go_code/toq_server/configs/aws_credentials
export AWS_PROFILE=admin
export AWS_REGION=us-east-1

ROOT_DIR=$(pwd)
BIN_DIR="$ROOT_DIR/aws/lambdas/bin"
POLICY_FILE="$ROOT_DIR/aws/step_functions/updated_policy.json"
DEFAULT_LAMBDA_ROLE_ARN="arn:aws:iam::058264253741:role/toq-media-processing-backend-staging"
FFMPEG_LAYER_NAME=${FFMPEG_LAYER_NAME:-"toq-ffmpeg-layer-staging"}
FFMPEG_LAYER_ARN=${FFMPEG_LAYER_ARN:-""}
FFMPEG_LAYER_ZIP="$BIN_DIR/ffmpeg-layer.zip"

# Roles/policies to update (backend + step functions)
IAM_ROLE_NAME=${IAM_ROLE_NAME:-"toq-media-processing-backend-staging"}
IAM_POLICY_NAME=${IAM_POLICY_NAME:-"toq-media-processing-backend-inline"}
STEP_ROLE_NAME=${STEP_ROLE_NAME:-"toq-media-processing-stepfunctions-staging"}
STEP_POLICY_NAME=${STEP_POLICY_NAME:-"MediaProcessingStepFunctionsPolicy"}

shopt -s nullglob
zip_files=("$BIN_DIR"/*.zip)
shopt -u nullglob

if [ ${#zip_files[@]} -eq 0 ]; then
	echo "❌ Nenhum artefato encontrado em $BIN_DIR. Execute scripts/build_lambdas.sh antes."
	exit 1
fi

# Publicar layer de FFmpeg se o zip existir
if [ -f "$FFMPEG_LAYER_ZIP" ]; then
	echo "Publicando layer FFmpeg a partir de $FFMPEG_LAYER_ZIP..."
	FFMPEG_LAYER_ARN=$(aws lambda publish-layer-version \
		--layer-name "$FFMPEG_LAYER_NAME" \
		--zip-file "fileb://$FFMPEG_LAYER_ZIP" \
		--compatible-runtimes provided.al2 \
		--query 'LayerVersionArn' --output text)
	echo "✅ Layer publicada: $FFMPEG_LAYER_ARN"
else
	echo "⚠️  Layer FFmpeg não encontrada em $FFMPEG_LAYER_ZIP. Prosseguindo sem atualizar layer."
fi

for zip_file in "${zip_files[@]}"; do
	lambda=$(basename "$zip_file" .zip)
	# Ignore layer artifact packaged as zip; handled separately above
	if [ "$lambda" = "ffmpeg-layer" ]; then
		echo "Skipping function deploy for ffmpeg-layer.zip (layer already published)."
		continue
	fi
	function_name="listing-media-${lambda}-staging"
	role_arn="${LAMBDA_ROLE_ARN:-$DEFAULT_LAMBDA_ROLE_ARN}"

	echo "Deploying $lambda..."

	# Create the function if it doesn't exist yet (new lambdas like video_thumbnails)
	if ! aws lambda get-function --function-name "$function_name" >/dev/null 2>&1; then
		echo "Function $function_name not found. Creating..."
		aws lambda create-function \
			--function-name "$function_name" \
			--runtime provided.al2 \
			--role "$role_arn" \
			--handler bootstrap \
			--timeout 300 \
			--memory-size 1024 \
			--zip-file "fileb://$zip_file" \
			${FFMPEG_LAYER_ARN:+--layers "$FFMPEG_LAYER_ARN"} \
			$( [ "$lambda" = "video_thumbnails" ] && echo --environment "Variables={FFMPEG_PATH=/opt/bin/ffmpeg}" ) >/dev/null
	else
		aws lambda update-function-code \
			--function-name "$function_name" \
			--zip-file "fileb://$zip_file" > /dev/null

			# Update config for video thumbnails to ensure ffmpeg path is set and layer attached if provided
			if [ "$lambda" = "video_thumbnails" ]; then
				if [ -n "$FFMPEG_LAYER_ARN" ]; then
					aws lambda update-function-configuration \
						--function-name "$function_name" \
						--layers "$FFMPEG_LAYER_ARN" \
						--environment "Variables={FFMPEG_PATH=/opt/bin/ffmpeg}" >/dev/null
				else
					echo "⚠️  FFMPEG_LAYER_ARN não definido; a função permanecerá sem layer."
					aws lambda update-function-configuration \
						--function-name "$function_name" \
						--environment "Variables={FFMPEG_PATH=/opt/bin/ffmpeg}" >/dev/null
				fi
			fi
	fi
done

if [ ! -f "$POLICY_FILE" ]; then
	echo "❌ Política não encontrada em $POLICY_FILE"
	exit 1
fi

echo "Updating IAM inline policy for role $IAM_ROLE_NAME..."
aws iam put-role-policy \
	--role-name "$IAM_ROLE_NAME" \
	--policy-name "$IAM_POLICY_NAME" \
	--policy-document "file://$POLICY_FILE" > /dev/null

echo "Updating IAM inline policy for role $STEP_ROLE_NAME..."
aws iam put-role-policy \
	--role-name "$STEP_ROLE_NAME" \
	--policy-name "$STEP_POLICY_NAME" \
	--policy-document "file://$POLICY_FILE" > /dev/null

echo "Updating Step Functions definition..."
aws stepfunctions update-state-machine \
	--state-machine-arn arn:aws:states:us-east-1:058264253741:stateMachine:listing-media-processing-sm-staging \
	--definition file://aws/step_functions/media_processing_pipeline.json > /dev/null

aws stepfunctions update-state-machine \
	--state-machine-arn arn:aws:states:us-east-1:058264253741:stateMachine:listing-media-finalization-sm-staging \
	--definition file://aws/step_functions/media_finalization_pipeline.json > /dev/null

echo "Deployment complete."
