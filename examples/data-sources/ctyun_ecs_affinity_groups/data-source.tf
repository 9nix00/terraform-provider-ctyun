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

data "ctyun_ecs_affinity_groups" "test" {

}

output "ctyun_ecs_affinity_groups_test" {
  value = data.ctyun_ecs_affinity_groups.test
}

