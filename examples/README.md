# Terraform Examples

This repository contains terraform examples that are slated for inclusion in the terraform repository elsewhere, but not yet ready.
I will likely archive / remove it in the future when they are at a place where they can be merged into `liquidweb/terraform-provider-liquidweb`.

## What is Terraform?

Terraform is a tool targeted at a Infrastructure as Code approach to managing asset inventory.
It offers a declarative language to create configurations describing infrastructure.
Given the configurations, you can rapidly create, remove, and recreate infrastructure.
Since the configuraitons are plaintext, it allows easy versioning of the infrastructure state with Version Control Software.

### Background Terms

That's a loaded paragraph, some terminology:

- Version Control Software (vcs) - like `git`, a tool that lets you track files over time and compare differences
  - Of note, most developers put their source code in VCS. This has many benefits.
- Infrastructure as Code (IaC) - managing servers via config files, often which you can commit to a repository
- Declarative Syntax / Language - describing what an system should be
- Asset Inventory - what assets you have. VPS's are an asset, but SSL certificates, LB's, and Block Storage are also assets.
- configuration files end in `.tf` and determine what is needed
- State - the current way a system is, the actual live snapshot of it, not the way it hsould be
- Lockfile - a file tracking what things terraform currently has

### Terraform Basic Commands

The focus of Terraform is create, recreate, and destroying what is needed.
Terraform can be used alone, and assets recreated as your schema changes.
But most of the time, multiple IaC tools are used to better describe a system.

The major background pieces it will create are:

- the lock file resides at `./.terraform.lock.hcl`
- the state file resides at `./terraform.tfstate`
- a backup state file at `./terraform.tfstate.backup`
- providers typically reside in `./terraform.d`

If you have something deployed, you want to save the

The major commands that terraform provides are:

- `init` - download required providers and set up state and lockfile
- `validate` - make sure configs are valid
- `plan` - show changes to modify state to match configs
- `apply` - run `plan`, then prompt to make those changes
- `destroy` show changes to remove everything, prompt, then remove everything
- `show` - display the current assets
- `taint` - mark an asset currently deployed, on next `apply` will be recreated
- `refresh` - update the state of assets (not supported with LiquidWeb's provider)
- `import` - add existing assets into current state (not supported with LiquidWeb's provider)

Terraform Modules also exist.
For example, [there is a wordpress module](https://registry.terraform.io/providers/yyamanoi1222/wordpress/latest/docs).
There are also other IaC tools that can be used with Terraform.
The idomatic approach to IaC typically involves multiple tools used in combination.
However for the purposes of these examples, the focus is on terraform.

### Installing and Examples

The [Hashicorp official instructions for installing terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli) are well written.
You are likely best off going there.
Do note, you will likely have a better time if you use the package install for your OS.
In other words, on macOS use `homebrew`, on Windows use `choco`, on Linux use your package manager.

The LiquidWeb provider does not require special installation.
If it is used in your configs, it should be automatically installed with `terraform init`.

[Documentation for that provider is published on the provider page](https://registry.terraform.io/providers/liquidweb/liquidweb/latest/docs).

You will need to provide credentials to a LiquidWeb account in order to use the LiquidWeb provider.
These credentials should be in the following environment variables:

```env
LWAPI_USERNAME="username"
LWAPI_PASSWORD="password"
```

There is also an `acme` SSL provider.
If your domain is hosted with LiquidWeb, you can use this to get an SSL.
[The documentation gives a basic example](https://registry.terraform.io/providers/vancluever/acme/latest/docs/guides/dns-providers-liquidweb).
If you wish to get an SSL with the `acme` provider with a DNS server, you must provide the following credentials:

```env
LIQUID_WEB_USERNAME="username"
LIQUID_WEB_PASSWORD="password"
LIQUID_WEB_ZONE="zone" # the DNS zone where the ACME challenge record must be created
```

For examples, please look at:

- [Basic server example](https://github.com/jakdept/liquidweb-terraform-example/tree/main/basic-example)
- [Basic wordpress example](https://github.com/jakdept/liquidweb-terraform-example/tree/main/simple-wordpress)
- [Wordpress Cluster example](https://github.com/jakdept/liquidweb-terraform-example/tree/main/scalable-wordpress)

Eventually, these will be moving into [the repository for the provider](https://github.com/liquidweb/terraform-provider-liquidweb/tree/master/examples)
and will be automatically published to other locations.

More examples will also be available.
More features will also likely be added to the provider in the future as well.
Requests for specific examples will also be opened in the future.
However, first, the focus is on simplifying some of the tools.
