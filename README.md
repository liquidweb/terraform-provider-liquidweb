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

Since we're using a private repo for `liquidweb-go`, our API client for LW APIs, `dep` will attempt to clone via https which will fail since terminal prompts are disabled during the "ensure process". There is likely a better way to solve this but for now:

```
git config --global url."git@git.liquidweb.com:".insteadOf "https://git.liquidweb.com/"
```

- `make key` -- create a new SSH key to provision Storm Servers with.
- `make dep` -- install Golang dependencies.
- `make devrun` -- build, init and apply cycle.
