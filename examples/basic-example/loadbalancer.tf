resource "liquidweb_network_load_balancer" "testing_lb" {
  depends_on = [data.liquidweb_network_zone.testing_zone]
  name       = "lb.0.terraform-testing.api.example.com"

  region = data.liquidweb_network_zone.testing_zone.region_id

  nodes = liquidweb_cloud_server.testing_servers[*].ip

  service {
    src_port  = 80
    dest_port = 80
  }

  service {
    src_port  = 1337
    dest_port = 1337
  }

  #session_persistence = false
  #ssl_termination = false
  strategy = "roundrobin"
}

output "lb_vip" {
  value = liquidweb_network_load_balancer.testing_lb.vip
}