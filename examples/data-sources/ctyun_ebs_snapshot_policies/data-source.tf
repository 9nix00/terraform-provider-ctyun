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

data "ctyun_ebs_snapshot_policies" "test" {

}

output "ctyun_ebs_snapshot_policies" {
  value = data.ctyun_ebs_snapshot_policies.test
}

