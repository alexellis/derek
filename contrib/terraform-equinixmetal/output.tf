output "summary" {
  value = <<CONFIGURATION

Your faasd instance "${var.name} is ready.

    IP Address: ${module.faasd.ipv4_address}
    
To continue, use the IP address above create a DNS A record for your domain "${var.domain}"
Give Caddy a moment to get a certificate and when ready, the faasd gateway is available at:

    ${module.faasd.gateway_url}

Authenticate with faas-cli:

    export PASSWORD=$(terraform output -raw basic_auth_password)
    export OPENFAAS_URL=${module.faasd.gateway_url}
    echo $PASSWORD | faas-cli login -u admin --password-stdin

CONFIGURATION
}

output "basic_auth_password" {
  description = "The basic auth password."
  value       = module.faasd.basic_auth_password
  sensitive   = true
}

