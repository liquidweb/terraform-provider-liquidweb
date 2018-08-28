# Liquid Web Terraform Provider

## Developing

Dependencies:

- Install [Terraform](https://www.terraform.io)
- Install [Go](https://www.golang.org)
- Install [dep](https://golang.github.io/dep)
- Create a `.lwapi.toml` file:

```
[lwApi]
username = "[yourusername]"
password = "[yourpassword]"
url = "https://api.stormondemand.com"
timeout = 15
```

There is a Terraform definition provided as an example, `main.tf`, which illustrates how to create various resources. It can also be used to actually test resource creation as well.

- `make key` -- create a new SSH key to provision Storm Servers with.
- `make dep` -- install Golang dependencies.
- `make devrun` -- build, init and apply cycle.
