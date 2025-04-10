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

data "ctyun_ecs_instances" "test" {
  page_size = 1
}

output "ctyun_ecs_instances_test" {
  value = data.ctyun_ecs_instances.test
}