resource "liquidweb_cloud_server" "webserver" {
  count = 3

  #config_id = "${data.liquidweb_storm_server_config.api.id}"
  config_id = 1757
  zone      = data.liquidweb_network_zone.zonec.network_zone_id
  #data.liquidweb_network_zone.api.id
  template       = "ROCKYLINUX_8_UNMANAGED"
  domain         = "wordpress-webserver${count.index}-p${random_id.server.dec}.us-midwest-2.${var.top_domain}"
  public_ssh_key = file("${path.root}/default.pub")
  password       = random_password.server.result

  lifecycle {
    create_before_destroy = false
  }

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
      "yum install -y wget curl nginx mysql mysql-common php82-php-fpm php82-php-mysqlnd php82-php-mbstring"
    ]
  }

  provisioner "file" {
    content     = "${acme_certificate.web_cert.certificate_pem}${acme_certificate.web_cert.issuer_pem}"
    destination = "/etc/pki/tls/certs/${var.site_name}.crt"
  }
  provisioner "file" {
    content     = "${acme_certificate.web_cert.private_key_pem}"
    destination = "/etc/pki/tls/private/${var.site_name}.key"
  }

  provisioner "file" {
    content     = data.template_file.site-conf.rendered
    destination = "/etc/nginx/conf.d/site.conf"
  }
  provisioner "file" {
    content     = data.template_file.php-conf.rendered
    destination = "/etc/opt/remi/php82/php-fpm.d/site.conf"
  }
  provisioner "file" {
    content     = data.template_file.install-wordpress.rendered
    destination = "/root/install-wordpress.sh"
  }
  provisioner "file" {
    content     = data.template_file.wp-config.rendered
    destination = "/root/wp-config.php"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /root/install-wordpress.sh",
      "/root/install-wordpress.sh",
      "systemctl enable nginx.service php82-php-fpm.service",
      "systemctl start nginx.service php82-php-fpm.service",
      "firewall-cmd --zone public --permanent --add-port 80/tcp",
      "firewall-cmd --zone public --permanent --add-port 443/tcp",
      "firewall-cmd --reload"
    ]
  }
}

output "webserver_hostnames" {
  value = liquidweb_cloud_server.webserver.*.domain
}

output "webserver_ips" {
  value = liquidweb_cloud_server.webserver.*.ip
}

