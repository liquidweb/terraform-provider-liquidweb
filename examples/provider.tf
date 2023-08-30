variable "liquidweb_config_path" {
  type = string
}

terraform {
  required_providers {
    liquidweb = {
      source = "local.providers/liquidweb/liquidweb"
      version = "~> 1.5.8"
    }
  }
}

provider "liquidweb" {
  config_path = var.liquidweb_config_path
}