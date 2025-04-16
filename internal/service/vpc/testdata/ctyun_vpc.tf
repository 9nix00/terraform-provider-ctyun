resource "ctyun_vpc" "%[1]s" {
  name        = "%[3]s"
  cidr        = "%[5]s"
  description = "%[4]s"
  enable_ipv6 = true
}

data "ctyun_vpcs" "%[2]s" {
  vpc_id = ctyun_vpc.%[1]s.id
}