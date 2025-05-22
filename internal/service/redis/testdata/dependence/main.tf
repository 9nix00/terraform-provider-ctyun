resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-redis"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-redis"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
  enable_ipv6 = true
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-redis"
  description = "terraform测试使用"
}

resource "ctyun_eip" "eip_test" {
  name                = "tf-eip-for-redis"
  bandwidth           = 10
  cycle_type          = "on_demand"
  demand_billing_type = "bandwidth"
}

data "ctyun_redis_specs" "test"{

}

locals {
  spec = data.ctyun_redis_specs.test.series_infos[0]
}