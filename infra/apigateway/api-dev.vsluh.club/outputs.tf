output "domain" {
  value = aws_api_gateway_domain_name.api.domain_name
}

output "regional_domain" {
  description = "Хост API Gateway, на который смотрит CNAME в Cloudflare"
  value       = aws_api_gateway_domain_name.api.regional_domain_name
}

output "certificate_arn" {
  value = aws_acm_certificate.api.arn
}
