---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "synology_vmm_guest Data Source - terraform-provider-synology"
subcategory: ""
description: |-
  
---

# synology_vmm_guest (Data Source)



## Example Usage

```terraform
terraform {
  required_providers {
    synology = {
      version = "0.1"
      source  = "github.com/arnouthoebreckx/synology"
    }
  }
}

provider "synology" {
  url      = "<SYNOLOGY_ADDRESS>"
  username = "<SYNOLOGY_USERNAME>"
  password = "<SYNOLOGY_PASSWORD>"
  # these variables can be set as env vars in SYNOLOGY_ADDRESS SYNOLOGY_USERNAME and SYNOLOGY_PASSWORD
}

resource "synology_vmm_guest" "my-guest" {
  autorun     = 2
  guest_name   = "terraform-guest"
  description  = "Virtual machine setup with terraform"
  storage_name = "synology - VM Storage 1"
  vram_size    = 1024
  vnics {
    network_name = "default"
  }
  vdisks {
    create_type = 0
    vdisk_size  = 10240
  }
}

data "synology_vmm_guest" "my-guest" {
  guest_name = "terraform-guest-tmp"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `guest_name` (String) The name of this guest

### Read-Only

- `auto_run` (Number) 0: off 1: last state 2: on
- `description` (String) The description of the guest.
- `guest_id` (String) The id of this guest.
- `id` (String) The ID of this resource.
- `status` (String) The guest status. (running/shutdown/inaccessiblen/booting/shutting_down/moving/stor_migrating/creating/importing/preparing/ha_standby/unknown/crashed/undefined
- `storage_id` (String) The id of storage where the guest resides.
- `storage_name` (String) The name of storage where the guest resides.
- `vcpu_num` (Number) The number of vCPU.
- `vdisks` (List of Object) (see [below for nested schema](#nestedatt--vdisks))
- `vnics` (List of Object) (see [below for nested schema](#nestedatt--vnics))
- `vram_size` (Number) The memory size of this guest in MB.

<a id="nestedatt--vdisks"></a>
### Nested Schema for `vdisks`

Read-Only:

- `controller` (Number)
- `unmap` (Boolean)
- `vdisk_id` (String)
- `vdisk_size` (Number)


<a id="nestedatt--vnics"></a>
### Nested Schema for `vnics`

Read-Only:

- `mac` (String)
- `model` (Number)
- `network_id` (String)
- `network_name` (String)
- `vnic_id` (String)

