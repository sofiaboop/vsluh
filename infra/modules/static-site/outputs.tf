output "bucket" {
  description = "Имя бакета — цель `aws s3 sync` в CI"
  value       = aws_s3_bucket.site.id
}

output "distribution_id" {
  description = "Id дистрибуции — цель инвалидации в CI"
  value       = aws_cloudfront_distribution.site.id
}

output "distribution_domain" {
  description = "Домен CloudFront, на который смотрит CNAME в Cloudflare"
  value       = aws_cloudfront_distribution.site.domain_name
}
