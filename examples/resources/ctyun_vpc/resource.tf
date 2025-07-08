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
  name        = "vpca-ccs"
  cidr        = "10.0.0.0/8"
  description = "terraform测试使用"
  enable_ipv6 = true
}