# Liquid Web Terraform Provider

[![Build Status](https://travis-ci.org/liquidweb/terraform-provider-liquidweb.svg?branch=master)](https://travis-ci.org/liquidweb/terraform-provider-liquidweb)

## Developing

If you want to develop for this, you need go1.21+ and `terraform` 1.5+.

Tests can be run with (from the root of this repository):

```bash
go test -v ./...
```

This can be built with:

```bash
go build
```

Then to use a bin you built, relative to your `terraform minfests` copy it to:

```bash
~/terraform.d/plugins/local.providers/liquidweb/liquidweb/1.6.0/darwin_amd64/terraform-provider-liquidweb
```

- `darwin_amd64` changes relevative to your platform
- `1.6.0` is relative to what version this is`

## Tracing

Tracing via Jaeger is available so various actions: successful API patterns, bottlenecks and problems can be identified recognized accordingly. It's important to capture where we're getting things right as much as wrong.

```shell
make jaeger
xdg-open http://localhost:16686/search
```

Tracing is enabled if `JAEGER_DISABLED` is set to `false`. This requires the `jaeger` container to be running and general use with an external Terraform project isn't yet supported.

## Using this Provider

All things for this section are relative to your terraform manifests.

This terraform provider is currently not published.
Thus, you need the provider compiled above at the location above.

You also currently need the following toml file with your credentials.

- Create a `.lwapi.toml` file in the root directory:

```toml
[lwApi]
username = "[yourusername]"
password = "[yourpassword]"
url = "https://api.liquidweb.com"
timeout = 15
```

These both then need to be included in with a provider block like:

```tf
terraform {
  required_providers {
    liquidweb = {
      source = "local.providers/liquidweb/liquidweb"
      version = "~> 1.5.8"
    }
  }
}

variable "liquidweb_config_path" {
  type = string
}

provider "liquidweb" {
  config_path = var.liquidweb_config_path
}
```

## Examples

In the `examples` directory there are Terraform manifests demonstrating usage.
Everything needs `provider.tf`, most others are only dependent on themselves.

There are also a few example projects in that folder.

### Cloud Servers

```terraform
data "liquidweb_network_zone" "testing" {
  name        = "Zone C"
  region_name = "US Central"
}

resource "liquidweb_cloud_server" "testing" {
  count = 1

  config_id      = 1090
  zone           = data.liquidweb_network_zone.testing.id
  template       = "UBUNTU_1804_UNMANAGED"                     // ubuntu 18.04
  domain         = "terraform-testing.api.${count.index}.masre.net"
  password       = "11111aA"
  public_ssh_key = file("./devkey.pub")
}
```

### Cloud Servers + Load Balancer

```terraform
data "liquidweb_network_zone" "testing" {
  name        = "Zone C"
  region_name = "US Central"
}

resource "" "testing" {
  count = 1

  config_id      = 1090
  zone           = data.liquidweb_network_zone.testing.id
  template       = "UBUNTU_1804_UNMANAGED"                     // ubuntu 18.04
  domain         = "terraform-testing.api.${count.index}.masre.net"
  password       = "11111aA"
  public_ssh_key = file("./devkey.pub")
}

resource "liquidweb_network_load_balancer" "testing" {
  name       = "testing"
  region = data.liquidweb_network_zone.testing.region_id

  nodes = .testing[*].ip

  service {
    src_port  = 80
    dest_port = 80
  }

  service {
    src_port  = 1337
    dest_port = 1337
  }

  strategy = "roundrobin"
}
```

### Cloud Configs

```terraform
data "liquidweb_network_zone" "testing" {
  name        = "Zone C"
  region_name = "US Central"
}

data "_config" "testing" {
  vcpu         = 2
  memory       = "2000"
  disk         = "100"
  network_zone = data.liquidweb_network_zone.testing.id
}
```

### DNS

```terraform
resource "liquidweb_network_dns_record" "testing" {
  name  = "terraform-testing.api.${count.index}.masre.net"
  type  = "A"
  rdata = "127.0.0.1"
  zone  = "masre.net"
}
```

### Block Volumes

```terraform
resource "liquidweb_network_block_volume" "testing" {
  attach = "2GHUN4"
  domain = "blarstacoman"
  size   = 10
}
```

### VIP

```terraform
resource "liquidweb_network_vip" "testing" {
  domain  = "terraform-testing-vip"
  zone    = 52
}
```
