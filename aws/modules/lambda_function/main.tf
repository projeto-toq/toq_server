terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

resource "aws_lambda_function" "this" {
  function_name = var.function_name
  role          = var.role_arn
  handler       = var.handler
  runtime       = var.runtime
  filename      = var.filename
  source_code_hash = filebase64sha256(var.filename)
  memory_size   = var.memory_size
  timeout       = var.timeout
  publish       = false
  architectures = ["x86_64"]

  dynamic "environment" {
    for_each = length(var.environment) > 0 ? [1] : []
    content {
      variables = var.environment
    }
  }

  dynamic "tracing_config" {
    for_each = var.tracing_mode != "" ? [1] : []
    content {
      mode = var.tracing_mode
    }
  }

  layers = var.layers

  dynamic "ephemeral_storage" {
    for_each = var.ephemeral_storage_size != 512 ? [1] : []
    content {
      size = var.ephemeral_storage_size
    }
  }

  logging_config {
    log_format = "Text"
  }
}

variable "function_name" {
  type = string
}

variable "role_arn" {
  type = string
}

variable "handler" {
  type = string
}

variable "runtime" {
  type = string
}

variable "filename" {
  type = string
}

variable "memory_size" {
  type    = number
  default = 128
}

variable "timeout" {
  type    = number
  default = 3
}

variable "environment" {
  type    = map(string)
  default = {}
}

variable "tracing_mode" {
  type    = string
  default = ""
}

variable "layers" {
  type    = list(string)
  default = []
}

variable "ephemeral_storage_size" {
  type    = number
  default = 512
}

output "arn" {
  value = aws_lambda_function.this.arn
}
