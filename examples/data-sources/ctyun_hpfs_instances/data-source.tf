terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

provider "ctyun" {
  env = "prod"
}

data "ctyun_hpfs_instances" "test" {
  sfs_status = "available"
  sfs_protocol = "hpfs"
  az_name = "bb9fdb42056f11eda1610242ac110002"
}
