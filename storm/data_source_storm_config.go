package storm

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceStormConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceStormConfigRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"blargs": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}
