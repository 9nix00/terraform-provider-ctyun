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

data "ctyun_ebs_backups" "test" {

}

output "ctyun_ebs_backups_test" {
  value = data.ctyun_ebs_backups.test
}

