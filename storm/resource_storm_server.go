package storm

import (
	"github.com/hashicorp/terraform/helper/schema"
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
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
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
