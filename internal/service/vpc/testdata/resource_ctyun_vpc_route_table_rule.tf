resource "ctyun_vpc" "vpc_test" {
  name        = "for-route-rule"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_vpc_route_table" "route" {
  vpc_id = ctyun_vpc.vpc_test.id
  name = "route-tf-1"
}

data "ctyun_vpc_route_table_rules" "rtest" {
  route_table_id = ctyun_vpc_route_table.route.id
}

locals {
  igw_rules = [for rule in data.ctyun_vpc_route_table_rules.rtest.rules : rule if rule.next_hop_type == "igw"]
  igw_id = length(local.igw_rules) > 0 ? local.igw_rules[0].next_hop_id : ""
}

resource "ctyun_vpc_route_table_rule" "%[1]s"{
  destination = "%[2]s"
  description = "%[3]s"
  next_hop_id = local.igw_id
  next_hop_type = "igw"
  route_table_id = ctyun_vpc_route_table.route.id
  ip_version = 4
}
