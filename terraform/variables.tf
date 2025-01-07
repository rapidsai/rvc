variable "domain_name" {
  description = "Custom domain name"
  type        = string
  default     = "version.gpuci.io"
}

variable "certificate_name" {
  description = "SSL certificate name"
  type        = string
  default     = "*.gpuci.io"
}
