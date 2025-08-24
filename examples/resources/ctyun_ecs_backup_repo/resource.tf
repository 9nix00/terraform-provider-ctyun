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



resource "ctyun_ecs_backup_repo" "test" {
  repository_name = "test111"
  cycle_count = "5"
  cycle_type  = "MONTH"
}


