package liquidweb

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceLWStormConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLWStormConfigRead,

		Schema: map[string]*schema.Schema{
			"filter": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"blargs": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// dataSourceLWStormConfigRead gets the available storm configs.
func dataSourceLWStormConfigRead(d *schema.ResourceData, meta interface{}) error {
	return nil
	//config := meta.(*Config)
	//config.Client.Call("v1/Storm/Config/list", opts)
}
