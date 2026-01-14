terraform {
  required_version = ">= 1.7"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# KMS key for media processing
module "kms" {
  source      = "./modules/kms_key"
  alias       = "alias/toq-media-processing-staging"
  description = "Chave KMS para processamento de mídias TOQ Server - Staging"
}

# Logs bucket (no versioning, lifecycle 7d)
module "s3_logs" {
  source      = "./modules/s3_bucket_logs"
  bucket_name = "toq-logs-staging"
  expire_days = 7
  enable_lifecycle = true
}

# Media bucket (versioned, SSE-KMS, logging to logs bucket, lifecycle raw->Glacier 180d)
module "s3_media" {
  source      = "./modules/s3_bucket_media"
  bucket_name = "toq-listing-medias"
  kms_key_arn = module.kms.arn
  log_bucket  = module.s3_logs.name
}

# User media bucket (SSE AES256, no versioning/lifecycle/logging)
module "s3_user_medias" {
  source      = "./modules/s3_bucket_basic"
  bucket_name = local.user_media_bucket
}

# SQS DLQ
module "sqs_dlq" {
  source                      = "./modules/sqs_queue"
  name                        = "listing-media-processing-dlq-staging"
  visibility_timeout_seconds  = 300
  message_retention_seconds   = 1209600
  receive_wait_time_seconds   = 0
  kms_master_key_id           = module.kms.alias
  kms_data_key_reuse_period_seconds = 300
}

# SQS main
module "sqs_main" {
  source                      = "./modules/sqs_queue"
  name                        = "listing-media-processing-staging"
  visibility_timeout_seconds  = 60
  message_retention_seconds   = 345600
  receive_wait_time_seconds   = 20
  kms_master_key_id           = module.kms.alias
  kms_data_key_reuse_period_seconds = 300
  redrive_policy              = jsonencode({
    deadLetterTargetArn = module.sqs_dlq.arn,
    maxReceiveCount     = 5
  })
}

locals {
  media_bucket = "toq-listing-medias"
  logs_bucket  = "toq-logs-staging"
  user_media_bucket = "toq-user-medias"
  kms_alias    = "alias/toq-media-processing-staging"
  account_id   = "058264253741"
}

# IAM roles
module "iam_lambda" {
  source             = "./modules/iam_role"
  name               = "toq-media-processing-lambda-staging"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume.json
  inline_policies    = [
    {
      name     = "lambda-logs-kms-s3-sqs"
      document = data.aws_iam_policy_document.lambda_policy.json
    },
    {
      name     = "lambda-start-sfn"
      document = data.aws_iam_policy_document.lambda_start_sfn.json
    }
  ]
}

module "iam_stepfunctions" {
  source             = "./modules/iam_role"
  name               = "toq-media-processing-stepfunctions-staging"
  assume_role_policy = data.aws_iam_policy_document.sfn_assume.json
  inline_policies    = [
    {
      name     = "sfn-lambda-mediaconvert-logs"
      document = data.aws_iam_policy_document.sfn_policy.json
    }
  ]
}

module "iam_mediaconvert" {
  source             = "./modules/iam_role"
  name               = "toq-media-processing-mediaconvert-staging"
  assume_role_policy = data.aws_iam_policy_document.mediaconvert_assume.json
  inline_policies    = [
    {
      name     = "mediaconvert-s3"
      document = data.aws_iam_policy_document.mediaconvert_policy.json
    }
  ]
}

# Lambda functions (artifacts expected at aws/lambdas/bin/*.zip)
module "lambda_validate" {
  source       = "./modules/lambda_function"
  function_name = "listing-media-validate-staging"
  role_arn      = module.iam_lambda.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  filename      = "${path.module}/lambdas/bin/validate.zip"
  memory_size   = 512
  timeout       = 60
  environment   = {
    MEDIA_BUCKET    = local.media_bucket
    ENV             = "staging"
    TRACE_HEADER_KEY = "traceparent"
  }
  tracing_mode  = "Active"
}

module "lambda_thumbnails" {
  source       = "./modules/lambda_function"
  function_name = "listing-media-thumbnails-staging"
  role_arn      = module.iam_lambda.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  filename      = "${path.module}/lambdas/bin/thumbnails.zip"
  memory_size   = 2048
  timeout       = 300
  environment   = {
    MEDIA_BUCKET    = local.media_bucket
    ENV             = "staging"
    TRACE_HEADER_KEY = "traceparent"
  }
  tracing_mode  = "Active"
}

module "lambda_video_thumbnails" {
  source       = "./modules/lambda_function"
  function_name = "listing-media-video_thumbnails-staging"
  role_arn      = module.iam_lambda.arn
  handler       = "bootstrap"
  runtime       = "provided.al2"
  filename      = "${path.module}/lambdas/bin/video_thumbnails.zip"
  memory_size   = 1024
  timeout       = 300
  environment   = {
    FFMPEG_PATH = "/opt/bin/ffmpeg"
  }
  tracing_mode  = "PassThrough"
  layers        = [var.ffmpeg_layer_arn]
}

module "lambda_consolidate" {
  source       = "./modules/lambda_function"
  function_name = "listing-media-consolidate-staging"
  role_arn      = module.iam_lambda.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  filename      = "${path.module}/lambdas/bin/consolidate.zip"
  memory_size   = 128
  timeout       = 3
  tracing_mode  = "PassThrough"
}

module "lambda_callback" {
  source       = "./modules/lambda_function"
  function_name = "listing-media-callback-staging"
  role_arn      = module.iam_lambda.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  filename      = "${path.module}/lambdas/bin/callback.zip"
  memory_size   = 128
  timeout       = 3
  tracing_mode  = "PassThrough"
  environment   = {
    ENV             = "staging"
    TRACE_HEADER_KEY = "traceparent"
  }
}

module "lambda_callback_dispatch" {
  source       = "./modules/lambda_function"
  function_name = "listing-media-callback-dispatch-staging"
  role_arn      = module.iam_lambda.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  filename      = "${path.module}/lambdas/bin/callback_dispatch.zip"
  memory_size   = 256
  timeout       = 60
  tracing_mode  = "Active"
  environment   = {
    CALLBACK_SECRET    = var.callback_secret
    CALLBACK_URL       = var.callback_url
    INTERNAL_API_TOKEN = var.internal_api_token
    ENV                = "staging"
    TRACE_HEADER_KEY   = "traceparent"
  }
}

module "lambda_zip" {
  source       = "./modules/lambda_function"
  function_name = "listing-media-zip-staging"
  role_arn      = module.iam_lambda.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  filename      = "${path.module}/lambdas/bin/zip.zip"
  memory_size   = 3008
  timeout       = 900
  tracing_mode  = "Active"
  ephemeral_storage_size = 2048
  environment   = {
    MEDIA_BUCKET    = local.media_bucket
    ENV             = "staging"
    TRACE_HEADER_KEY = "traceparent"
  }
}

# Step Functions
module "sfn_processing" {
  source       = "./modules/stepfunctions"
  name         = "listing-media-processing-sm-staging"
  role_arn     = module.iam_stepfunctions.arn
  definition_path = "${path.module}/step_functions/media_processing_pipeline.json"
  logging_level = "ALL"
  logging_include_execution_data = true
  logging_log_group_arn = aws_cloudwatch_log_group.sfn_processing.arn
  tracing_enabled = true
}

module "sfn_finalization" {
  source       = "./modules/stepfunctions"
  name         = "listing-media-finalization-sm-staging"
  role_arn     = module.iam_stepfunctions.arn
  definition_path = "${path.module}/step_functions/media_finalization_pipeline.json"
  logging_level = "OFF"
  logging_include_execution_data = false
  logging_log_group_arn = ""
  tracing_enabled = false
}

resource "aws_cloudwatch_log_group" "sfn_processing" {
  name              = "/aws/stepfunctions/listing-media-processing-sm-staging"
  retention_in_days = 30
}

data "aws_iam_policy_document" "lambda_assume" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "sfn_assume" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["states.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "mediaconvert_assume" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["mediaconvert.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "lambda_policy" {
  statement {
    sid     = "Logs"
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    resources = ["arn:aws:logs:${var.aws_region}:${local.account_id}:log-group:/aws/lambda/listing-media-*:*"]
  }

  statement {
    sid     = "S3Media"
    actions = [
      "s3:GetObject",
      "s3:PutObject",
      "s3:DeleteObject",
      "s3:ListBucket",
      "s3:GetObjectAttributes",
      "s3:HeadObject"
    ]
    resources = [
      "arn:aws:s3:::${local.media_bucket}",
      "arn:aws:s3:::${local.media_bucket}/*"
    ]
  }

  statement {
    sid     = "SQSMainAndDLQ"
    actions = [
      "sqs:SendMessage",
      "sqs:ReceiveMessage",
      "sqs:DeleteMessage",
      "sqs:GetQueueAttributes",
      "sqs:GetQueueUrl"
    ]
    resources = [module.sqs_main.arn, module.sqs_dlq.arn]
  }

  statement {
    sid     = "KMS"
    actions = ["kms:Encrypt", "kms:Decrypt", "kms:GenerateDataKey", "kms:DescribeKey"]
    resources = [module.kms.arn]
  }

  statement {
    sid     = "StepFunctionsTasks"
    actions = ["states:SendTaskSuccess", "states:SendTaskFailure", "states:SendTaskHeartbeat"]
    resources = ["*"]
  }

  statement {
    sid     = "MediaConvert"
    actions = ["mediaconvert:CreateJob", "mediaconvert:GetJob", "mediaconvert:ListJobs", "mediaconvert:DescribeEndpoints"]
    resources = ["*"]
  }

  statement {
    sid     = "XRay"
    actions = ["xray:PutTraceSegments", "xray:PutTelemetryRecords"]
    resources = ["*"]
  }
}

data "aws_iam_policy_document" "lambda_start_sfn" {
  statement {
    actions   = ["states:StartExecution"]
    resources = ["arn:aws:states:${var.aws_region}:${local.account_id}:stateMachine:listing-media-processing-sm-staging"]
  }
}

data "aws_iam_policy_document" "sfn_policy" {
  statement {
    sid     = "InvokeLambdas"
    actions = ["lambda:InvokeFunction", "lambda:InvokeAsync"]
    resources = [
      module.lambda_validate.arn,
      module.lambda_thumbnails.arn,
      module.lambda_video_thumbnails.arn,
      module.lambda_consolidate.arn,
      module.lambda_callback.arn,
      module.lambda_callback_dispatch.arn,
      module.lambda_zip.arn
    ]
  }

  statement {
    sid     = "MediaConvert"
    actions   = ["mediaconvert:CreateJob", "mediaconvert:GetJob", "mediaconvert:DescribeEndpoints"]
    resources = ["*"]
  }

  statement {
    sid       = "PassRoleForMediaConvert"
    actions   = ["iam:PassRole"]
    resources = ["arn:aws:iam::${local.account_id}:role/toq-media-processing-mediaconvert-staging"]
    condition {
      test     = "StringEquals"
      variable = "iam:PassedToService"
      values   = ["mediaconvert.amazonaws.com"]
    }
  }

  statement {
    sid     = "SQSSend"
    actions = ["sqs:SendMessage"]
    resources = [module.sqs_main.arn, module.sqs_dlq.arn]
  }

  statement {
    sid     = "StartFinalizationWorkflow"
    actions = ["states:StartExecution"]
    resources = ["arn:aws:states:${var.aws_region}:${local.account_id}:stateMachine:listing-media-finalization-sm-staging"]
  }

  statement {
    sid     = "Logs"
    actions = ["logs:CreateLogDelivery", "logs:GetLogDelivery", "logs:UpdateLogDelivery", "logs:DeleteLogDelivery", "logs:ListLogDeliveries", "logs:PutLogEvents", "logs:PutResourcePolicy", "logs:DescribeResourcePolicies", "logs:DescribeLogGroups"]
    resources = ["*"]
  }

  statement {
    sid     = "XRay"
    actions = ["xray:PutTraceSegments", "xray:PutTelemetryRecords", "xray:GetSamplingRules", "xray:GetSamplingTargets"]
    resources = ["*"]
  }
}

data "aws_iam_policy_document" "mediaconvert_policy" {
  statement {
    sid     = "S3"
    actions = ["s3:GetObject", "s3:PutObject", "s3:ListBucket"]
    resources = [
      "arn:aws:s3:::${local.media_bucket}",
      "arn:aws:s3:::${local.media_bucket}/*"
    ]
  }

  statement {
    sid     = "KMS"
    actions = ["kms:Decrypt", "kms:Encrypt", "kms:GenerateDataKey", "kms:DescribeKey"]
    resources = [module.kms.arn]
  }
}

data "aws_iam_policy_document" "ec2_assume" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "backend_policy_core" {
  statement {
    sid     = "S3MediaBucketAccess"
    actions = [
      "s3:PutObject",
      "s3:GetObject",
      "s3:HeadObject",
      "s3:DeleteObject",
      "s3:GetObjectAttributes",
      "s3:ListBucket"
    ]
    resources = [
      "arn:aws:s3:::${local.media_bucket}",
      "arn:aws:s3:::${local.media_bucket}/*"
    ]
  }

  statement {
    sid     = "AllowUserMediaBucket"
    actions = [
      "s3:PutObject",
      "s3:GetObject",
      "s3:DeleteObject",
      "s3:ListBucket",
      "s3:PutObjectAcl",
      "s3:GetObjectAcl"
    ]
    resources = [
      "arn:aws:s3:::${local.user_media_bucket}",
      "arn:aws:s3:::${local.user_media_bucket}/*"
    ]
  }

  statement {
    sid     = "SQSQueueAccess"
    actions = [
      "sqs:SendMessage",
      "sqs:GetQueueAttributes",
      "sqs:GetQueueUrl"
    ]
    resources = [module.sqs_main.arn, module.sqs_dlq.arn]
  }

  statement {
    sid     = "KMSKeyAccess"
    actions = ["kms:Decrypt", "kms:Encrypt", "kms:GenerateDataKey", "kms:DescribeKey"]
    resources = [module.kms.arn]
  }

  statement {
    sid     = "StepFunctionsAccess"
    actions = ["states:StartExecution", "states:DescribeExecution", "states:GetExecutionHistory"]
    resources = ["arn:aws:states:${var.aws_region}:${local.account_id}:stateMachine:listing-media-processing-sm-staging"]
  }
}

data "aws_iam_policy_document" "backend_policy_extended" {
  statement {
    sid     = "InvokeLambdas"
    actions = ["lambda:InvokeFunction"]
    resources = [
      module.lambda_validate.arn,
      module.lambda_thumbnails.arn,
      module.lambda_video_thumbnails.arn,
      module.lambda_zip.arn,
      module.lambda_callback_dispatch.arn,
      module.lambda_consolidate.arn
    ]
  }

  statement {
    sid     = "MediaConvertJobs"
    actions = ["mediaconvert:CreateJob", "mediaconvert:GetJob", "mediaconvert:DescribeEndpoints"]
    resources = ["*"]
  }

  statement {
    sid     = "PassRoleForMediaConvert"
    actions = ["iam:PassRole"]
    resources = ["arn:aws:iam::${local.account_id}:role/toq-media-processing-mediaconvert-staging"]
    condition {
      test     = "StringEquals"
      variable = "iam:PassedToService"
      values   = ["mediaconvert.amazonaws.com"]
    }
  }

  statement {
    sid     = "SQSSendMessage"
    actions = ["sqs:SendMessage"]
    resources = [module.sqs_main.arn, module.sqs_dlq.arn]
  }

  statement {
    sid     = "StartFinalizationWorkflow"
    actions = ["states:StartExecution"]
    resources = ["arn:aws:states:${var.aws_region}:${local.account_id}:stateMachine:listing-media-finalization-sm-staging"]
  }

  statement {
    sid     = "CloudWatchLogs"
    actions = [
      "logs:CreateLogDelivery", "logs:GetLogDelivery", "logs:UpdateLogDelivery", "logs:DeleteLogDelivery",
      "logs:ListLogDeliveries", "logs:PutResourcePolicy", "logs:DescribeResourcePolicies", "logs:DescribeLogGroups"
    ]
    resources = ["*"]
  }

  statement {
    sid     = "XRayTracing"
    actions = ["xray:PutTraceSegments", "xray:PutTelemetryRecords", "xray:GetSamplingRules", "xray:GetSamplingTargets"]
    resources = ["*"]
  }
}

module "iam_backend" {
  source             = "./modules/iam_role"
  name               = "toq-media-processing-backend-staging"
  assume_role_policy = data.aws_iam_policy_document.ec2_assume.json
  inline_policies = [
    {
      name     = "MediaProcessingBackendPolicy"
      document = data.aws_iam_policy_document.backend_policy_core.json
    },
    {
      name     = "toq-media-processing-backend-inline"
      document = data.aws_iam_policy_document.backend_policy_extended.json
    }
  ]
}
