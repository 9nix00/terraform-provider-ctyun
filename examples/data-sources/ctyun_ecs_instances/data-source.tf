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

data "ctyun_ecs_instances" "test" {
  page_size = 1
}

output "ctyun_ecs_instances_test" {
  value = data.ctyun_ecs_instances.test
}