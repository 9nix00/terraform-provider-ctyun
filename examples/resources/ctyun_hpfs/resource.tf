terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}


resource "ctyun_hpfs" "test" {
  sfs_protocol = "hpfs"
  cycle_type = "on_demand"
  sfs_name = "hpfs-test"
  sfs_size = 512
}

