variable "liquidweb_config_path" {
  type = "string"
}

provider "liquidweb" {
  config_path = "${var.liquidweb_config_path}"
}

resource "liquidweb_network_vip" "new_vip" {
  domain  = "terraform-testing-vip"
  zone    = 28
}

output "vip_name" {
  value = "${liquidweb_network_vip.new_vip.domain}"
}

output "vip_ip" {
  value = "${liquidweb_network_vip.new_vip.ip}"
}