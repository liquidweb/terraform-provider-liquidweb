variable "liquidweb_config_path" {
  type = "string"
}

variable "api_server_password" {
  type = "string"
}

provider "storm" {
  config_path = "${var.liquidweb_config_path}"
}

//data "liquidweb_network_zone" "api" {
//  active    = true
//  available = true
//  vcpu      = 1
//  memory    = "2G"
//  disk      = "100G"
//  zone      = 12
//}
//
//data "liquidweb_storm_server_config" "api" {
//  active    = true
//  available = true
//  vcpu      = 1
//  memory    = "2G"
//  disk      = "100G"
//  zone      = "${data.storm_network_zone.api_zone.id}"
//}

resource "liquidweb_storm_server" "api_servers" {
  count          = 2
  config_id      = "${data.liquidweb_storm_config.api.id}"
  template       = "UBUNTU_1804_UNMANAGED"                 // ubuntu 18.04
  domain         = "api.${count.index + 1}.mwx.masre.net"
  password       = "${var.api_server_password}"
  public_ssh_key = "${file("./devkey.pub")}"
  zone           = "${data.liquidweb_network_zone.api.id}"
}
