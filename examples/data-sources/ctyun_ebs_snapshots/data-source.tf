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

data "ctyun_ebs_snapshots" "test" {

}

output "ctyun_ebs_snapshots_test" {
  value = data.ctyun_ebs_snapshots.test
}

