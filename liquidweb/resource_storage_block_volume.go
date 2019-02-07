package liquidweb

import (
	storage "git.liquidweb.com/masre/liquidweb-go/storage"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceStorageBlockVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateBlockVolume,
		Read:   resourceReadBlockVolume,
		Update: resourceUpdateBlockVolume,
		Delete: resourceDeleteBlockVolume,

		Schema: map[string]*schema.Schema{
			"attach": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"cross_attach": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"attached_to": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"uniq_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"zone": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"zone_availability": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceCreateBlockVolume(d *schema.ResourceData, m interface{}) error {
	opts := buildBlockVolumeOpts(d, m)
	config := m.(*Config)

	result, err := config.LWAPI.StorageBlockVolume.Create(opts)
	if err != nil {
		return err
	}

	if result.HasError() {
		return result
	}

	d.SetId(result.UniqID)

	return resourceReadBlockVolume(d, m)
}

func resourceReadBlockVolume(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	blockVolumeItem := blockVolumeDetails(config, d.Id())
	if blockVolumeItem.HasError() {
		return blockVolumeItem
	}

	updateBlockVolumeResource(d, blockVolumeItem)
	return nil
}

func resourceUpdateBlockVolume(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	// Resize if size has changed.
	if d.HasChange("size") {
		size := d.Get("size").(int)

		resizeOpts := &storage.BlockVolumeParams{
			NewSize: size,
			UniqID:  d.Id(),
		}

		blockVolumeResize := config.LWAPI.StorageBlockVolume.Resize(resizeOpts)
		if blockVolumeResize.HasError() {
			return blockVolumeResize
		}
	}

	updateOps := &storage.BlockVolumeParams{
		CrossAttach: d.Get("cross_attach").(bool),
		Domain:      d.Get("domain").(string),
		UniqID:      d.Id(),
	}

	blockVolumeItem := config.LWAPI.StorageBlockVolume.Update(updateOps)
	if blockVolumeItem.HasError() {
		return blockVolumeItem
	}

	updateBlockVolumeResource(d, blockVolumeItem)
	return nil
}

func resourceDeleteBlockVolume(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	params := &storage.BlockVolumeParams{UniqID: d.Id()}

	deleteResponse := config.LWAPI.StorageBlockVolume.Delete(params)
	if deleteResponse.HasError() {
		return deleteResponse
	}

	return nil
}

// buildBlockVolumeOpts builds options for a create/update block volume API call.
func buildBlockVolumeOpts(d *schema.ResourceData, m interface{}) *storage.BlockVolumeParams {
	params := &storage.BlockVolumeParams{
		Attach:      d.Get("attach").(string),
		CrossAttach: d.Get("cross_attach").(bool),
		Domain:      d.Get("domain").(string),
		Region:      d.Get("region").(int),
		Size:        d.Get("size").(int),
		Zone:        d.Get("zone").(int),
	}

	return params
}

// blockVolumeDetails gets a dns record's details from the API.
func blockVolumeDetails(config *Config, id string) *storage.BlockVolumeItem {
	return config.LWAPI.StorageBlockVolume.Details(id)
}

// updateBlockVolumeResource updates the resource data for the DNS Record.
func updateBlockVolumeResource(d *schema.ResourceData, dr *storage.BlockVolumeItem) {
	d.Set("cross_attach", dr.CrossAttach)
	d.Set("domain", dr.Domain)
	d.Set("size", dr.Size)
	d.Set("attached_to", dr.AttachedTo)
	d.Set("label", dr.Label)
	d.Set("status", dr.Status)
	d.Set("uniq_id", dr.UniqID)
	d.Set("zone_availability", dr.ZoneAvailability)
}
