terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

resource "aws_iam_role" "this" {
  name               = var.name
  assume_role_policy = var.assume_role_policy
}

resource "aws_iam_role_policy" "inline" {
  count  = length(var.inline_policies)
  name   = var.inline_policies[count.index].name
  role   = aws_iam_role.this.id
  policy = var.inline_policies[count.index].document
}

resource "aws_iam_role_policy_attachment" "managed" {
  count      = length(var.managed_policy_arns)
  role       = aws_iam_role.this.name
  policy_arn = var.managed_policy_arns[count.index]
}

variable "name" {
  type        = string
  description = "Role name"
}

variable "assume_role_policy" {
  type        = string
  description = "Assume role policy JSON"
}

variable "inline_policies" {
  description = "List of inline policies"
  type = list(object({
    name      = string
    document  = string
  }))
  default = []
}

variable "managed_policy_arns" {
  description = "List of managed policy arns"
  type        = list(string)
  default     = []
}

output "arn" {
  value = aws_iam_role.this.arn
}

output "name" {
  value = aws_iam_role.this.name
}
