terraform {
  required_providers {
    liquidweb = {
      source  = "liquidweb/liquidweb"
      version = ">= 1.7.0"
    }
    acme = {
      source = "vancluever/acme"
      version = "2.17.1"
    }
  }
}

resource "random_id" "server" {
  byte_length = 1
  # count       = 1
}

resource "random_password" "server" {
  length  = 20
  special = false
}

resource "random_password" "wordpress_dbpass" {
  length  = 20
  special = true
}

resource "random_password" "wordpress_salt" {
  length = 32
  special = true
}

data "liquidweb_network_zone" "zonec" {
  name        = "Zone C"
  region_name = "US Central"
}

data "template_file" "install-wordpress" {
  template = file("${path.module}/templates/install-wordpress.sh")
  vars = {
    user = var.username
  }
}

data "template_file" "wp-config" {
  template = file("${path.module}/templates/wp-config.php")
  vars = {
    dbname = var.wordpress_dbname
    dbuser = var.wordpress_dbuser
    dbpass = random_password.wordpress_dbpass.result
    salt = random_password.wordpress_salt.result
    dbhost = liquidweb_cloud_server.dbserver.ip
  }
}

data "template_file" "site-conf" {
  template = file("${path.module}/templates/nginx.conf")
  vars = {
    domain = var.site_name
    user = var.username
  }
}

data "template_file" "php-conf" {
  template = file("${path.module}/templates/php-fpm.conf")
  vars = {
    user = var.username
  }
}

resource "liquidweb_network_dns_record" "webserver_dns" {
  count = 3
  name  = liquidweb_cloud_server.webserver[count.index].domain
  type  = "A"
  rdata = liquidweb_cloud_server.webserver[count.index].ip
  zone  = var.top_domain
}

resource "liquidweb_network_dns_record" "wordpress_record" {
  name  = var.site_name
  type  = "A"
  rdata = liquidweb_network_load_balancer.loadbalancer.vip
  zone  = var.top_domain
}

output "domain_a_name" {
  value = liquidweb_network_dns_record.wordpress_record.name
}
