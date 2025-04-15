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

resource "ctyun_vpc" "vpc_test" {
  name        = "vpc-test-mc1"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
  region_id   = "200000001852"
}

data "ctyun_vpcs" "test" {
  region_id = "200000001852"
  vpc_id = ctyun_vpc.vpc_test.id
  # page_no = 1
  # page_size = 1
}

output "ctyun_test" {
  value = data.ctyun_vpcs.test
}

