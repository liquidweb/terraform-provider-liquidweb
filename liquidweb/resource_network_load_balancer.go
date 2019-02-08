package liquidweb

import (
	network "git.liquidweb.com/masre/liquidweb-go/network"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNetworkLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateNetworkLoadBalancer,
		Read:   resourceReadNetworkLoadBalancer,
		Update: resourceUpdateNetworkLoadBalancer,
		Delete: resourceDeleteNetworkLoadBalancer,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"nodes": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     schema.TypeString,
			},
			"region": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"services": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"src_port": {
							Type:     schema.TypeString,
							Required: true,
						},
						"dest_port": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"session_persistence": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ssl_cert": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssl_includes": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ssl_int": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssl_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssl_termination": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"strategy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"vip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCreateNetworkLoadBalancer(d *schema.ResourceData, m interface{}) error {
	opts := buildNetworkLoadBalancerOpts(d, m)
	config := m.(*Config)

	result, err := config.LWAPI.NetworkLoadBalancer.Create(opts)
	if err != nil {
		return err
	}

	if result.HasError() {
		return result
	}

	d.SetId(result.UniqID)

	return resourceReadNetworkLoadBalancer(d, m)
}

func resourceReadNetworkLoadBalancer(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	loadBalancerItem := loadBalancerDetails(config, d.Id())
	if loadBalancerItem.HasError() {
		return loadBalancerItem
	}

	updateLoadBalancerResource(d, loadBalancerItem)
	return nil
}

func resourceUpdateNetworkLoadBalancer(d *schema.ResourceData, m interface{}) error {
	opts := buildNetworkLoadBalancerOpts(d, m)
	config := m.(*Config)
	loadBalancerItem := config.LWAPI.NetworkLoadBalancer.Update(opts)
	if loadBalancerItem.HasError() {
		return loadBalancerItem
	}

	updateLoadBalancerResource(d, loadBalancerItem)
	return nil
}

func resourceDeleteNetworkLoadBalancer(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	deleteResponse := config.LWAPI.NetworkLoadBalancer.Delete(d.Id())
	if deleteResponse.HasError() {
		return deleteResponse
	}

	return nil
}

// buildNetworkLoadBalancerOpts builds options for a create/update load balancer API call.
func buildNetworkLoadBalancerOpts(d *schema.ResourceData, m interface{}) network.LoadBalancerParams {
	params := network.LoadBalancerParams{
		Name:               d.Get("name").(string),
		Nodes:              d.Get("nodes").([]network.LoadBalancerNodeParams),
		Region:             d.Get("region").(int),
		Services:           d.Get("services").([]network.LoadBalancerServiceParams),
		SessionPersistence: d.Get("session_persistence").(bool),
		SSLCert:            d.Get("ssl_cert").(string),
		SSLIncludes:        d.Get("ssl_includes").(bool),
		SSLInt:             d.Get("ssl_int").(string),
		SSLKey:             d.Get("ssl_key").(string),
		SSLTermination:     d.Get("ssl_termination").(bool),
		Strategy:           d.Get("strategy").(string),
	}

	return params
}

// loadBalancerDetails gets a load balancer's details from the API.
func loadBalancerDetails(config *Config, id string) *network.LoadBalancerItem {
	return config.LWAPI.NetworkLoadBalancer.Details(id)
}

// updateLoadBalancerResource updates the resource data for the load balancer.
func updateLoadBalancerResource(d *schema.ResourceData, lb *network.LoadBalancerItem) {
	d.Set("name", lb.Name)
	d.Set("nodes", lb.Nodes)
	d.Set("region", lb.RegionID)
	d.Set("services", lb.Services)
	d.Set("session_persistence", lb.SessionPersistence)
	d.Set("ssl_includes", lb.SSLIncludes)
	d.Set("strategy", lb.Strategy)
	d.Set("vip", lb.VIP)
}
