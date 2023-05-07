package provider

import (
	"context"
	"github.com/arnouthoebreckx/terraform-provider-synology/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVmmGuest() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVmmGuestRead,
		Schema: map[string]*schema.Schema{
			"auto_run": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "0: off 1: last state 2: on",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the guest.",
			},
			"guest_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of this guest.",
			},
			"guest_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of this guest",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The guest status. (running/shutdown/inaccessiblen/booting/shutting_down/moving/stor_migrating/creating/importing/preparing/ha_standby/unknown/crashed/undefined",
			},
			"storage_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of storage where the guest resides.",
			},
			"storage_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of storage where the guest resides.",
			},
			"vcpu_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of vCPU.",
			},
			"vram_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The memory size of this guest in MB.",
			},
			"vnics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mac": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "MAC address of this vNIC.",
						},
						"model": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "1: VirtIO 2: e1000 3: rtl8139",
						},
						"network_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the network group which this vNIC connects to.",
						},
						"network_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the network group which this vNIC connects to.",
						},
						"vnic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of this vNIC.",
						},
					},
				},
			},
			"vdisks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"controller": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "1: VirtIO 2: IDE 3: SATA",
						},
						"unmap": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Determine whether to enable space reclamation.",
						},
						"vdisk_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of this vDisk.",
						},
						"vdisk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The vDisk's size of this guest in MB.",
						},
					},
				},
			},
		},
	}
}

func dataSourceVmmGuestRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(client.SynologyClient)
	service := GuestService{synologyClient: client}
	name := d.Get("guest_name").(string)

	guest, err := service.Read(name)
	if err != nil {
		return diag.FromErr(err)
	}

	mapFromGuestToData(d, guest)

	return diags
}
