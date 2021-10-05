terraform {
  required_providers {
    metal = {
      source = "equinix/metal"
    }
  }
}

# Configure the Equinix Metal Provider.
provider "metal" {
  auth_token = var.auth_token
}

data "metal_project" "project" {
  name = var.project
}

module "faasd" {
  source = "github.com/jsiebens/terraform-equinix-faasd"

  project_id = data.metal_project.project.id
  name       = var.name
  domain     = var.domain
  email      = var.email

  plan                = var.plan
  metro               = var.metro
  ufw_enabled         = var.ufw_enabled
  project_ssh_key_ids = []
}

