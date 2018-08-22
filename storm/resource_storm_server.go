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
			"public_ssh_key": &schema.Schema{
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
