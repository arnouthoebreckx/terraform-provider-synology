package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/arnouthoebreckx/terraform-provider-synology/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFileItem() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFileCreateItem,
		ReadContext:   resourceFileReadItem,
		UpdateContext: resourceFileUpdateItem,
		DeleteContext: resourceFileDeleteItem,
		Schema: map[string]*schema.Schema{
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Contents of the file",
			},
			"filename": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The filename including path",
			},
		},
	}
}

func resourceFileCreateItem(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(client.SynologyClient)

	content := d.Get("content").(string)
	filename := d.Get("filename").(string)

	service := FileItemService{synologyClient: client}
	err := service.Create(filename, []byte(content))
	if err != nil {
		return diag.FromErr(err)
	}
	resourceFileReadItem(ctx, d, m)
	return diags
}

func resourceFileReadItem(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(client.SynologyClient)

	filename := d.Get("filename").(string)

	service := FileItemService{synologyClient: client}

	content, err := service.Read(filename)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("filename", filename)
	d.Set("content", string(content))
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func resourceFileUpdateItem(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceFileCreateItem(ctx, d, m)
}

func resourceFileDeleteItem(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(client.SynologyClient)

	filename := d.Get("filename").(string)

	service := FileItemService{synologyClient: client}

	err := service.Delete(filename)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
