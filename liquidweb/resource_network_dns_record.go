package liquidweb

import (
	"fmt"
	"strconv"
	"time"

	network "git.liquidweb.com/masre/liquidweb-go/network"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	lwapi "github.com/liquidweb/go-lwApi"
)

var dnsRecordFields = []string{
	"adminEmail",
	"created",
	"exchange",
	"expiry",
	"fullData",
	"id",
	"last_updated",
	"minimum",
	"name",
	"nameserver",
	"port",
	"prio",
	"rdata",
	"refreshInterval",
	"regionOverrides",
	"retry",
	"serial",
	"target",
	"ttl",
	"type",
	"weight",
	"zone_id",
}

func resourceNetworkDNSRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateNetworkDNSRecord,
		Read:   resourceReadNetworkDNSRecord,
		Update: resourceUpdateNetworkDNSRecord,
		Delete: resourceDeleteNetworkDNSRecord,

		Schema: map[string]*schema.Schema{
			"adminEmail": &schema.Schema{
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
			"fullData": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeInt,
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
			"refreshInterval": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"regionOverrides": &schema.Schema{
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
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"weight": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"zone_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceCreateNetworkDNSRecord(d *schema.ResourceData, m interface{}) error {
	opts := buildNetworkDNSRecordOpts(d, m)
	config := m.(*Config)

	result := config.LWAPI.NetworkDNS.Create(opts)
	if result.HasError() {
		return result
	}

	id := strconv.Itoa(result.ID)
	d.SetId(id)

	return resourceUpdateStormServer(d, m)
}

func resourceReadNetworkDNSRecord(d *schema.ResourceData, m interface{}) error {
	uid := d.Id()
	config := m.(*Config)
	server, err := stormServerDetails(config, uid)
	if err != nil {
		errClass, ok := err.(lwapi.LWAPIError)
		if ok && errClass.ErrorClass == "LW::Exception::RecordNotFound" {
			d.SetId("")
			return nil
		}
		return err
	}

	updateStormServerResource(d, server)

	return nil
}

func resourceUpdateNetworkDNSRecord(d *schema.ResourceData, m interface{}) error {
	opts := buildUpdateStormServerOpts(d, m)
	validOpts := pickStormServerUpdateOpts(opts)
	config := m.(*Config)
	_, err := config.Client.Call("v1/Storm/Server/update", validOpts)
	if err != nil {
		return err
	}

	stateChange := &resource.StateChangeConf{
		Delay:          10 * time.Second,
		Pending:        stormServerStates,
		Refresh:        refreshStormServer(config, d.Id()),
		Target:         []string{"Running"},
		Timeout:        20 * time.Minute,
		NotFoundChecks: 240,
		MinTimeout:     5 * time.Second,
	}
	_, err = stateChange.WaitForState()
	if err != nil {
		return err
	}

	return resourceReadStormServer(d, m)
}

func resourceDeleteNetworkDNSRecord(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	uid := d.Id()
	opts := make(map[string]interface{})
	opts["uniq_id"] = uid
	_, err := config.Client.Call("v1/Storm/Server/destroy", opts)
	if err != nil {
		return err
	}

	stateChange := &resource.StateChangeConf{
		Delay:          10 * time.Second,
		Pending:        stormServerStates,
		Refresh:        refreshStormServer(config, d.Id()),
		Target:         []string{"Destroying"},
		Timeout:        20 * time.Minute,
		NotFoundChecks: 240,
		MinTimeout:     5 * time.Second,
	}
	_, err = stateChange.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to be destroyed: %s", uid, err)
	}

	return nil
}

// buildNetworkDNSRecordOpts builds options for a create/update DNS record API call.
func buildNetworkDNSRecordOpts(d *schema.ResourceData, m interface{}) *network.DNSRecordParams {
	params := &network.DNSRecordParams{
		Name:  d.Get("name").(string),
		Prio:  d.Get("prio").(int),
		RData: d.Get("rdata").(string),
		Type:  d.Get("type").(string),
	}

	// Add Zone ID if provided.
	zid := d.Get("zone_id").(int)
	if zid > 0 {
		params.ZoneID = zid
	}

	// Add Zone if Zone ID isn't set.
	zone := d.Get("zone").(string)
	if zid == 0 {
		params.Zone = zone
	}

	ttl := d.Get("ttl").(int)
	if ttl > 0 {
		params.TTL = ttl
	}

	return params
}
