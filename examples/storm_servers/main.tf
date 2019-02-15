variable "liquidweb_config_path" {
  type = "string"
}

provider "liquidweb" {
  config_path = "${var.liquidweb_config_path}"
}

data "liquidweb_network_zone" "api" {
  name        = "Zone C"
  region_name = "US Central"
}

resource "liquidweb_storm_server" "api_servers" {
  count = 1

  //config_id      = "${data.liquidweb_storm_server_config.api.id}"
  config_id      = 1090
  zone           = "${data.liquidweb_network_zone.api.id}"
  template       = "UBUNTU_1804_UNMANAGED"                            // ubuntu 18.04
  domain         = "terraform-testing.2.api.${count.index}.masre.net"
  password       = "11111aA"
  public_ssh_key = "${file("${path.module}/devkey.pub")}"
}

output "api_server_ips" {
  value = "${join(",", concat(liquidweb_storm_server.api_servers.*.ip))}"
}
