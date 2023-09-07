package liquidweb

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider implements the provider definition.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},
		DataSourcesMap: map[string]*schema.Resource{
			"liquidweb_network_zone":        dataSourceLWNetworkZone(),
			"liquidweb_cloud_server_config": dataSourceServerConfig(),

			// backwards compat
			"liquidweb_storm_server_config": dataSourceServerConfig(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"liquidweb_cloud_block_storage":   resourceBlockStorage(),
			"liquidweb_cloud_server":          resourceServer(),
			"liquidweb_network_dns_record":    resourceDNSRecord(),
			"liquidweb_network_load_balancer": resourceLoadBalancer(),
			"liquidweb_network_vip":           resourceFloatingIP(),

			// backwards compat
			"liquidweb_storage_block_volume": resourceBlockStorage(),
			"liquidweb_storm_server":         resourceServer(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return GetConfig()
}
