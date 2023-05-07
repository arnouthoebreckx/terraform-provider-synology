terraform {
  required_providers {
    synology = {
      version = "0.2.0"
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

data "synology_vmm_storages" "my-guest" {
}