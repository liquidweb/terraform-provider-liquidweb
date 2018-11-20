package liquidweb

import (
	"fmt"
	"log"
	"time"

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
	rawr, err := config.Client.Call("v1/Storm/Server/create", opts)
	if err != nil {
		return err
	}
	resp := rawr.(map[string]interface{})
	uid := resp["uniq_id"].(string)
	d.SetId(uid)

	stateChange := &resource.StateChangeConf{
		Delay:          10 * time.Second,
		Pending:        stormServerStates,
		Refresh:        refreshStormServer(config, uid),
		Target:         []string{"Running"},
		Timeout:        20 * time.Minute,
		NotFoundChecks: 240,
		MinTimeout:     5 * time.Second,
	}
	// https://godoc.org/github.com/hashicorp/terraform/helper/resource#StateRefreshFunc
	// we need to figure out why returning the updated instance isn't updating the server state. Added a call to update at the end of the refresh just for good measure for now.
	_, err = stateChange.WaitForState()
	if err != nil {
		return err
	}

	return resourceUpdateStormServer(d, m)
}

func resourceReadStormServer(d *schema.ResourceData, m interface{}) error {
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

func resourceUpdateStormServer(d *schema.ResourceData, m interface{}) error {
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

func resourceDeleteStormServer(d *schema.ResourceData, m interface{}) error {
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

// NetworkDNSRecordOptsOpts are options passed to API calls
type NetworkDNSRecordOptsOpts struct {
	Name   string
	Prio   int
	RData  string
	TTL    int
	Type   string
	Zone   string
	ZoneID int
}

// buildNetworkDNSRecordOpts builds options for a create/update DNS record API call.
func buildNetworkDNSRecordOpts(d *schema.ResourceData, m interface{}) map[string]interface{} {
	so := &NetworkDNSRecordOptsOpts{
		Name:   d.Get("name").(int),
		Prio:   d.Get("prio").(string),
		RData:  d.Get("rdata").(int),
		TTL:    d.Get("ttl").(string),
		Type:   d.Get("type").(string),
		Zone:   d.Get("zone").(string),
		ZoneID: d.Get("zone_id").(int),
	}

	// The API client uses a string map for parameters.
	var opts = make(map[string]interface{})
	opts["name"] = so.Name
	opts["rdata"] = so.RData
	opts["type"] = so.Type

	// Add Zone ID if provided.
	if len(so.ZoneID) > 0 {
		opts["zone_id"] = so.ZoneID
	}

	// Add Zone if Zone ID isn't set.
	if _, ok := opts["zone_id"] > 0; !ok {
		opts["zone"] = so.Zone
	}

	if so.Prio >= 0 {
		opts["prio"] = so.Prio
	}

	if so.TTL >= 0 {
		opts[ttl] = so.TTL
	}

	return opts
}

// buildUpdateStormServerOpts builds options for an update server API call.
func buildUpdateStormServerOpts(d *schema.ResourceData, m interface{}) map[string]interface{} {
	so := &StormServerOpts{
		BackupEnabled:  d.Get("backup_enabled").(int),
		BackupPlan:     d.Get("backup_plan").(string),
		BackupQuota:    d.Get("backup_quota").(int),
		BandwidthQuota: d.Get("bandwidth_quota").(int),
		Domain:         d.Get("domain").(string),
		UniqID:         d.Id(),
	}
	// The Storm API client uses a string map for parameters.
	var opts = make(map[string]interface{})
	opts["backup_enabled"] = so.BackupEnabled
	if len(so.BackupPlan) > 0 {
		opts["backup_plan"] = so.BackupPlan
	}
	if so.BackupQuota > 0 {
		opts["backup_quota"] = so.BackupQuota
	}
	opts["bandwidth_quota"] = so.BandwidthQuota
	opts["domain"] = so.Domain
	opts["uniq_id"] = so.UniqID

	return opts
}

// pickUpdateOpts returns a set of options valid for an update request.
func pickStormServerUpdateOpts(opts map[string]interface{}) map[string]interface{} {
	allowed := [6]string{"backup_enabled", "backup_plan", "backup_quota", "bandwidth_quota", "domain", "uniq_id"}
	validOpts := make(map[string]interface{})

	for _, af := range allowed {
		f, ok := opts[af]
		if ok {
			validOpts[af] = f
		}
	}

	return validOpts
}

// pickDetailsOpts returns a set of options valid for a details request.
func pickStormServerDetailsOpts(opts map[string]interface{}) map[string]interface{} {
	allowed := [6]string{"uniq_id"}
	validOpts := make(map[string]interface{})

	for _, af := range allowed {
		f, ok := opts[af]
		if ok {
			validOpts[af] = f
		}
	}

	return validOpts
}

// refreshStormServer queries the API for status returns the current status.
// If the status is "Running" query for its details and return them.
func refreshStormServer(config *Config, uid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opts := make(map[string]interface{})
		opts["uniq_id"] = uid
		rawr, err := config.Client.Call("v1/Storm/Server/status", opts)
		if err != nil {
			return nil, "", err
		}
		resp := rawr.(map[string]interface{})
		status, ok := resp["status"]
		if !ok {
			return nil, "", fmt.Errorf("problem getting server status")
		}
		state := status.(string)
		log.Printf("status returned: %v", state)

		// Get server details if it's running.
		if state == "Running" {
			rawr, err := stormServerDetails(config, uid)
			if err != nil {
				return nil, "", err
			}

			return rawr, state, nil
		}

		// Return an empty string as if nil is returned the resource will be considered "not found".
		// See
		return "", state, nil
	}
}

// serverDetails gets server details from the API.
func stormServerDetails(config *Config, uid string) (interface{}, error) {
	opts := make(map[string]interface{})
	opts["uniq_id"] = uid
	return config.Client.Call("v1/Storm/Server/details", opts)
}

// updateStormServerResource updates the resource data for the storm server.
func updateStormServerResource(d *schema.ResourceData, server interface{}) {
	ss := server.(map[string]interface{})

	fields := stormServerFields

	for _, field := range fields {
		f, ok := ss[field]
		if ok {
			d.Set(field, f)
		}
		if field == "uniq_id" {
			d.SetId(f.(string))
		}
	}
}
