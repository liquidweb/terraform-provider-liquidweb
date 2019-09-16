[![Build Status](https://travis-ci.org/liquidweb/terraform-provider-liquidweb.svg?branch=master)](https://travis-ci.org/liquidweb/terraform-provider-liquidweb)

# Liquid Web Terraform Provider

## Developing

Dependencies:

- Create a `.lwapi.toml` file in the root directory:

```toml
[lwApi]
username = "[yourusername]"
password = "[yourpassword]"
url = "https://api.stormondemand.com"
timeout = 15
```

- `make shell` -- drop into a development shell so you can build/test the provider.

The following run inside the development shell:

- `make build` -- build the provider
- `make init` -- initialize terraform
- `EXAMPLE=./examples/storm_servers make plan` -- plan an example project
- `EXAMPLE=./examples/storm_servers make apply` -- apply an example project
- `EXAMPLE=./examples/storm_servers make destroy` -- destroy an example project
- `make test_release` -- test a release (requires goreleaser to be installed)

There are also `devplan` and `devapply` make tasks that will do a build and subsequent init followed by plan/apply.

## Tracing

Tracing via Jaeger is available so various actions: successful API patterns, bottlenecks and problems can be identified recognized accordingly. It's important to capture where we're getting things right as much as wrong.

```shell
make jaeger
xdg-open http://localhost:16686/search
```

In the development container:

```shell
JAEGER_DISABLED=false EXAMPLE=./examples/storm_servers make apply
```

Tracing is enabled if `JAEGER_DISABLED` is set to `false`. This requires the `jaeger` container to be running and general use with an external Terraform project isn't yet supported.

## Examples

In the `examples` directory there are Terraform projects illustrating how to create various resources. There are a handful of Make tasks that are helpful:

- `EXAMPLE=./examples/storm_servers make key` -- create a new SSH key to provision Storm Servers with (only relevant for the storm server and load balancer example).
- `EXAMPLE=./examples/storm_servers make devplan` -- build, init and plan.
- `EXAMPLE=./examples/storm_servers make devapply` -- build, init and apply cycle to create resources.
- `EXAMPLE=./examples/storm_servers make destroy` -- destroy resources.

#### Storm Servers

```terraform
data "liquidweb_network_zone" "testing" {
  name        = "Zone C"
  region_name = "US Central"
}

resource "liquidweb_storm_server" "testing" {
  count = 1

  config_id      = 1090
  zone           = data.liquidweb_network_zone.testing.id
  template       = "UBUNTU_1804_UNMANAGED"                     // ubuntu 18.04
  domain         = "terraform-testing.api.${count.index}.masre.net"
  password       = "11111aA"
  public_ssh_key = file("./devkey.pub")
}
```

#### Storm Servers + Load Balancer

```terraform
data "liquidweb_network_zone" "testing" {
  name        = "Zone C"
  region_name = "US Central"
}

resource "liquidweb_storm_server" "testing" {
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

  nodes = liquidweb_storm_server.testing[*].ip

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

#### Storm Configs

```terraform
data "liquidweb_network_zone" "testing" {
  name        = "Zone C"
  region_name = "US Central"
}

data "liquidweb_storm_server_config" "testing" {
  vcpu         = 2
  memory       = "2000"
  disk         = "100"
  network_zone = data.liquidweb_network_zone.testing.id
}
```

#### DNS

```terraform
resource "liquidweb_network_dns_record" "testing" {
  name  = "terraform-testing.api.${count.index}.masre.net"
  type  = "A"
  rdata = "127.0.0.1"
  zone  = "masre.net"
}
```

#### Block Volumes

```terraform
resource "liquidweb_network_block_volume" "testing" {
  attach = "2GHUN4"
  domain = "blarstacoman"
  size   = 10
}
```

#### VIP

```terraform
resource "liquidweb_network_vip" "testing" {
  domain  = "terraform-testing-vip"
  zone    = 52
}
```
