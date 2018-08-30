package liquidweb

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	lwapi "github.com/liquidweb/go-lwApi"
)

var stormServerFields = []string{
	"accnt",
	"active",
	"backup_enabled",
	"backup_plan",
	"backup_quota",
	"backup_size",
	"bandwidth_quota",
	"config_description",
	"config_id",
	"create_date",
	"diskspace",
	"domain",
	"ip",
	"ip_count",
	"manage_level",
	"memory",
	"template",
	"template_description",
	"type",
	"uniq_id",
	"vcpu",
	"zone",
}

var stormServerStates = []string{
	"Building",
	"Cloning",
	"Resizing",
	"Moving",
	"Booting",
	"Stopping",
	"Restarting",
	"Rebooting",
	"Shutting Down",
	"Restoring Backup",
	"Creating Image",
	"Deleting Image",
	"Restoring Image",
	"Re-Imaging",
	"Updating Firewall",
	"Updating Network",
	"Adding IPs",
	"Removing IP",
	"Destroying",
}

func resourceStormServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateServer,
		Read:   resourceReadStormServer,
		Update: resourceUpdateStormServer,
		Delete: resourceDeleteStormServer,

		Schema: map[string]*schema.Schema{
			"accnt": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"active": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"backup_enabled": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"backup_plan": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"backup_quota": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"backup_size": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_quota": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"config_description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"config_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"image_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"ip_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"manage_level": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"memory": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ssh_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"template_description": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"uniq_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vcpu": &schema.Schema{
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

func resourceCreateServer(d *schema.ResourceData, m interface{}) error {
	opts := buildCreateStormServerOpts(d, m)
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
	_, err = stateChange.WaitForState()
	if err != nil {
		return err
	}

	return nil
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

	return nil
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
		Target:         []string{"Running"},
		Timeout:        20 * time.Minute,
		NotFoundChecks: 240,
		MinTimeout:     5 * time.Second,
	}
	_, err = stateChange.WaitForState()
	if err != nil {
		return err
	}

	return nil
}

// StormServerOpts are options passed to Storm API calls
type StormServerOpts struct {
	BackupEnabled  int
	BackupPlan     string
	BackupQuota    int
	BandwidthQuota int
	ConfigID       int
	Domain         string
	ImageID        int
	Password       string
	PublicSSHKey   string
	Template       string
	UniqID         string
	Zone           int
}

// buildCreateStormServerOpts builds options for a create server API call.
func buildCreateStormServerOpts(d *schema.ResourceData, m interface{}) map[string]interface{} {
	so := &StormServerOpts{
		ConfigID:     d.Get("config_id").(int),
		Domain:       d.Get("domain").(string),
		ImageID:      d.Get("image_id").(int),
		Template:     d.Get("template").(string),
		Password:     d.Get("password").(string),
		PublicSSHKey: d.Get("public_ssh_key").(string),
		Zone:         d.Get("zone").(int),
	}
	// The Storm API client uses a string map for parameters.
	var opts = make(map[string]interface{})
	opts["config_id"] = so.ConfigID
	opts["domain"] = so.Domain

	// Add Template if provided.
	if len(so.Template) > 0 {
		opts["template"] = so.Template
	}

	// Add Image if provided.
	if so.ImageID != 0 {
		opts["image_id"] = so.ImageID
	}

	opts["password"] = so.Password
	opts["public_ssh_key"] = so.PublicSSHKey
	opts["template"] = so.Template
	opts["zone"] = so.Zone

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

		// Get server details if it's running.
		if state == "Running" {
			rawr, err := stormServerDetails(config, uid)
			if err != nil {
				return nil, "", err
			}

			return rawr, state, nil
		}

		return nil, state, nil
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
