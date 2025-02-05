# Credits

Fork from https://github.com/sergief/terraform-provider-synology, started from his project to expand and add other APIs aswell.

# Terraform Provider Synology

Provider for managing Synology Resources.

## Build and install
### Using the local environment
If you have the `Go` environment installed, you can simply run the makefile:
```bash
make clean build test release
```
### Using the docker-compose image
You don't need to install the golang development environment if you already have a working `docker` environment, simply run:
```bash
docker-compose build --no-cache && docker-compose run app make clean build test release
```

After this, pick up the version specifically compiled for your OS and architecture from `./bin` and put it in `$HOME/.terraform.d/plugins/github.com/arnouthoebreckx/synology/0.1/$OS_$ARCH/terraform-provider-synology`


## Acceptance Tests
Run the following command setting the required environment variables (no docker support):
```bash
SYNOLOGY_ADDRESS=http://aaa.bbb.ccc.dddd:5000 SYNOLOGY_USERNAME=test_user SYNOLOGY_PASSWORD=test_password make testacc
```

## Terraform Resources

```terraform
terraform {
  required_providers {
    synology = {
      version = "0.2.0"
      source = "github.com/arnouthoebreckx/synology"
    }
  }
}

provider "synology" {
    url = "http://192.168.1.5:5000"
    username = "testuser"
    password = "testpass"
    # these variables can be set as env vars in SYNOLOGY_ADDRESS SYNOLOGY_USERNAME and SYNOLOGY_PASSWORD
}
```
## Virtual Machine Manager
This resource creates a text file in a Synology Filestation.
Example:

```
resource "synology_vmm_guest" "my-guest" {
  guest_name = "terraform-guest"
  storage_name = "storage1"
  vnics {
    network_name = "default"
  }
  vdisks {
    create_type = 0
    vdisk_size = 10240
  }
}

```

### File

This resource creates a text file in a Synology Filestation.
Example:

```
resource "synology_file" "hello-world" {
  filename = "/home/downloaded/hello-world.txt"
  content = "Hello World"
}
```

### Folder

This resource creates a folder in a Synology Filestation.
Example:
```terraform
terraform {
  required_providers {
    synology = {
      version = "0.2.0"
      source = "github.com/arnouthoebreckx/synology"
    }
  }
}

provider "synology" {
    url = "http://192.168.1.5:5000"
    username = "test"
    password = "test"
    # these variables can be set as env vars in SYNOLOGY_ADDRESS SYNOLOGY_USERNAME and SYNOLOGY_PASSWORD
}

resource "synology_folder" "my-folder" {
  path = "/home/downloaded/sample-folder"
}
```

## Known issues

### Vdisk size after creation seems to always follow the closest power of 2 that is greater than or equal to input
For example.

if you put

```
  vdisks {
    create_type = 0
    vdisk_size  = 20000
  }
```

After another plan it will propose the change to 

```
vdisk_size  = 20480 -> 20000
```

Best way around this is by creating with the size you want and then you either add:

```
ignore_changes = [
  vdisk_size
]
```

or you can remove the resource and update the disk_size to what was proposed.