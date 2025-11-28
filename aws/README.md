# AWS Resources

This directory contains the source code and definitions for AWS resources used in the TOQ Server media processing pipeline.

## Structure

- **lambdas/**: Contains the source code for AWS Lambda functions.
  - `validate/`: Code for the validation Lambda (Node.js).
  - `zip_generator/`: Code for the ZIP generation Lambda (Python).
  - `thumbnail_generator/`: Code for the thumbnail generation Lambda (Node.js).

- **step_functions/**: Contains the definitions for AWS Step Functions.
  - `media_processing_pipeline.json`: The JSON definition of the media processing state machine.

## Usage

These resources are deployed to AWS. If you need to modify them, update the code here and redeploy using the appropriate AWS CLI commands or CI/CD pipeline.

### Redeploying Lambdas (Example)

To update the ZIP generator lambda:
```bash
cd aws/lambdas/zip_generator
zip -r ../zip_generator.zip .
aws lambda update-function-code --function-name listing-media-zip-staging --zip-file fileb://../zip_generator.zip
```
