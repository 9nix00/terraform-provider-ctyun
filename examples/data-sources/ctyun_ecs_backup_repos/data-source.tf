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

data "ctyun_ecs_backup_repos" "test" {

}


output "ctyun_ecs_backup_repos_test" {
  value = data.ctyun_ecs_backup_repos.test
}

