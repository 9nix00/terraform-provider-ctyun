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

data "ctyun_ecs_snapshots" "test" {

}

output "ctyun_ecs_snapshots_test" {
  value = data.ctyun_ecs_snapshots.test
}

