variable "liquidweb_config_path" {
  type = "string"
}

provider "liquidweb" {
  config_path = "${var.liquidweb_config_path}"
}

data "liquidweb_network_zone" "testing" {
  name        = "Zone C"
  region_name = "US Central"
}

data "liquidweb_storm_server_config" "testing" {
  vcpu         = 2
  memory       = "2000"
  disk         = "100"
  network_zone = "${data.liquidweb_network_zone.testing.id}"
}

output "testing_storm_config_id" {
  value = "${data.liquidweb_storm_server_config.testing.id}"
}
