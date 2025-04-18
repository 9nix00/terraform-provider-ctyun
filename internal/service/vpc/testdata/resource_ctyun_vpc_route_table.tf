resource "ctyun_vpc" "vpc_test" {
  name        = "for-route-table"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_vpc_route_table" "%[1]s" {
  vpc_id = ctyun_vpc.vpc_test.id
  name = "%[2]s"
}
