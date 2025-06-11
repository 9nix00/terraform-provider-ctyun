resource "ctyun_eip" "eip_test" {
  name                = "tf-eip-for-nat"
  bandwidth           = 1
  cycle_type          = "on_demand"
  demand_billing_type = "upflowc"
}

resource "ctyun_eip" "eip_test1" {
  name                = "tf-eip-for-nat1"
  bandwidth           = 1
  cycle_type          = "on_demand"
  demand_billing_type = "upflowc"
}

resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-nat"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_nat" "nat_test"{
  vpc_id = ctyun_vpc.vpc_test.id
  spec = 1
  name = "tf-nat"
  description = "terraform测试使用"
  cycle_type = "on_demand"
}

resource "ctyun_subnet" "subnet_test1" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-nat1"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_subnet" "subnet_test2" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-nat2"
  cidr        = "192.168.128.0/24"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
}