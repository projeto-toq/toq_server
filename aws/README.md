# AWS Resources - TOQ Server

This directory contains the source code and definitions for AWS resources used in the TOQ Server media processing pipeline.

## Structure

- **lambdas/**: Contains the source code for AWS Lambda functions.
  - `go_src/`: Go source code for all lambdas (Hexagonal Architecture).
    - `cmd/`: Entry points for each lambda (`validate`, `thumbnails`, `zip`, `consolidate`, `callback`).
    - `internal/`: Core logic, ports, and adapters.
  - `bin/`: Compiled binaries (artifacts for deployment).

- **step_functions/**: Contains the definitions for AWS Step Functions.
  - `media_processing_pipeline.json`: The JSON definition of the media processing state machine.

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

Run the following commands to update the code for each function:

```bash
# Update Thumbnails Lambda
aws lambda update-function-code --function-name listing-media-thumbnails-staging --zip-file fileb://aws/lambdas/bin/thumbnails.zip

# Update Zip Lambda
aws lambda update-function-code --function-name listing-media-zip-staging --zip-file fileb://aws/lambdas/bin/zip.zip

# Update Validate Lambda
aws lambda update-function-code --function-name listing-media-validate-staging --zip-file fileb://aws/lambdas/bin/validate.zip

# Update Consolidate Lambda
aws lambda update-function-code --function-name listing-media-consolidate-staging --zip-file fileb://aws/lambdas/bin/consolidate.zip

# Update Callback Lambda
aws lambda update-function-code --function-name listing-media-callback-staging --zip-file fileb://aws/lambdas/bin/callback.zip
```

## Architecture Notes

- **Runtime**: `provided.al2023` (Go custom runtime).
- **Architecture**: Hexagonal (Ports & Adapters).
- **Image Processing**: Uses `disintegration/imaging` for high-quality resizing and EXIF rotation handling.
- **S3 Paths**:
  - Raw: `/{listingId}/raw/{mediaType}/{uuid}.{ext}`
  - Processed: `/{listingId}/processed/{mediaType}/{size}/{uuid}.{ext}`
  - Zip: `/{listingId}/processed/zip/{batchId}.zip`
