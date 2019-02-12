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

resource "liquidweb_storm_server" "testing" {
  count = 1

  config_id      = 1090
  zone           = "${data.liquidweb_network_zone.testing.id}"
  template       = "UBUNTU_1804_UNMANAGED"                          // ubuntu 18.04
  domain         = "terraform-testing.api.${count.index}.masre.net"
  password       = "11111aA"
  public_ssh_key = "${file("${path.module}/devkey.pub")}"
}

resource "liquidweb_network_load_balancer" "testing_some_space_balls" {
  depends_on = ["data.liquidweb_network_zone.testing"]
  name       = "spaceballz44"

  region = "${data.liquidweb_network_zone.testing.region_id}"

  nodes = [
    "${liquidweb_storm_server.testing.ip}",
  ]

  services = [
    {
      src_port  = 80
      dest_port = 80
    },
    {
      src_port  = 1337
      dest_port = 1337
    },
  ]

  #session_persistence = false
  #ssl_termination = false
  strategy = "roundrobin"
}

output "space_balls" {
  value = "${liquidweb_network_load_balancer.testing_some_space_balls.vip}"
}

output "region_id" {
  value = "${data.liquidweb_network_zone.testing.region_id}"
}
