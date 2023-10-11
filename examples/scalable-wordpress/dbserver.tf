resource "liquidweb_cloud_server" "dbserver" {
  #config_id = "${data.liquidweb_storm_server_config.api.id}"
  config_id = 1757
  zone      = data.liquidweb_network_zone.zonec.network_zone_id
  #data.liquidweb_network_zone.api.id
  template       = "ROCKYLINUX_8_UNMANAGED"
  domain         = "wordpress-db01-p${random_id.server.dec}.us-midwest-2.${var.top_domain}"
  public_ssh_key = file("${path.root}/default.pub")
  password       = random_password.server.result

  connection {
    type  = "ssh"
    user  = "root"
    agent = true
    host  = self.ip
  }

  provisioner "remote-exec" {
    inline = [
      "yum install -y epel-release",
      "yum install -y http://rpms.remirepo.net/enterprise/remi-release-8.rpm",
      "yum install -y wget curl mysql mysql-common mysql-server"
    ]
  }

  provisioner "remote-exec" {
    inline = [
      "systemctl start mysqld.service",
      "systemctl enable mysqld.service",
      "firewall-cmd --zone public --permanent --add-port 3306/tcp",
      "firewall-cmd --reload"
    ]
  }
}

output "dbserver_hostnames" {
  value = liquidweb_cloud_server.dbserver.ip
}

output "dbserver_ips" {
  value = liquidweb_cloud_server.dbserver.ip
}

