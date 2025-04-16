resource "ctyun_vpc" "vpc_test" {
  name        = "vpc-test-tf"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_security_group" "%[1]s" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "%[3]s"
  description = "%[4]s"
}

data "ctyun_security_groups" "%[2]s" {
  security_group_id = ctyun_security_group.%[1]s.id
}