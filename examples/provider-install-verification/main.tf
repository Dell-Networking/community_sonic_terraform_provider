terraform {
  required_providers {
    sonic = {
      source = "registry.terraform.io/dell/community-sonic"
    }
  }
}

provider "sonic" {
  host = "https://192.168.1.122:443"
  username = "admin"
  password = "admin123"
  insecure = true
}

data "sonic_switch" "example" {}
