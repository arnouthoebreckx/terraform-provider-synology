package provider

import (
	"context"
	"github.com/arnouthoebreckx/terraform-provider-synology/client"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceVmmStorages() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceVmmStoragesRead,
		Schema: map[string]*schema.Schema{
			"storages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"used": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volume_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceVmmStoragesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(client.SynologyClient)
	service := StorageGuestService{synologyClient: client}

	storageResponse, err := service.Read()
	if err != nil {
		return diag.FromErr(err)
	}

	var storages []map[string]interface{}
	for _, s := range storageResponse.Storages {
		m := map[string]interface{}{
			"host_id":      s.HostID,
			"host_name":    s.HostName,
			"size":         s.Size,
			"status":       s.Status,
			"storage_id":   s.StorageID,
			"storage_name": s.StorageName,
			"used":         s.Used,
			"volume_path":  s.VolumePath,
		}
		storages = append(storages, m)
	}

	if err := d.Set("storages", storages); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
