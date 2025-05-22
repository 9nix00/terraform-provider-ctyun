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


data "ctyun_ecs_flavors" "ecs_flavor_test" {
  cpu    = 4
  ram    = 8
  arch   = "x86"
  series = "C"
  type   = "CPU_C7"
}

resource "ctyun_ccse_node_pool" "example" {
  cluster_id               = "19b4be67777e40e690b97c3a8664a1f9"
  node_pool_name           = "default-pool"
  cycle_type              = "on_demand"
  auto_renew_status        = 1
  instance_type            = "ecs"
  mirror_name             = "CTyunOS-23.01-CCND_CCSE_40_08-x86_64"
  mirror_id                = "3f80d8c0-8eb5-4afa-a506-13ba68b61872"
  mirror_type              = 1
  password                 = "P@ss2wsx"
  use_affinity_group       = true
  affinity_group_id      = "e9d3239a-207a-4006-aa84-3945265bac27"
  item_def_name            = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].name
  max_pod_num              = 110

  sys_disk = {
    type = "SATA"
    size = 3000
  }

  data_disks = [
    {
      type = "SSD"
      size = 4000
    }
  ]
}