
resource "ctyun_express_connect" "express_connect_dependence" {
  name        = "express_connect_dependence"
  description = "云间高速开发测试专用"

}
resource "ctyun_ec_cloud_gateway" "cloud_gateway_dependence" {
  ec_id       = ctyun_express_connect.express_connect_dependence.id
  name        = "cloud_gateway_dependence"
  description = "云间高速开发测试专用"
}
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


resource "ctyun_subnet" "subnet_test" {
  count       = local.data_vpc_id=="" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-subnet-for-ec-1"
  cidr        = "192.168.10.0/24"
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

resource "ctyun_vpc" "vpc_test_for_instance" {
  name        = "tf-vpc-for-instance"
  cidr        = "192.168.0.0/16"
  description = "terraform-ec vpc instance测试使用"
  enable_ipv6 = true
  # region_id   = "200000002401"
}


resource "ctyun_subnet" "subnet_test_for_instance" {
  # region_id   = "200000002401"
  vpc_id      = ctyun_vpc.vpc_test_for_instance.id
  name        = "tf-subnet-for-vpc_instance"
  cidr        = "192.168.1.0/24"
  description = "terraform-ec vpc instance测试使用"
  dns = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_express_connect_vpc_instance" "instance_test" {
  # region_id   = "200000002401"
  ec_id       = ctyun_express_connect.express_connect_dependence.id
  cgw_id      = ctyun_ec_cloud_gateway.cloud_gateway_dependence.id
  rtb_id      = ctyun_ec_cloud_gateway.cloud_gateway_dependence.rtb_id
  vpc_id      = ctyun_vpc.vpc_test_for_instance.id
  route_learn = 1
  route_sync  = 1
  subnets     = [ctyun_subnet.subnet_test_for_instance.id]
}


resource "ctyun_ec_cloud_gateway" "cloud_gateway_xinan1" {
  ec_id       = ctyun_express_connect.express_connect_dependence.id
  name        = "cloud_gateway_xinan1"
  description = "云间高速开发测试专用"
  region_id   = "200000002368"
  region_name = "cn-xinan1-xn1A-public-ctcloud"
}

resource "ctyun_ec_cloud_gateway" "cloud_gateway_hgh7" {
  ec_id       = ctyun_express_connect.express_connect_dependence.id
  name        = "cloud_gateway_hgh7"
  description = "云间高速开发测试专用"
  region_id   = "200000003329"
  region_name = "cn-zj-hgh7-1a-public-ctcloud"
}

resource "ctyun_ec_cloud_gateway" "cloud_gateway_huhehaote3" {
  ec_id       = ctyun_express_connect.express_connect_dependence.id
  name        = "cloud_gateway_xinan1"
  description = "云间高速开发测试专用"
  region_id   = "200000003573"
  region_name = "cn-nm-het3-1a-public-ctcloud"
}

resource "ctyun_ec_cloud_gateway" "cloud_gateway_wulumuqi7" {
  ec_id       = ctyun_express_connect.express_connect_dependence.id
  name        = "cloud_gateway_hgh7"
  description = "云间高速开发测试专用"
  region_id   = "200000004098"
  region_name = "cn-xj-urc7-1a-public-ctcloud"
}



resource "ctyun_ec_packet" "packet_test" {
  ec_id        = ctyun_express_connect.express_connect_dependence.id
  packet_name  = "packet_region_peer_test"
  bandwidth    = 10
  cycle_type   = "MONTH"
  cycle_count  = 1
}


# resource "ctyun_express_connect_region_peer" "region_peer_test" {
#   name        = "region_peer_test"
#   ec_id       = ctyun_express_connect.express_connect_dependence.id
#   src_cgw_id  = ctyun_ec_cloud_gateway.cloud_gateway_huhehaote3.id
#   dst_cgw_id  = ctyun_ec_cloud_gateway.cloud_gateway_wulumuqi7.id
#   packet_id   = ctyun_ec_packet.packet_test.id
#   rate        = 1
#   route_learn = 1
# }
