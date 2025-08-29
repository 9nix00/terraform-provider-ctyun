data "ctyun_vpcs" "vpc_test" {
  page_size = 50
}

locals {
  vpcs        = [for vpc in data.ctyun_vpcs.vpc_test.vpcs : vpc if vpc.name == "tf-vpc-for-paas"]
  data_vpc_id = length(local.vpcs) > 0 ? local.vpcs[0].vpc_id : ""

  vpcs_sfs = [for vpc in data.ctyun_vpcs.vpc_test.vpcs : vpc if vpc.name == "tf-vpc-for-iaas"]
  data_vpc_id_sfs = length(local.vpcs_sfs) > 0 ? local.vpcs_sfs[0].vpc_id : ""
}

resource "ctyun_vpc" "vpc_test" {
  count       = local.data_vpc_id == "" ? 1 : 0
  name        = "tf-vpc-for-paas"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_vpc" "vpc_test_iaas" {
  count       = local.data_vpc_id_sfs == "" ? 1 : 0
  name        = "tf-vpc-for-iaas"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}


locals {
  real_vpc_id = local.data_vpc_id == "" ? try(ctyun_vpc.vpc_test[0].id, "") : local.data_vpc_id
  read_iaas_vpc_id = local.data_vpc_id_sfs == ""? try(ctyun_vpc.vpc_test_iaas[0].id, ""):local.data_vpc_id_sfs
}


data "ctyun_subnets" "subnet_test" {
  vpc_id = local.real_vpc_id
}

locals {
  subnets = [
    for subnet in data.ctyun_subnets.subnet_test.subnets : subnet if subnet.name == "tf-subnet-for-paas"
  ]
  data_subnet_id = length(local.subnets) > 0 ? local.subnets[0].subnet_id : ""
}

resource "ctyun_subnet" "subnet_test" {
  count       = local.data_vpc_id=="" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-subnet-for-paas"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

locals {
  real_subnet_id = local.data_subnet_id == "" ? try(ctyun_subnet.subnet_test[0].id, "") : local.data_subnet_id
}


resource "ctyun_sfs" "sfs_test" {
  sfs_type     = "capacity"
  sfs_protocol = "nfs"
  name         = "sfs-for-group"
  sfs_size     = 500
  cycle_type   = "on_demand"
  vpc_id       = local.real_vpc_id
  subnet_id    = local.real_subnet_id
}

resource "ctyun_sfs_permission_group" "group_test" {
  name = "sfs-test1"
  description = "单元测试1"
}

resource "ctyun_sfs_permission_group" "group_test1" {
  name = "sfs-test2"
  description = "单元测试2"
}