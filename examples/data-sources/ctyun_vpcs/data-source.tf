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

data "ctyun_vpcs" "test" {
  page_no = 1
}

output "ctyun_test" {
  value = data.ctyun_vpcs.test
}

