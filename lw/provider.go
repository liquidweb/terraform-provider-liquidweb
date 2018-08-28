package lw

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// Provider implements the provider definition.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"config_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path to the LiquidWeb API configuration file.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"storm_server": resourceServer(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return GetConfig(d.Get("config_path").(string))
}
