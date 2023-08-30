package liquidweb

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"liquidweb_network_zone":        dataSourceLWNetworkZone(),
			"liquidweb_cloud_server_config": dataSourceLWStormServerConfig(),

			// backwards compat
			"liquidweb_storm_server_config": dataSourceLWStormServerConfig(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"liquidweb_cloud_block_storage":   resourceStorageBlockVolume(),
			"liquidweb_cloud_server":          resourceStormServer(),
			"liquidweb_network_dns_record":    resourceNetworkDNSRecord(),
			"liquidweb_network_load_balancer": resourceNetworkLoadBalancer(),
			"liquidweb_network_vip":           resourceNetworkVIP(),

			// backwards compat
			"liquidweb_storage_block_volume": resourceStorageBlockVolume(),
			"liquidweb_storm_server":         resourceStormServer(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return GetConfig(d.Get("config_path").(string))
}
