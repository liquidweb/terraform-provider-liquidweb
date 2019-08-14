package liquidweb

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	opentracing "github.com/opentracing/opentracing-go"

	"git.liquidweb.com/masre/liquidweb-go/storm"
)

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
				Computed: true,
				Optional: true,
			},
			"backup_plan": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"backup_quota": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"backup_size": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_quota": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"config_description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"config_id": &schema.Schema{
				Type:     schema.TypeInt,
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
			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vcpu": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"zone": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCreateServer(d *schema.ResourceData, m interface{}) error {
	serverParams := storm.ServerParams{
		ConfigID:     d.Get("config_id").(int),
		Domain:       d.Get("domain").(string),
		ImageID:      d.Get("image_id").(int),
		Template:     d.Get("template").(string),
		Password:     d.Get("password").(string),
		PublicSSHKey: d.Get("public_ssh_key").(string),
		Zone:         d.Get("zone").(int),
	}

	config := m.(*Config)

	tracer := opentracing.GlobalTracer()
	sp := tracer.StartSpan("create-storm-server")
	defer sp.Finish()

	result, err := config.LWAPI.StormServer.Create(serverParams)
	if err != nil {
		traceError(sp, err)
		return err
	}

	d.SetId(result.UniqID)

	stateChange := &resource.StateChangeConf{
		Delay:          10 * time.Second,
		Pending:        storm.ServerStates,
		Refresh:        refreshStormServer(config, d.Id()),
		Target:         []string{"Running"},
		Timeout:        30 * time.Minute,
		NotFoundChecks: 2,
		MinTimeout:     5 * time.Second,
	}
	// https://godoc.org/github.com/hashicorp/terraform/helper/resource#StateRefreshFunc
	// we need to figure out why returning the updated instance isn't updating the server state. Added a call to update at the end of the refresh just for good measure for now.
	statusSpan := opentracing.StartSpan("status-storm-server", opentracing.ChildOf(sp.Context()))
	defer statusSpan.Finish()

	_, err = stateChange.WaitForState()
	if err != nil {
		traceError(statusSpan, err)
		return err
	}

	return resourceReadStormServer(d, m)
}

func resourceReadStormServer(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	tracer := opentracing.GlobalTracer()
	sp := tracer.StartSpan("read-storm-server")
	defer sp.Finish()

	result, err := config.LWAPI.StormServer.Details(d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "LW::Exception::RecordNotFound") {
			d.SetId("")
			return nil
		}
		return err
	}
	updateStormServerResource(d, result)

	return nil
}

func resourceUpdateStormServer(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	params := storm.ServerParams{
		BackupEnabled:  d.Get("backup_enabled").(int),
		BackupPlan:     d.Get("backup_plan").(string),
		BackupQuota:    d.Get("backup_quota").(int),
		BandwidthQuota: d.Get("bandwidth_quota").(string),
		Domain:         d.Get("domain").(string),
		UniqID:         d.Id(),
	}
	tracer := opentracing.GlobalTracer()
	sp := tracer.StartSpan("update-storm-server")
	defer sp.Finish()

	result, err := config.LWAPI.StormServer.Update(params)
	if err != nil {
		return err
	}

	updateStormServerResource(d, result)

	stateChange := &resource.StateChangeConf{
		Delay:          10 * time.Second,
		Pending:        storm.ServerStates,
		Refresh:        refreshStormServer(config, d.Id()),
		Target:         []string{"Running"},
		Timeout:        20 * time.Minute,
		NotFoundChecks: 240,
		MinTimeout:     5 * time.Second,
	}
	statusSpan := opentracing.StartSpan("status-storm-server", opentracing.ChildOf(sp.Context()))
	defer statusSpan.Finish()

	server, err := stateChange.WaitForState()
	if err != nil {
		return err
	}

	updateStormServerResource(d, server.(*storm.Server))

	return nil
}

func resourceDeleteStormServer(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	tracer := opentracing.GlobalTracer()
	sp := tracer.StartSpan("destroy-storm-server")
	defer sp.Finish()

	_, err := config.LWAPI.StormServer.Destroy(d.Id())
	if err != nil {
		return err
	}

	stateChange := &resource.StateChangeConf{
		Delay:          10 * time.Second,
		Pending:        []string{"pending destruction"},
		Refresh:        refreshDestroyStormServerStatus(config, d.Id()),
		Target:         []string{"destroyed"},
		Timeout:        20 * time.Minute,
		NotFoundChecks: 240,
		MinTimeout:     5 * time.Second,
	}
	statusSpan := opentracing.StartSpan("status-storm-server", opentracing.ChildOf(sp.Context()))
	defer statusSpan.Finish()

	_, err = stateChange.WaitForState()
	if err != nil {
		traceError(statusSpan, err)
		return fmt.Errorf(
			"Error waiting for instance (%s) to be destroyed: %s", d.Id(), err)
	}
	d.SetId("")

	return nil
}

// refreshDestroyStormServerStatus queries the API for the status of the server destroy.
func refreshDestroyStormServerStatus(config *Config, uid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := config.LWAPI.StormServer.Status(uid)
		if err != nil {
			if strings.Contains(err.Error(), "LW::Exception::RecordNotFound") {
				return s, "destroyed", nil
			}
			return nil, "", err
		}

		return nil, "pending destruction", nil
	}
}

// refreshStormServer queries the API for status.
// If the status is "Running" query for its details and return them.
func refreshStormServer(config *Config, uid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		result, err := config.LWAPI.StormServer.Status(uid)
		if err != nil {
			return nil, "", err
		}

		if len(result.Status) == 0 {
			return nil, "NotReady", nil
		}

		// Get server details if it's running.
		if result.Status == "Running" {
			details, err := config.LWAPI.StormServer.Details(uid)
			if err != nil {
				return nil, "", err
			}

			// Ensure we have an IP, otherwise return a pseudo-status until it does have an IP.
			if len(details.IP) == 0 {
				return nil, "NotReady", nil
			}

			return details, result.Status, nil
		}

		// Return an empty string as if nil is returned the resource will be considered "not found".
		// See
		return "", result.Status, nil
	}
}

// updateStormServerResource updates the resource data for the storm server.
func updateStormServerResource(d *schema.ResourceData, server *storm.Server) {
	d.SetId(server.UniqID)
	d.Set("accnt", server.ACCNT)
	d.Set("active", server.Active)
	d.Set("backup_enabled", server.BackupEnabled)
	d.Set("backup_plan", server.BackupPlan)
	d.Set("backup_quota", server.BackupQuota)
	d.Set("backup_size", server.BackupSize)
	d.Set("bandwidth_quota", server.BandwidthQuota)
	d.Set("config_description", server.ConfigDescription)
	d.Set("config_id", server.ConfigID)
	d.Set("domain", server.Domain)
	d.Set("ip", server.IP.String())
	d.Set("ip_count", server.IPCount)
	d.Set("manage_level", server.ManageLevel)
	d.Set("memory", server.Memory)
	d.Set("template", server.Template)
	d.Set("template_description", server.TemplateDescription)
	d.Set("type", server.Type)
	d.Set("vcpu", server.VCPU)
	d.Set("zone", server.Zone)
}
