variable "liquidweb_config_path" {
  type = string
}

terraform {
  required_providers {
    liquidweb = {
      source = "registry.terraform.io/liquidweb/liquidweb"
      version = "~> 1.6.2"
    }
  }
}

provider "liquidweb" {
  config_path = var.liquidweb_config_path
}