
resource "ctyun_express_connect" "express_connect_dependence" {
  name        = "express_connect_dependence"
  description = "云间高速开发测试专用"
data "ctyun_vpcs" "vpc_test" {
  page_size = 50
}

locals {
  vpcs        = [for vpc in data.ctyun_vpcs.vpc_test.vpcs : vpc if vpc.name == "tf-vpc-for-ec"]
  data_vpc_id = length(local.vpcs) > 0 ? local.vpcs[0].vpc_id : ""
}

resource "ctyun_vpc" "vpc_test" {
  count       = local.data_vpc_id == "" ? 1 : 0
  name        = "tf-vpc-for-ec"
  cidr        = "192.168.0.0/16"
  description = "terraform-ec测试使用"
  enable_ipv6 = true
}

locals {
  real_vpc_id = local.data_vpc_id == "" ? try(ctyun_vpc.vpc_test[0].id, "") : local.data_vpc_id
}


data "ctyun_subnets" "subnet_test" {
  vpc_id = local.real_vpc_id
}

locals {
  subnets = [
    for subnet in data.ctyun_subnets.subnet_test.subnets : subnet if subnet.name == "tf-subnet-for-ec-1"
  ]
  data_subnet_id = length(local.subnets) > 0 ? local.subnets[0].subnet_id : ""

  subnets2 = [
    for subnet in data.ctyun_subnets.subnet_test.subnets : subnet if subnet.name == "tf-subnet-for-ec-2"
  ]
  data_subnet_id2 = length(local.subnets2) > 0 ? local.subnets2[0].subnet_id : ""
}

resource "ctyun_ec_cloud_gateway" "cloud_gateway_dependence" {
  ec_id    = ctyun_express_connect.express_connect_dependence.id
  name     = "cloud_gateway_dependence"
  description = "云间高速开发测试专用"
  region_id = "200000002401"
  region_name = "cn-hn-cs42-hncs1A-public-ctcloud"
resource "ctyun_subnet" "subnet_test" {
  count       = local.data_vpc_id=="" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-subnet-for-ec-1"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_subnet" "subnet_test2" {
  count       = local.data_vpc_id=="" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-subnet-for-ec-2"
  cidr        = "192.168.2.0/24"
  description = "terraform测试使用"
  dns = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

locals {
  real_subnet_id = local.data_subnet_id == "" ? try(ctyun_subnet.subnet_test[0].id, "") : local.data_subnet_id
  real_subnet_id2 = local.data_subnet_id2 == "" ? try(ctyun_subnet.subnet_test2[0].id, "") : local.data_subnet_id2
}
