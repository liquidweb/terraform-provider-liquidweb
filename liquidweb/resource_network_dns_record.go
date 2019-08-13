package liquidweb

import (
	"strconv"

	network "git.liquidweb.com/masre/liquidweb-go/network"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNetworkDNSRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateNetworkDNSRecord,
		Read:   resourceReadNetworkDNSRecord,
		Update: resourceUpdateNetworkDNSRecord,
		Delete: resourceDeleteNetworkDNSRecord,

		Schema: map[string]*schema.Schema{
			"admin_email": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"exchange": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiry": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"full_data": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"minimum": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nameserver": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"prio": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"rdata": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"refresh_interval": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_overrides": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"retry": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"serial": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"target": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"weight": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"zone_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCreateNetworkDNSRecord(d *schema.ResourceData, m interface{}) error {
	opts := buildNetworkDNSRecordOpts(d, m)
	config := m.(*Config)

	result, err := config.LWAPI.NetworkDNS.Create(opts)
	if err != nil {
		return err
	}

	id := strconv.Itoa(int(result.ID))
	d.SetId(id)

	return resourceReadNetworkDNSRecord(d, m)
}

func resourceReadNetworkDNSRecord(d *schema.ResourceData, m interface{}) error {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	config := m.(*Config)
	dnsRecord, err := config.LWAPI.NetworkDNS.Details(id)
	if err != nil {
		return err
	}

	updateDNSRecordResource(d, dnsRecord)
	return nil
}

func resourceUpdateNetworkDNSRecord(d *schema.ResourceData, m interface{}) error {
	opts := buildNetworkDNSRecordOpts(d, m)
	// Attach ID to params.
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	opts.ID = id

	// API call for update does not accept zone info or type.
	opts.ZoneID = 0
	opts.Zone = ""
	opts.Type = ""

	config := m.(*Config)
	dnsRecord, err := config.LWAPI.NetworkDNS.Update(opts)
	if err != nil {
		return err
	}

	updateDNSRecordResource(d, dnsRecord)
	return nil
}

func resourceDeleteNetworkDNSRecord(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	params := &network.DNSRecordParams{ID: id}

	_, err = config.LWAPI.NetworkDNS.Delete(params)
	if err != nil {
		return err
	}

	return nil
}

// buildNetworkDNSRecordOpts builds options for a create/update DNS record API call.
func buildNetworkDNSRecordOpts(d *schema.ResourceData, m interface{}) *network.DNSRecordParams {
	params := &network.DNSRecordParams{
		Zone:  d.Get("zone").(string),
		Name:  d.Get("name").(string),
		Prio:  d.Get("prio").(int),
		RData: d.Get("rdata").(string),
		Type:  d.Get("type").(string),
		TTL:   d.Get("ttl").(int),
	}

	return params
}

// updateDNSRecordResource updates the resource data for the DNS Record.
func updateDNSRecordResource(d *schema.ResourceData, dr *network.DNSRecord) {
	d.Set("admin_email", dr.AdminEmail)
	d.Set("created", dr.Created)
	d.Set("exchange", dr.Exchange)
	d.Set("expiry", dr.Expiry)
	d.Set("full_data", dr.FullData)
	d.Set("last_updated", dr.LastUpdated)
	d.Set("minimum", dr.Minimum)
	d.Set("name", dr.Name)
	d.Set("nameserver", dr.Nameserver)
	d.Set("port", dr.Port)
	d.Set("prio", dr.Prio)
	d.Set("rdata", dr.RData)
	d.Set("refresh_interval", dr.RefreshInterval)
	d.Set("region_overrides", dr.RegionOverrides)
	d.Set("retry", dr.Retry)
	d.Set("serial", dr.Serial)
	d.Set("target", dr.Target)
	d.Set("ttl", dr.TTL)
	d.Set("type", dr.Type)
	d.Set("weight", dr.Weight)
	d.Set("zone_id", dr.ZoneID)
}
