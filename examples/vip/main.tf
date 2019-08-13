variable "liquidweb_config_path" {
  type = "string"
}

provider "liquidweb" {
  config_path = var.liquidweb_config_path
}

resource "liquidweb_network_vip" "new_vip" {
  count  = 5
  domain = "terraform-testing-vip-${count.index + 1}"
  zone   = 28
}

output "vip_names" {
  value = join(",", liquidweb_network_vip.new_vip.*.domain)
}

output "vip_ips" {
  value = join(",", liquidweb_network_vip.new_vip.*.ip)
}
