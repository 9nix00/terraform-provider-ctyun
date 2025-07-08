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

resource "ctyun_ecs_affinity_group" "test" {
  affinity_group_name = "tf-test-group"
  affinity_group_policy = "anti-affinity"
}
