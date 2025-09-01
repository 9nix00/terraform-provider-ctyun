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


data "ctyun_scaling_ecs_list" "scaling_ecs_list" {
  group_id = 109737
}
