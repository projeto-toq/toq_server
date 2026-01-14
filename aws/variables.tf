variable "aws_region" {
  type        = string
  description = "AWS region"
  default     = "us-east-1"
}

variable "ffmpeg_layer_arn" {
  type        = string
  description = "Existing ffmpeg layer ARN for video thumbnails"
  default     = "arn:aws:lambda:us-east-1:058264253741:layer:toq-ffmpeg-layer-staging:5"
}

variable "callback_secret" {
  type        = string
  description = "Secret for callback dispatch lambda"
  default     = "toq-media-callback-secret-staging-2024"
}

variable "callback_url" {
  type        = string
  description = "Callback URL"
  default     = "https://api.gca.dev.br/api/v2/listings/media/callback"
}

variable "internal_api_token" {
  type        = string
  description = "Internal API token"
  default     = "toq-internal-token-staging-change-me"
}
