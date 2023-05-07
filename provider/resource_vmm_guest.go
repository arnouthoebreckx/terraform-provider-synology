package provider

import (
	"context"
	"fmt"
	"github.com/arnouthoebreckx/terraform-provider-synology/client"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVmmGuest() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGuestCreateItem,
		ReadContext:   resourceGuestReadItem,
		UpdateContext: resourceGuestUpdateItem,
		DeleteContext: resourceGuestDeleteItem,
		Schema: map[string]*schema.Schema{
			"guest_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The guest name",
			},
			"guest_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The guest name",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. The description of the guest.",
			},
			"poweron": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Optional. Default VM is not powered on.",
			},
			"autorun": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Optional. 0: off 1: last state 2: on",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The guest status. (running/shutdown/inaccessiblen/booting/shutting_down/moving/stor_migrating/creating/importing/preparing/ha_standby/unknown/crashed/undefined",
			},
			"storage_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional. The name of storage where the guest resides. Note: At least storage_id or storage_name should be given.",
			},
			"storage_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Optional. The id of storage where the guest resides. Note: At least storage_id or storage_name should be given.",
			},
			"vcpu_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Optional. The vCPU number",
			},
			"vram_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1024,
				Description: "Optional. The memory size in MB",
			},
			"vnics": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mac": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Optional. MAC address. If not specified, a MAC address of this vNIC will be randomly generated.",
						},
						"model": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "1: VirtIO 2: e1000 3: rtl8139",
						},
						"network_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Optional. Connected network group id. At least network_id or network_name should be given. Note: network_id can be an empty string to represent not being connected.",
						},
						"network_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Optional. Connected network group name. At least network_id or network_name should be given.",
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
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_type": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    true,
							Description: "0: Create an empty vDisk, 1: Clone an existing image",
						},
						"vdisk_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "Optional. If create_type is 0, this field must be set. The created vDisk size in MB.",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								if size, ok := val.(int); ok {
									if size < 10240 {
										errs = append(errs, fmt.Errorf("%q must be greater than 10240", key))
									}
								}
								return
							},
						},
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							ForceNew:    true,
							Description: "Optional. If create_type is 1, at least image_id or image_name should be given. The id of the image that is to be cloned. Note: Image type should be disk.",
						},
						"image_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Optional. If create_type is 1, at least image_id or image_name should be given. The name of the image that is to be cloned. Note: Image type should be disk.",
						},
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
					},
				},
			},
		},
	}
}

func mapFromGuestToData(d *schema.ResourceData, guest client.Guest) {
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set("autorun", fmt.Sprint(guest.Autorun))
	d.Set("description", guest.Description)
	d.Set("guest_id", guest.GuestId)
	d.Set("guest_name", guest.GuestName)
	d.Set("status", guest.Status)
	d.Set("storage_id", guest.StorageId)
	d.Set("storage_name", guest.StorageName)
	d.Set("vcpu_num", guest.VcpuNum)

	vdisks := make([]interface{}, len(guest.Vdisks))
	for i, vdisk := range guest.Vdisks {
		vdiskMap := make(map[string]interface{})
		vdiskMap["controller"] = vdisk.Controller
		vdiskMap["unmap"] = vdisk.Unmap
		vdiskMap["vdisk_id"] = vdisk.VdiskId
		vdiskMap["vdisk_size"] = vdisk.VdiskSize
		vdisks[i] = vdiskMap
	}
	d.Set("vdisks", vdisks)

	vnics := make([]interface{}, len(guest.Vnics))
	for i, vnic := range guest.Vnics {
		vnicMap := make(map[string]interface{})
		vnicMap["mac"] = vnic.Mac
		vnicMap["model"] = vnic.Model
		vnicMap["network_id"] = vnic.NetworkID
		vnicMap["network_name"] = vnic.NetworkName
		vnicMap["vnic_id"] = vnic.VnicID
		vnics[i] = vnicMap
	}

	d.Set("vnics", vnics)
	d.Set("vram_size", guest.VramSize)

}

func resourceGuestCreateItem(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(client.SynologyClient)

	name := d.Get("guest_name").(string)
	storageID := d.Get("storage_id").(string)
	storageName := d.Get("storage_name").(string)
	validateIDName(storageID, storageName)

	vnics := removeEmptyEntries(d.Get("vnics").([]interface{}))
	vdisks := removeEmptyEntries(d.Get("vdisks").([]interface{}))
	validateListIDName(vnics, "network_id", "network_name")
	validateListIDName(vdisks, "image_id", "image_name")

	service := GuestService{synologyClient: client}

	err := service.Create(name, storageID, storageName, vnics, vdisks)
	if err != nil {
		return diag.FromErr(err)
	}

	autorun := d.Get("autorun").(int)
	description := d.Get("description").(string)
	vcpuNum := d.Get("vcpu_num").(int)
	vramSize := d.Get("vram_size").(int)
	err = service.Set(name, "", autorun, description, vcpuNum, vramSize)
	if err != nil {
		return diag.FromErr(err)
	}

	powerOn := d.Get("poweron").(bool)
	err = service.Power(name, powerOn)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceGuestReadItem(ctx, d, m)

	return diags
}

func resourceGuestReadItem(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(client.SynologyClient)
	service := GuestService{synologyClient: client}

	name := d.Get("guest_name").(string)

	guest, err := service.Read(name)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("Guest information: %+v", guest)

	mapFromGuestToData(d, guest)

	return diags
}

func resourceGuestUpdateItem(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(client.SynologyClient)
	service := GuestService{synologyClient: client}

	name := d.Get("guest_name").(string)
	oldName := name
	if d.HasChange("guest_name") {
		old, new := d.GetChange("guest_name")
		oldName = old.(string)
		name = new.(string)
	}

	autorun := d.Get("autorun").(int)
	description := d.Get("description").(string)
	vcpuNum := d.Get("vcpu_num").(int)
	vramSize := d.Get("vram_size").(int)

	// Turn off VM to make changes to name
	err := service.Power(oldName, false)
	if err != nil {
		return diag.FromErr(err)
	}

	// Add a sleep here to give time for the guest to shut down
	time.Sleep(10 * time.Second)

	err = service.Set(oldName, name, autorun, description, vcpuNum, vramSize)
	if err != nil {
		return diag.FromErr(err)
	}

	// Add a sleep here to give time for the guest to update
	time.Sleep(2 * time.Second)

	powerOn := d.Get("poweron").(bool)
	err = service.Power(name, powerOn)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceGuestReadItem(ctx, d, m)

	return diags
}

func resourceGuestDeleteItem(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(client.SynologyClient)
	service := GuestService{synologyClient: client}
	name := d.Get("guest_name").(string)

	// Incase of recreate turn off the VM
	err := service.Power(name, false)
	if err != nil {
		return diag.FromErr(err)
	}
	time.Sleep(10 * time.Second)

	err = service.Delete(name)

	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func getShutdownNeeded(d *schema.ResourceData, shutdown_updates []string) bool {
	if d.HasChanges(shutdown_updates...) {
		return true
	}
	return false
}
