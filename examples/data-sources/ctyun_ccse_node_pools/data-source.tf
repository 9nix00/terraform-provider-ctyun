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

data "ctyun_ccse_node_pools" "test" {
  cluster_id               = "19b4be67777e40e690b97c3a8664a1f9"
}

output "t" {
  value = data.ctyun_ccse_node_pools.test
}