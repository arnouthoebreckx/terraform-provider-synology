package provider

import (
	"context"

	"github.com/arnouthoebreckx/terraform-provider-synology/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SYNOLOGY_ADDRESS", ""),
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SYNOLOGY_USERNAME", ""),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SYNOLOGY_PASSWORD", ""),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"synology_vmm_guest":    dataSourceVmmGuest(),
			"synology_vmm_storages": datasourceVmmStorages(),
			"synology_vmm_networks": datasourceVmmNetworks(),
			"synology_vmm_hosts":    datasourceVmmHosts(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"synology_file":      resourceFileItem(),
			"synology_folder":    resourceFolderItem(),
			"synology_vmm_guest": resourceVmmGuest(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	url := d.Get("url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	var diags diag.Diagnostics
	synologyClient := client.NewClient()

	err := synologyClient.Connect(url, username, password)
	if err != nil {
		return synologyClient, diag.FromErr(err)
	}

	return synologyClient, diags
}
