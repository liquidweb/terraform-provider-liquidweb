# Storm Terraform Provider

## Developing

- Install Terraform
- Install Go
- Install [dep](https://golang.github.io/dep)
- Create a `.lwapi.toml` file:

```
[lwApi]
username = "tf_dev"
password = "Sji2zy4Aubs5hDNb"
url = "https://api.stormondemand.com"
timeout = 15
```

- `make key` -- for testing creation of storm instances.
- `make dep` -- install dependencies.
- `make devrun` -- build, init and apply cycle.
