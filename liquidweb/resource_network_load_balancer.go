package liquidweb

import (
	network "git.liquidweb.com/masre/liquidweb-go/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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
			"nodes": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
				Optional: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"services": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"src_port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"dest_port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
					},
				},
			},
			"session_persistence": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
				Default:  false,
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

	d.SetId(result.UniqID)

	return resourceReadNetworkLoadBalancer(d, m)
}

func resourceReadNetworkLoadBalancer(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	loadBalancer, err := config.LWAPI.NetworkLoadBalancer.Details(d.Id())
	if err != nil {
		return err
	}

	updateLoadBalancerResource(d, loadBalancer)
	return nil
}

func resourceUpdateNetworkLoadBalancer(d *schema.ResourceData, m interface{}) error {
	opts := buildNetworkLoadBalancerUpdateOpts(d, m)
	config := m.(*Config)

	loadBalancer, err := config.LWAPI.NetworkLoadBalancer.Update(opts)
	if err != nil {
		return err
	}

	updateLoadBalancerResource(d, loadBalancer)
	return nil
}

func resourceDeleteNetworkLoadBalancer(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	_, err := config.LWAPI.NetworkLoadBalancer.Delete(d.Id())
	if err != nil {
		return err
	}

	return nil
}

// buildNetworkLoadBalancerOpts builds options for a create load balancer API call.
func buildNetworkLoadBalancerOpts(d *schema.ResourceData, m interface{}) network.LoadBalancerParams {
	params := network.LoadBalancerParams{
		Name:               d.Get("name").(string),
		Nodes:              expandSetToStrings(d.Get("nodes").(*schema.Set).List()),
		Region:             d.Get("region").(int),
		Services:           expandServicesSet(d.Get("services").(*schema.Set).List()),
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

// buildNetworkLoadBalancerUpdateOpts builds options for a update load balancer API call.
func buildNetworkLoadBalancerUpdateOpts(d *schema.ResourceData, m interface{}) network.LoadBalancerParams {
	params := network.LoadBalancerParams{
		UniqID:             d.Id(),
		Name:               d.Get("name").(string),
		Nodes:              expandSetToStrings(d.Get("nodes").(*schema.Set).List()),
		Services:           expandServicesSet(d.Get("services").(*schema.Set).List()),
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

// updateLoadBalancerResource updates the resource data for the load balancer.
func updateLoadBalancerResource(d *schema.ResourceData, lb *network.LoadBalancer) {
	d.Set("name", lb.Name)
	d.Set("nodes", lb.Nodes)
	d.Set("region", lb.RegionID)
	d.Set("services", lb.Services)
	d.Set("session_persistence", lb.SessionPersistence)
	d.Set("ssl_includes", lb.SSLIncludes)
	d.Set("strategy", lb.Strategy)
	d.Set("vip", lb.VIP.String())
}

// expandServicesSet expands types.TypeSet into an actual services set.
func expandServicesSet(services []interface{}) []network.LoadBalancerServiceParams {
	expandedServices := make([]network.LoadBalancerServiceParams, len(services))
	for i, v := range services {
		service := v.(map[string]interface{})
		expandedServices[i] = network.LoadBalancerServiceParams{
			SrcPort:  service["src_port"].(int),
			DestPort: service["dest_port"].(int),
		}
	}

	return expandedServices
}
