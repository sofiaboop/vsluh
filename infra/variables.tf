variable "project_name" {
  description = "Префикс имён ресурсов"
  type        = string
  default     = "vsluh"
}

variable "environment" {
  description = "Окружение: dev или prod"
  type        = string

  validation {
    condition     = contains(["dev", "prod"], var.environment)
    error_message = "environment должен быть dev или prod."
  }
}

variable "aws_region" {
  description = "Регион AWS для бакетов"
  type        = string
  default     = "eu-central-1"
}

variable "root_domain" {
  description = "Корневой домен зоны Cloudflare"
  type        = string
  default     = "vsluh.club"
}

variable "cloudflare_zone_id" {
  description = "ID зоны Cloudflare для root_domain"
  type        = string
}
