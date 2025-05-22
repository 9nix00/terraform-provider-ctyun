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


resource "ctyun_ecs_affinity_group" "test" {
  affinity_group_name = "tf-test-group"
  affinity_group_policy = "anti-affinity"
}
