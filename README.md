# Liquid Web Terraform Provider

## Developing

Dependencies:

- Install [Terraform](https://www.terraform.io)
- Install [Go](https://www.golang.org)
- Install [dep](https://golang.github.io/dep)
- Create a `.lwapi.toml` file in the root directory:

```
[lwApi]
username = "[yourusername]"
password = "[yourpassword]"
url = "https://api.stormondemand.com"
timeout = 15
```

- `make ensure` -- install Golang dependencies.
- `make build` -- build the provider

## Examples

In the `examples` directory there are Terraform projects illustrating how to create various resources. There are a handful of Make tasks that are helpful:

- `EXAMPLE=./examples/storm_servers make key` -- create a new SSH key to provision Storm Servers with (only relevant for the storm server example).
- `EXAMPLE=./examples/storm_servers make devplan` -- build, init and plan.
- `EXAMPLE=./examples/storm_servers make devapply` -- build, init and apply cycle to create resources.
- `EXAMPLE=./examples/storm_servers make destroy` -- destroy resources.

#### Storm Servers

```
resource "liquidweb_storm_server" "api_servers" {
  count = 1

  config_id      = 1090
  zone           = "${data.liquidweb_network_zone.api.id}"
  template       = "UBUNTU_1804_UNMANAGED"                     // ubuntu 18.04
  domain         = "terraform-testing.api.${count.index}.masre.net"
  password       = "11111aA"
  public_ssh_key = "${file("./devkey.pub")}"
}
```

#### Storm Configs

```
data "liquidweb_network_zone" "testing" {
  name        = "Zone C"
  region_name = "US Central"
}

data "liquidweb_storm_server_config" "testing" {
  vcpu         = 2
  memory       = "2000"
  disk         = "100"
  network_zone = "${data.liquidweb_network_zone.testing.id}"
}
```

#### DNS

```
resource "liquidweb_network_dns_record" "testing_a_record" {
  name  = "terraform-testing.api.${count.index}.masre.net"
  type  = "A"
  rdata = "127.0.0.1"
  zone  = "masre.net"
}
```
