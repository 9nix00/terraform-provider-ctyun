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
}

resource "ctyun_vpc_route_table" "route" {
  vpc_id = ctyun_vpc.vpc_test.id
  name = "route-t1f"
}

data "ctyun_vpc_route_table_rules" "rtest" {
  route_table_id = ctyun_vpc_route_table.route.route_table_id
}

