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

data "ctyun_ecs_backup_policies" "test" {

}


output "ctyun_ecs_backup_policies_test" {
  value = data.ctyun_ecs_backup_policies.test
}



