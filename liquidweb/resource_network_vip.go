package liquidweb

import (
	"strings"
	network "git.liquidweb.com/masre/liquidweb-go/network"
	"github.com/hashicorp/terraform/helper/schema"
	opentracing "github.com/opentracing/opentracing-go"
	opentracinglog "github.com/opentracing/opentracing-go/log"
)

func resourceNetworkVIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateNetworkVIP,
		Read:   resourceReadNetworkVIP,
		Delete: resourceDeleteNetworkVIP,

		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"active": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"active_status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"uniq_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"subaccnt": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"destroyed": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCreateNetworkVIP(d *schema.ResourceData, m interface{}) error {
	opts := buildNetworkVIPOpts(d, m)
	config := m.(*Config)
	
	tracer := opentracing.GlobalTracer()
	sp := tracer.StartSpan("create-network-vip")
	defer sp.Finish()
	
	result, err := config.LWAPI.NetworkVIP.Create(opts)
	if err != nil {
		sp.SetTag("error", "true")
		sp.LogFields(opentracinglog.String("error", err.Error()))
		return err
	}

	if result.HasError() {
			sp.SetTag("error", "true")
			sp.LogFields(
				opentracinglog.String("error", result.Error()),
			)
		return result
	}

	d.SetId(result.UniqID)
	d.Set("zone", d.Get("zone"))

	return resourceReadNetworkVIP(d, m)
}

func resourceReadNetworkVIP(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	vipItem := VIPDetails(config, d.Id())

	if vipItem.HasError() {
		// If VIP was destroyed outside of Terraform, set id to nil value and return nil.
		if strings.Contains(vipItem.Error(), "LW::Exception::RecordNotFound") {
			d.SetId("")
			return nil
		}

		return vipItem
	}

	updateVIPResource(d, vipItem)

	return nil
}

func resourceDeleteNetworkVIP(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	tracer := opentracing.GlobalTracer()
	sp := tracer.StartSpan("destroy-network-vip")
	defer sp.Finish()

	deleteResponse := config.LWAPI.NetworkVIP.Destroy(d.Id())
	if deleteResponse.HasError() {
		sp.SetTag("error", "true")
		sp.LogFields(
			opentracinglog.String("error", deleteResponse.Error()),
		)
		return deleteResponse
	}

	return nil
}

// buildNetworkVIPOpts builds options for a create VIP API call.
func buildNetworkVIPOpts(d *schema.ResourceData, m interface{}) network.VIPParams {
	params := network.VIPParams{
		Domain: d.Get("domain").(string),
		Zone:   d.Get("zone").(int),
	}

	return params
}

// VIPDetails gets a VIP's details from the API.
func VIPDetails(config *Config, id string) *network.VIPItem {
	return config.LWAPI.NetworkVIP.Details(id)
}

// updateVIPResource updates the resource data for the VIP.
func updateVIPResource(d *schema.ResourceData, dr *network.VIPItem) {
	d.Set("domain", dr.Domain)
	d.Set("active", dr.Active)
	d.Set("activeStatus", dr.ActiveStatus)
	d.Set("uniq_id", dr.UniqID)
	d.Set("destroyed", dr.Destroyed)
	d.Set("ip", dr.IP)
}
