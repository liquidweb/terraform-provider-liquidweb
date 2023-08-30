package liquidweb

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	network "github.com/liquidweb/liquidweb-go/network"
	opentracing "github.com/opentracing/opentracing-go"
)

func resourceFloatingIP() *schema.Resource {
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

	result, err := config.LWAPI.NetworkVIP.Create(opts)
	if err != nil {
		traceError(sp, err)
		return err
	}

	sp.Finish()

	d.SetId(result.UniqID)
	d.Set("zone", d.Get("zone"))

	return resourceReadNetworkVIP(d, m)
}

func resourceReadNetworkVIP(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	vip, err := config.LWAPI.NetworkVIP.Details(d.Id())

	if err != nil {
		// If VIP was destroyed outside of Terraform, set ID to an empty string so Terraform treats it as destroyed.
		if strings.Contains(err.Error(), "LW::Exception::RecordNotFound") {
			d.SetId("")
			return nil
		}

		return err
	}

	updateVIPResource(d, vip)

	return nil
}

func resourceDeleteNetworkVIP(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	tracer := opentracing.GlobalTracer()
	sp := tracer.StartSpan("destroy-network-vip")

	_, err := config.LWAPI.NetworkVIP.Destroy(d.Id())
	if err != nil {
		traceError(sp, err)
		return err
	}
	sp.Finish()

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

// updateVIPResource updates the resource data for the VIP.
func updateVIPResource(d *schema.ResourceData, dr *network.VIP) {
	d.Set("domain", dr.Domain)
	d.Set("active", dr.Active)
	d.Set("activeStatus", dr.ActiveStatus)
	d.Set("uniq_id", dr.UniqID)
	d.Set("destroyed", dr.Destroyed)
	d.Set("ip", dr.IP)
}
