package liquidweb

import (
	"fmt"
	"log"

	"github.com/liquidweb/liquidweb-go/storm"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceServerConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerConfigRead,

		Schema: map[string]*schema.Schema{
			// Filters
			"network_zone": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"config_id": &schema.Schema{
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
			},
			// Attributes
			"active": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"available": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"category": {
				Type:     schema.TypeString,
				Default:  "storm",
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk": { // disk space (in GB) included with VM
				Type:     schema.TypeInt,
				Optional: true,
			},
			"featured": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"memory": { // memory (in mb) included with vm
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vcpu": { // CPU threads / vcpus included with vm
				Type:     schema.TypeInt,
				Optional: true,
			},
			"zone_availability": { // which zones this config is availble in
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

// dataSourceServerConfigRead gets the available server configs.
func dataSourceServerConfigRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	params := storm.ConfigListParams{}
	params.PageSize = 200 // Grab all configs at once.

	available, ok := d.GetOk("available")
	if ok {
		name, rOk := available.(bool)
		if rOk {
			params.Available = name
		}
	}

	result, err := config.LWAPI.StormConfig.List(params)
	if err != nil {
		return err
	}

	// Filter list based on various criteria.
	filteredConfigs := filterServerConfig(result, d)
	log.Printf("blars: %+v", filteredConfigs)
	if len(filteredConfigs) != 1 {
		return fmt.Errorf("Search returned %d results, please revise so only one is returned", len(filteredConfigs))
	}

	item := filteredConfigs[0]

	d.SetId(item.ID.String())
	d.Set("active", item.Active)
	d.Set("available", item.Available)
	d.Set("category", item.Category)
	d.Set("description", item.Description)
	d.Set("disk", item.Disk)
	d.Set("featured", item.Featured)
	d.Set("memory", item.Memory)
	d.Set("vcpu", item.VCPU)
	d.Set("zone_availability", item.ZoneAvailability)

	return nil
}

func filterServerConfig(configList *storm.ConfigList, d *schema.ResourceData) []storm.Config {
	active := d.Get("active").(bool)
	available := d.Get("available").(bool)
	category := d.Get("category").(string)
	_, descriptionOk := d.GetOk("description")
	var description string
	if descriptionOk {
		description = d.Get("description").(string)
	}

	_, diskOk := d.GetOk("disk")
	var disk int
	if diskOk {
		disk = d.Get("disk").(int)
	}

	_, featuredOk := d.GetOk("featured")
	var featured bool
	if featuredOk {
		featured = d.Get("featured").(bool)
	}

	_, memoryOk := d.GetOk("memory")
	var memory int
	if memoryOk {
		memory = d.Get("memory").(int)
	}

	_, vcpuOk := d.GetOk("vcpu")
	var vcpu int
	if vcpuOk {
		vcpu = d.Get("vcpu").(int)
	}

	_, networkZoneOk := d.GetOk("network_zone")
	var networkZone string
	if networkZoneOk {
		networkZone = d.Get("network_zone").(string)
	}

	_, configIDOk := d.GetOk("config_id")
	var configID int
	if configIDOk {
		configID = d.Get("config_id").(int)
	}

	filteredConfigs := []storm.Config{}

	for _, c := range configList.Items {
		if active != bool(c.Active) {
			continue
		}

		if available != bool(c.Available) {
			continue
		}

		if category != c.Category {
			continue
		}

		if descriptionOk && c.Description != description {
			continue
		}

		if featuredOk && bool(c.Featured) != featured {
			continue
		}

		// Check minimums on various resources.
		if memoryOk && int(c.Memory) != memory {
			continue
		}

		if diskOk && int(c.Disk) != disk {
			continue
		}

		if vcpuOk && int(c.VCPU) != vcpu {
			continue
		}

		if networkZoneOk {
			if !c.ZoneAvailability[networkZone] {
				continue
			}
		}

		if configIDOk && int(c.ID) != configID {
			continue
		}

		filteredConfigs = append(filteredConfigs, c)
	}
	log.Printf("filteredConfigs %+v", filteredConfigs)
	return filteredConfigs
}
