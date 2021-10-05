variable "auth_token" {
  type = string
}

variable "project" {
  type = string
}

variable "name" {
  description = "The faasd device name."
  type        = string
}

variable "domain" {
  description = "A public domain for the faasd instance. This will the use of Caddy and a Let's Encrypt certificate"
  type        = string
}

variable "email" {
  description = "Email used to order a certificate from Let's Encrypt"
  type        = string
}

variable "metro" {
  description = "Metro area for the new faasd device."
  type        = string
  default     = "am"
}

variable "plan" {
  description = "The faasd device plan slug."
  type        = string
  default     = "c3.small.x86"
}

variable "ufw_enabled" {
  description = "Flag to indicate if UFW (Uncomplicated Firewall) should be enabled or not."
  type        = bool
  default     = true
}

