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

data "liquidweb_storm_server_config" "api" {
  vcpu         = 2
  memory       = "2000"
  disk         = "100"
  network_zone = "${data.liquidweb_network_zone.api.id}"
}

resource "liquidweb_storm_server" "api_servers" {
  count = 1

  //config_id      = "${data.liquidweb_storm_server_config.api.id}"
  config_id      = 1090
  zone           = "${data.liquidweb_network_zone.api.id}"
  template       = "UBUNTU_1804_UNMANAGED"                     // ubuntu 18.04
  domain         = "api.dev6.${count.index + 1}.mwx.masre.net"
  password       = "11111aA"
  public_ssh_key = "${file("./devkey.pub")}"
}

resource "liquidweb_network_dns_record" "api_servers" {
  count = "${liquidweb_storm_server.api_servers.count}"

  name  = "api-${count.index}"
  type  = "A"
  rdata = "${element(liquidweb_storm_server.api_servers.*.ip, count.index)}"
  zone  = "masre.net"
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

output "api_server_ips" {
  value = "${join(",", concat(liquidweb_storm_server.api_servers.*.ip))}"
}
