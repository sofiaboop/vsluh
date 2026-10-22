variable "project" {
  description = "Имя продукта, идёт в теги"
  type        = string
  default     = "vsluh"
}

variable "environment" {
  description = "Окружение: dev или prod. Совпадает с именем стейджа API Gateway"
  type        = string
  default     = "dev"
}

variable "domain" {
  description = "Публичный хост API"
  type        = string
  default     = "api-dev.vsluh.club"
}

variable "zone_name" {
  description = "Имя зоны Cloudflare, в которой живёт хост"
  type        = string
  default     = "vsluh.club"
}
