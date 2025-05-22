terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}


provider "ctyun" {
  region_id  = "bb9fdb42056f11eda1610242ac110002"
  az_name    = "cn-huadong1-jsnj1A-public-ctcloud"
}

data "ctyun_ecs_affinity_groups" "test" {

}

output "ctyun_ecs_affinity_groups_test" {
  value = data.ctyun_ecs_affinity_groups.test
}

