package liquidweb

import (
	"fmt"

	lwapi "git.liquidweb.com/masre/liquidweb-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceLWNetworkZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLWNetworkZoneRead,

		Schema: map[string]*schema.Schema{
			"region_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"valid_source_hvs": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// dataSourceLWNetworkZoneRead gets the available storm configs.
func dataSourceLWNetworkZoneRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	params := &lwapi.ZoneListParams{}

	regionName, ok := d.GetOk("region_name")
	if ok {
		name, rOk := regionName.(string)
		if rOk {
			params.Region = name
		}
	}

	result, err := config.LWAPI.NetworkZone.List(params)
	if err != nil {
		return err
	}

	if result.HasError() {
		return result
	}

	if result.ItemCount != 1 {
		return fmt.Errorf("Search returned %d results, please revise so only one is returned", result.ItemCount)
	}

	item := result.Items[0]
	d.SetId(string(item.ID))
	d.Set("is_default", item.IsDefault)
	d.Set("name", item.Name)
	d.Set("region", item.Region)
	d.Set("status", item.Status)
	d.Set("valid_source_hvs", item.ValidSourceHVS)

	return nil
}
