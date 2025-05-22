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

data "ctyun_ebs_volumes" "test" {
  region_id = "200000001852"
}

output "ctyun_ebs_volumes_test" {
  value = data.ctyun_ebs_volumes.test
}

