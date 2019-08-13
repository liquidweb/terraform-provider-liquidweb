variable "liquidweb_config_path" {
  type = "string"
}

provider "liquidweb" {
  config_path = "${var.liquidweb_config_path}"
}

resource "liquidweb_network_dns_record" "testing_a_record" {
  name  = "terraform-testing.api.lwmasre.net"
  type  = "A"
  rdata = "127.0.0.1"
  zone  = "lwmasre.net"
}

output "api_server_a_name" {
  value = "${liquidweb_network_dns_record.testing_a_record.name}"
}
