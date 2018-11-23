package liquidweb

import (
	"fmt"

	lwnetwork "git.liquidweb.com/masre/liquidweb-go/network"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceLWNetworkZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLWNetworkZoneRead,

		Schema: map[string]*schema.Schema{
			"network_zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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

// dataSourceLWNetworkZoneRead gets the available network zones.
func dataSourceLWNetworkZoneRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	params := &lwnetwork.ZoneListParams{}

	_, ok := d.GetOk("region_name")
	if ok {
		params.Region = d.Get("region_name").(string)
	}

	result, err := config.LWAPI.NetworkZone.List(params)

	if err != nil {
		return err
	}

	filteredNetworkZones := filterLWNetworkZones(result, d)

	if len(filteredNetworkZones) != 1 {
		return fmt.Errorf("Search returned %d results, please revise so only one is returned", len(filteredNetworkZones))
	}

	item := filteredNetworkZones[0]
	d.SetId(item.ID.String())
	d.Set("network_zone_id", item.ID.String())
	d.Set("is_default", item.IsDefault)
	d.Set("name", item.Name)
	d.Set("region", item.Region)
	d.Set("status", item.Status)
	d.Set("valid_source_hvs", item.ValidSourceHVS)
	d.Set("region_name", item.Region.Name)

	return nil
}

func filterLWNetworkZones(zoneList *lwnetwork.ZoneList, d *schema.ResourceData) []lwnetwork.Zone {
	_, nameOk := d.GetOk("name")
	var name string
	if nameOk {
		name = d.Get("name").(string)
	}

	filteredNetworkZones := []lwnetwork.Zone{}

	for _, z := range zoneList.Items {
		if nameOk && name != z.Name {
			continue
		}

		filteredNetworkZones = append(filteredNetworkZones, z)
	}

	return filteredNetworkZones
}
