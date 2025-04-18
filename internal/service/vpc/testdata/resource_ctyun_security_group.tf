resource "ctyun_vpc" "vpc_test" {
  name        = "vpc-test-tf"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_security_group" "%[1]s" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "%[2]s"
  description = "%[3]s"
}
