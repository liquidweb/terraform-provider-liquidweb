resource "liquidweb_network_load_balancer" "loadbalancer" {
  # depends_on = [data.liquidweb_network_zone.testing_zone]
  name       = "wordpress-loadbalancer1-p${random_id.server.dec}${var.top_domain}"

  region = data.liquidweb_network_zone.zonec.region_id

  nodes = liquidweb_cloud_server.webserver[*].ip

  service {
    src_port  = 80
    dest_port = 80
  }

  service {
    src_port  = 443
    dest_port = 443
  }

  #session_persistence = false
  #ssl_termination = false
  strategy = "roundrobin"
}

output "lb_vip" {
  value = liquidweb_network_load_balancer.loadbalancer.vip
}