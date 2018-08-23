package storm

import (
	"github.com/hashicorp/terraform/helper/schema"
	lwapi "github.com/liquidweb/go-lwApi"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"config_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"image_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
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
			},
			"zone": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	opts := buildServerOpts(d, m)
	config := m.(*Config)
	rawr, err := config.Client.Call("v1/Storm/Server/create", opts)
	if err != nil {
		return err
	}
	r := rawr.(map[string]interface{})
	uid := r["uniq_id"].(string)
	d.SetId(uid)
	return nil
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	opts := buildServerOpts(d, m)
	uid := d.Id()
	opts["uniq_id"] = uid
	validOpts := pickDetailsOpts(opts)
	config := m.(*Config)
	_, err := config.Client.Call("v1/Storm/Server/details", validOpts)
	if err != nil {
		errClass, ok := err.(lwapi.LWAPIError)
		if ok && errClass.ErrorClass == "LW::Exception::RecordNotFound" {
			d.SetId("")
			return nil
		}
		return err
	}

	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	opts := buildServerOpts(d, m)
	uid := d.Id()
	opts["uniq_id"] = uid
	validOpts := pickUpdateOpts(opts)
	config := m.(*Config)
	_, err := config.Client.Call("v1/Storm/Server/update", validOpts)
	if err != nil {
		return err
	}

	return nil
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	uid := d.Id()
	opts := make(map[string]interface{})
	opts["uniq_id"] = uid
	_, err := config.Client.Call("v1/Storm/Server/destroy", opts)
	if err != nil {
		return err
	}
	return nil
}

// ServerOpts are options passed to Storm API calls
type ServerOpts struct {
	ConfigID     int
	Domain       string
	ImageID      int
	Password     string
	PublicSSHKey string
	Template     string
	Zone         int
}

func buildServerOpts(d *schema.ResourceData, m interface{}) map[string]interface{} {
	so := &ServerOpts{
		ConfigID:     d.Get("config_id").(int),
		Domain:       d.Get("domain").(string),
		ImageID:      d.Get("image_id").(int),
		Template:     d.Get("template").(string),
		Password:     d.Get("password").(string),
		PublicSSHKey: d.Get("public_ssh_key").(string),
		Zone:         d.Get("zone").(int),
	}
	// The Storm API client uses a string map for parameters.
	var smOpts = make(map[string]interface{})
	smOpts["config_id"] = so.ConfigID
	smOpts["domain"] = so.Domain

	// Add Template if provided.
	if len(so.Template) > 0 {
		smOpts["template"] = so.Template
	}

	// Add Image if provided.
	if so.ImageID != 0 {
		smOpts["image_id"] = so.ImageID
	}

	smOpts["password"] = so.Password
	smOpts["public_ssh_key"] = so.PublicSSHKey
	smOpts["template"] = so.Template
	smOpts["zone"] = so.Zone

	return smOpts
}

// pickUpdateOpts returns a set of options valid for an update request.
func pickUpdateOpts(opts map[string]interface{}) map[string]interface{} {
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
func pickDetailsOpts(opts map[string]interface{}) map[string]interface{} {
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
