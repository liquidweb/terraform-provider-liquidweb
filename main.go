package main

import (
	"git.liquidweb.com/masre/terraform-provider-liquidweb/liquidweb"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return liquidweb.Provider()
		},
	})
}
