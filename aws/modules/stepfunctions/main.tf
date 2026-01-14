terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

resource "aws_sfn_state_machine" "this" {
  name     = var.name
  role_arn = var.role_arn
  type     = "STANDARD"
  definition = file(var.definition_path)

  logging_configuration {
    level                  = var.logging_level
    include_execution_data = var.logging_include_execution_data
    log_destination        = var.logging_log_group_arn != "" ? var.logging_log_group_arn : null
  }

  tracing_configuration {
    enabled = var.tracing_enabled
  }
}

variable "name" {
  type        = string
  description = "State machine name"
}

variable "role_arn" {
  type        = string
  description = "Execution role ARN"
}

variable "definition_path" {
  type        = string
  description = "Path to JSON definition"
}

variable "logging_level" {
  type        = string
  description = "Logging level"
  default     = "OFF"
}

variable "logging_include_execution_data" {
  type        = bool
  description = "Include execution data"
  default     = false
}

variable "logging_log_group_arn" {
  type        = string
  description = "Log group ARN"
  default     = ""
}

variable "tracing_enabled" {
  type        = bool
  description = "X-Ray tracing enabled"
  default     = false
}

output "arn" {
  value = aws_sfn_state_machine.this.arn
}
