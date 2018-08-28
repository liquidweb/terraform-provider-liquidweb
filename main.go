package main

import (
	"git.liquidweb.com/masre/terraform-provider-liquidweb/lw"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return lw.Provider()
		},
	})
}
