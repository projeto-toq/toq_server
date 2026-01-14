terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

resource "aws_kms_key" "this" {
  description             = var.description
  deletion_window_in_days = var.deletion_window_in_days
  enable_key_rotation     = true
}

resource "aws_kms_alias" "this" {
  name          = var.alias
  target_key_id = aws_kms_key.this.id
}

variable "alias" {
  type        = string
  description = "Alias (with prefix alias/)"
}

variable "description" {
  type        = string
  description = "Key description"
}

variable "deletion_window_in_days" {
  type        = number
  description = "Deletion window"
  default     = 30
}

output "arn" {
  value = aws_kms_key.this.arn
}

output "alias" {
  value = aws_kms_alias.this.name
}
