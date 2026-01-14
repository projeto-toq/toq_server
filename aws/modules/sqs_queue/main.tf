terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

resource "aws_sqs_queue" "this" {
  name                              = var.name
  visibility_timeout_seconds        = var.visibility_timeout_seconds
  message_retention_seconds         = var.message_retention_seconds
  receive_wait_time_seconds         = var.receive_wait_time_seconds
  kms_master_key_id                 = var.kms_master_key_id
  kms_data_key_reuse_period_seconds = var.kms_data_key_reuse_period_seconds
  delay_seconds                     = var.delay_seconds
  max_message_size                  = var.max_message_size
  redrive_policy                    = var.redrive_policy
}

variable "name" {
  type        = string
  description = "Queue name"
}

variable "visibility_timeout_seconds" {
  type        = number
  description = "Visibility timeout"
  default     = 60
}

variable "message_retention_seconds" {
  type        = number
  description = "Retention period"
  default     = 345600
}

variable "receive_wait_time_seconds" {
  type        = number
  description = "Long poll wait time"
  default     = 20
}

variable "kms_master_key_id" {
  type        = string
  description = "KMS key alias or ARN"
  default     = ""
}

variable "kms_data_key_reuse_period_seconds" {
  type        = number
  description = "Data key reuse period"
  default     = 300
}

variable "delay_seconds" {
  type        = number
  description = "Queue delay"
  default     = 0
}

variable "max_message_size" {
  type        = number
  description = "Max message size"
  default     = 262144
}

variable "redrive_policy" {
  type        = string
  description = "JSON redrive policy"
  default     = null
}

output "arn" {
  value = aws_sqs_queue.this.arn
}

output "url" {
  value = aws_sqs_queue.this.id
}
