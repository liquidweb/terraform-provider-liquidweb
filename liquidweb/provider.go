package liquidweb

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
		DataSourcesMap: map[string]*schema.Resource{
			"liquidweb_network_zone": dataSourceLWNetworkZone(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"liquidweb_storm_server": resourceStormServer(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return GetConfig(d.Get("config_path").(string))
}
