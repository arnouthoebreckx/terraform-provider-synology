package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/arnouthoebreckx/terraform-provider-synology/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFolderItem() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFolderCreateItem,
		ReadContext:   resourceFolderReadItem,
		UpdateContext: resourceFolderUpdateItem,
		DeleteContext: resourceFolderDeleteItem,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path of the folder",
			},
		},
	}
}

func resourceFolderCreateItem(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(client.SynologyClient)

	path := d.Get("path").(string)

	service := FolderItemService{synologyClient: client}
	err := service.Create(path)
	if err != nil {
		return diag.FromErr(err)
	}
	resourceFolderReadItem(ctx, d, m)
	return diags
}

func resourceFolderReadItem(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	path := d.Get("path").(string)

	d.Set("path", path)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func resourceFolderUpdateItem(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceFolderCreateItem(ctx, d, m)
}

func resourceFolderDeleteItem(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(client.SynologyClient)

	path := d.Get("path").(string)

	service := FolderItemService{synologyClient: client}

	err := service.Delete(path)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
