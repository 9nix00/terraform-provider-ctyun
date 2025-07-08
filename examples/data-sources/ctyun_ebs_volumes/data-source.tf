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
  disk_name = "lx-test-29"
}

output "ctyun_ebs_volumes_test" {
  value = data.ctyun_ebs_volumes.test
}

