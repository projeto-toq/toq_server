# AWS Resources - TOQ Server

This directory contains the source code and definitions for AWS resources used in the TOQ Server media processing pipeline.

## Structure

- **lambdas/**: Contains the source code for AWS Lambda functions.
  - `go_src/`: Go source code for all lambdas (Hexagonal Architecture).
    - `cmd/`: Entry points for each lambda (`validate`, `thumbnails`, `zip`, `consolidate`, `callback`).
    - `internal/`: Core logic, ports, and adapters.
  - `bin/`: Compiled binaries (artifacts for deployment).

- **step_functions/**: Contains the definitions for AWS Step Functions.
  - `media_processing_pipeline.json`: State machine responsible for raw asset processing (`listing-media-processing-sm-staging`).
  - `media_finalization_pipeline.json`: State machine responsible for ZIP generation (`listing-media-finalization-sm-staging`).

## Development

The project uses Go 1.25+ for Lambda functions.

### Prerequisites
- Go 1.25+
- AWS CLI v2
- Zip utility

### Building Lambdas

Use the provided script to compile all lambdas and create deployment artifacts:

```bash
./scripts/build_lambdas.sh
```

This script will:
1. Clean the `aws/lambdas/bin` directory.
2. Compile each lambda in `aws/lambdas/go_src/cmd/*`.
3. Create optimized `.zip` files ready for AWS Lambda (using `provided.al2023` runtime).

## Deployment

To deploy updates to AWS, use the AWS CLI. Ensure you have the correct credentials configured.

### Credentials
Credentials should be located in `configs/aws_credentials`. You can export them as environment variables:

```bash
export AWS_SHARED_CREDENTIALS_FILE=$(pwd)/configs/aws_credentials
export AWS_PROFILE=admin
export AWS_REGION=us-east-1
```

### Updating Functions

Run the following script to update the code for all functions:

```bash
./scripts/deploy_lambdas.sh
```

The script now also updates both Step Functions definitions:
- `listing-media-processing-sm-staging` ⟶ `aws/step_functions/media_processing_pipeline.json` (garante `provider`, `traceparent`, `outputs[].errorCode`).
- `listing-media-finalization-sm-staging` ⟶ `aws/step_functions/media_finalization_pipeline.json` (garante payload `zipBundles` e `assetsZipped`).

To update only the processing machine:

```bash
aws stepfunctions update-state-machine \
  --state-machine-arn arn:aws:states:us-east-1:058264253741:stateMachine:listing-media-processing-sm-staging \
  --definition file://aws/step_functions/media_processing_pipeline.json
```

To update only the finalization machine:

```bash
aws stepfunctions update-state-machine \
  --state-machine-arn arn:aws:states:us-east-1:058264253741:stateMachine:listing-media-finalization-sm-staging \
  --definition file://aws/step_functions/media_finalization_pipeline.json
```

## Architecture Notes

- **Runtime**: `provided.al2023` (Go custom runtime).
- **Architecture**: Hexagonal (Ports & Adapters).
- **Image Processing**: Uses `disintegration/imaging` for high-quality resizing and EXIF rotation handling.
- **S3 Paths**:
  - Raw: `/{listingIdentityId}/raw/{mediaTypeSegment}/{reference}-{filename}`
  - Processed: `/{listingIdentityId}/processed/{mediaTypeSegment}/{size}/{filename}`
  - Zip: `/{listingIdentityId}/processed/zip/listing-media.zip` (ZIP único sobrescrito a cada finalização bem-sucedida)

## Troubleshooting - Finalização de ZIP
1. Chame `POST /api/v2/listings/media/uploads/complete` e capture `listing_identity_id` no log. Em caso de sucesso, procure por `service.media.complete.started_zip`; o log traz `job_id` e `execution_arn` de `listing-media-finalization-sm-staging`.
2. Com o `job_id` em mãos, consulte a tabela MySQL:
  ```bash
  mysql -h 127.0.0.1 -P 3306 -utoq -ptoq toq_server \
    -e "SELECT id,status,external_id,started_at,last_error FROM media_processing_jobs WHERE id = <job_id>;"
  ```
  O campo `external_id` deve coincidir com o `execution_arn` retornado pelo Step Functions de finalização.
3. Para inspecionar o workflow, execute:
  ```bash
  aws stepfunctions describe-execution --execution-arn <execution_arn>
  aws stepfunctions get-execution-history --execution-arn <execution_arn> --reverse-order --max-results 200
  ```
  Falhas em `CreateZipBundle` ou `FinalizeAndCallback` aparecem no histórico; o backend replica os detalhes em `media_processing_jobs.callback_body`.
4. Se o campo `status` permanecer `PENDING` ou `RUNNING` por mais de alguns minutos, revise as Lambdas de zip/consolidation (`aws logs tail /aws/lambda/listing-media-zip-staging --follow`).
5. Após `status=SUCCEEDED`, o callback atualiza os assets e disponibiliza o ZIP em `/{listingIdentityId}/processed/zip/listing-media.zip`. Confirme com `aws s3 ls s3://toq-listing-medias/<listingIdentityId>/processed/zip/`.
