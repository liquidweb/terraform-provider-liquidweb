resource "random_id" "server" {
  byte_length = 1
  count = 2
}

resource "liquidweb_cloud_server" "testing_servers" {
  count = 2

  #config_id = "${data.liquidweb_cloud_server_config.api.id}"
  config_id = 1757
  zone      = 27
  #data.liquidweb_network_zone.api.id
  template       = "UBUNTU_1804_UNMANAGED" // ubuntu 18.04
  domain         = "terraform-host${random_id.server[count.index].dec}.us-midwest-2.hostbaitor.com"
  public_ssh_key = file("${path.root}/devkey.pub")
  password       = "1Aaaaaaaaa"
}

output "instances" {
  value = liquidweb_cloud_server.testing_servers.*.ip
}
