# Simple Wordpress Terraform Example

The files in this directoy provide a simple wordpress terraform deployment.
This deployment is not intended for production - there are no backups, there's a few quirks, but it could be used as-is and makes a good starting point.
This page attempts to explain how the parts fit together.

<!-- vscode-markdown-toc -->
* [Pre-Requisites](#pre-requisites)
* [Files in deployment](#files-in-deployment)
  * [`vars.tf`](#varstf)
  * [`data.tf`](#datatf)
  * [`ssl.tf`](#ssltf)
  * [`server.tf`](#servertf)

<!-- vscode-markdown-toc-config
	numbering=true
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->

## Pre-Requisites

* You must already have a Liquid Web account.
* Your Liquid Web account or API user must be set to environment variables:
  * `LWAPI_USERNAME` - your account username
  * `LWAPI_PASSWORD` - your account username
* For the ACME TLS provider, your account must be set to
  * `LIQUID_WEB_USERNAME` - your account username (for ACME)
  * `LIQUID_WEB_PASSWORD` - your account username (for ACME)
  * `LIQUID_WEB_URL` set to `"https://api.liquidweb.com"`
  * `LIQUID_WEB_ZONE` set to your domain name
* You must have a DNS zone created for the domain you want to use
* You should create a `.tfvars` to change the domain - see below

Once you have those prerequisites, deploy from this directory with:

```bash
terraform init
terraform apply
```

Tear this down with:

```bash
terraform destroy
```

## Files in deployment

This example is made up of 4 files in an attempt to simplify the example.
Technically, this could be in one `.tf` file instead.
Below, the content in each file is explained.

`output` sections in each file simply show things at the end.
They do not impact the deployment, just the output you see.
Relevant useful pieces are shown from each item.

```hcl
output "instances" {
  value = liquidweb_cloud_server.simple_server.*.ip
}
```

### `vars.tf`

This is the first file, and likely the simplest, it just defines some default variables.
The two things that you likely only want to change are:

```hcl
variable "site_name" {
  type = string
  default = "simple.hostbaitor.com"
}

variable "top_domain" {
  type = string
  default = "hostbaitor.com"
}
```

* If you copy these files somewhere else, you can modify those values
* You can change them at cli `terraform apply -var "top_domain=example.com"`
* Environment variables - `export TF_VAR_TOP_DOMAIN="example.com"`
* We'd recommend using a [tfvars file](https://developer.hashicorp.com/terraform/language/values/variables#variable-definitions-tfvars-files)

```hcl
top_domain = "example.com"
site_name = "wordpress.example.com"
```

### `data.tf`

This file is focused on setting up the provider and a lot of pieces of data used elsewhere.

First, the required providers are defined, that section is below.
This is used when you run `terraform init`.

```hcl
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
```

Next a bunch of random passwords and similar are generated:

```hcl
resource "random_password" "server" {
  length  = 20
  special = false
}
```

Finally, templates are rendered.
Templates can be rendered inline, but here they are made into objects so terraform will detect changes.

* `template` is a path to the location of the template (relative to the cwd)
* `vars` is a hash of variables available to be used in the template

```hcl
data "template_file" "wp-config" {
  template = file("${path.module}/templates/wp-config.php")
  vars = {
    dbhost = var.wordpress_dbhost
    dbname = var.wordpress_dbname
    dbuser = var.wordpress_dbuser
    dbpass = random_password.wordpress_dbpass.result
    salt = random_password.wordpress_salt.result
  }
}
```

Finally, there are a few DNS records created for both the site and the server.
These are dependent on the server being created, and happen after that.
But `terraform` automatically creates the DNS records for ease of use.

* `name` is where the DNS record is created
* `rdata` is what the DNS record should point to
* `zone` is the zone to create the DNS record in

All of these are string fields, and can be changed to be whatever else is needed.

```hcl
resource "liquidweb_network_dns_record" "server_dns" {
  name  = liquidweb_cloud_server.simple_server.domain
  type  = "A"
  rdata = liquidweb_cloud_server.simple_server.ip
  zone  = var.top_domain
}
```

### `ssl.tf`

`ssl.tf` incrementally creates the ssl, then outputs the domain name for the SSL.
With this method of SSL generation, you will need to run `terraform apply` to renew.

Relevant links for this section:

* [`acme_certificate` provider](https://registry.terraform.io/providers/vancluever/acme/latest/docs/resources/certificate#using-dns-challenges)
* [`liquidweb` challenge for use with zones hosted with LiquidWeb](https://registry.terraform.io/providers/vancluever/acme/latest/docs/guides/dns-providers-liquidweb)
* [source for the `terraform` provider (uses `lego`)](https://github.com/vancluever/terraform-provider-acme/tree/main/acme)
* [docs for `liquidweb` part within `lego`](https://go-acme.github.io/lego/dns/liquidweb/)
* [source for the `liquidweb` DNS challenge within `lego`](https://github.com/go-acme/lego/tree/master/providers/dns/liquidweb)
* [Let's Encrypt Endpoints](https://letsencrypt.org/docs/acme-protocol-updates/#api-endpoints)

First, the Let's Encrypt URL is configured.

```hcl
provider "acme" {
 server_url = "https://acme-v02.api.letsencrypt.org/directory"
}

Next a private key is created, and used to register with Let's Encrypt:
```hcl
resource "tls_private_key" "private_key" {
  algorithm = "RSA"
}

resource "acme_registration" "reg" {
  account_key_pem = tls_private_key.private_key.private_key_pem
  email_address   = "nobody@${var.site_name}"
}
```

An SSL is requested, and the `liquidweb` DNS challenge is used.

```hcl
resource "acme_certificate" "web_cert" {
  account_key_pem           = acme_registration.reg.account_key_pem
  common_name               = "${var.site_name}"
 key_type = "4096"

# subject_alternative_names = ["www2.example.com"]

  dns_challenge {
    provider = "liquidweb"
  }
}
```

This produces a `acme_certificate.web_cert` resource for use later.

### `server.tf`

The file `server.tf` contains one resource, but it's the most complex one - the server.
This resource depends on random values, templates, the TLS cert, and an SSH key.

First, the server options are specified - this determines how the server is created.

* `zone` (required) determines which zone to create the server in
* `config_id` (required) is the numeric id of the server type to create
* `template` (required) is the base image to use
* `domain` (required) is the server's hostname
* `password` (optional) will be the server's root passwordot provided
* `public_ssh_key` (optional) an SSH key inserted for root
* `lifecycle.create_before_desetroy` determines what to do when recreating a server`

You need either `public_ssh_key` or `password`, but not both.
If no `password` is provide, a random one will be set.

```hcl
  zone      = data.liquidweb_network_zone.zonec.network_zone_id
  config_id = 1757
  template       = "ROCKYLINUX_8_UNMANAGED"
  domain         = "wordpress-host${random_id.server.dec}.us-midwest-2.${var.top_domain}"
  public_ssh_key = file("${path.root}/default.pub")
  password       = random_password.server.result

  lifecycle {
    create_before_destroy = false
  }
```

Next, the connection to that server is configured.
The SSH key at `default.pub` should be replaced with your pub key, then this should use that key.
If you would rather, you can use the root password or [see the docs](https://developer.hashicorp.com/terraform/language/resources/provisioners/connection).

```hcl
  connection {
    type  = "ssh"
    user  = "root"
    agent = true
    host  = self.ip
  }
```

Once the server is online and the connection is good, some commands are run:

```hcl
  provisioner "remote-exec" {
    inline = [
      "yum install -y epel-release",
      "yum install -y http://rpms.remirepo.net/enterprise/remi-release-8.rpm",
      "yum install -y wget curl nginx mysql mysql-common mysql-server php82-php-fpm php82-php-mysqlnd php82-php-mbstring"
    ]
  }
```

Some files are also written to the vm.
All files here are from templates, but direct files could be done as well.
But see the [`file provisioner`](https://developer.hashicorp.com/terraform/language/resources/provisioners/file) docs and you can use a straight file instead as well.

```hcl
  provisioner "file" {
    content     = data.template_file.site-conf.rendered
    destination = "/etc/nginx/conf.d/site.conf"
```

Once all of those are done, the server's up and ready to go through the Wordpress setup.
