variable "liquidweb_config_path" {
  type = "string"
}

variable "api_server_password" {
  type = "string"
}

provider "liquidweb" {
  config_path = "${var.liquidweb_config_path}"
}

data "liquidweb_network_zone" "api" {
  region_name = "US West"
}

#
#data "liquidweb_storm_server_config" "api" {
#  active    = true
#  available = true
#  vcpu      = 1
#  memory    = "2G"
#  disk      = "100G"
#  zone      = "${data.storm_network_zone.api_zone.id}"
#}

resource "liquidweb_storm_server" "api_servers" {
  count = 2

  //config_id      = "${data.liquidweb_storm_config.api.id}"
  config_id      = 1090
  template       = "UBUNTU_1804_UNMANAGED"                // ubuntu 18.04
  domain         = "api.${count.index + 1}.mwx.masre.net"
  password       = "${var.api_server_password}"
  public_ssh_key = "${file("./devkey.pub")}"

  //zone           = "${data.liquidweb_network_zone.api.id}"
  zone = 12
}

//
//resource "liquidweb_storm_server" "database_servers" {
//  count = 2
//
//  //config_id      = "${data.liquidweb_storm_config.api.id}"
//  config_id      = 1090
//  template       = "UBUNTU_1804_UNMANAGED"               // ubuntu 18.04
//  domain         = "db.${count.index + 1}.mwx.masre.net"
//  password       = "${var.api_server_password}"
//  public_ssh_key = "${file("./devkey.pub")}"
//
//  //zone           = "${data.liquidweb_network_zone.api.id}"
//  zone = 12
//}
//
//resource "liquidweb_storm_server" "rabbit_servers" {
//  count = 3
//
//  //config_id      = "${data.liquidweb_storm_config.api.id}"
//  config_id      = 1090
//  template       = "UBUNTU_1804_UNMANAGED"                   // ubuntu 18.04
//  domain         = "rabbit.${count.index + 1}.mwx.masre.net"
//  password       = "${var.api_server_password}"
//  public_ssh_key = "${file("./devkey.pub")}"
//
//  //zone           = "${data.liquidweb_network_zone.api.id}"
//  zone = 12
//}
//
//resource "liquidweb_storm_server" "zookeeper_servers" {
//  count = 3
//
//  //config_id      = "${data.liquidweb_storm_config.api.id}"
//  config_id      = 1090
//  template       = "UBUNTU_1804_UNMANAGED"                      // ubuntu 18.04
//  domain         = "zookeeper.${count.index + 1}.mwx.masre.net"
//  password       = "${var.api_server_password}"
//  public_ssh_key = "${file("./devkey.pub")}"
//
//  //zone           = "${data.liquidweb_network_zone.api.id}"
//  zone = 12
//}
//

output "api_network_zone" {
  value = "${data.liquidweb_network_zone.api.name}"
}
