resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-iaas"
  cidr        = "192.168.0.0/16"
  description = "terraform-iaas测试使用"
  enable_ipv6 = true
}

resource "ctyun_vpc" "vpc_test1" {
  name        = "tf-vpc-for-iaas1"
  cidr        = "192.168.0.0/16"
  description = "terraform-iaas测试使用1"
  enable_ipv6 = true
}

resource "ctyun_vpc" "vpc_test2" {
  name        = "tf-vpc-for-iaas2"
  cidr        = "192.168.0.0/16"
  description = "terraform-iaas测试使用2"
  enable_ipv6 = true
}


resource "ctyun_vpc" "vpc_test3" {
  name        = "tf-vpc-for-iaas3"
  cidr        = "192.168.0.0/16"
  description = "terraform-iaas测试使用3"
  enable_ipv6 = true
}


resource "ctyun_vpc" "vpc_test4" {
  name        = "tf-vpc-for-iaas4"
  cidr        = "192.168.0.0/16"
  description = "terraform-iaas测试使用4"
  enable_ipv6 = true
}

resource "ctyun_vpc" "vpc_test5" {
  name        = "tf-vpc-for-iaas5"
  cidr        = "192.168.0.0/16"
  description = "terraform-iaas测试使用5"
  enable_ipv6 = true
}


resource "ctyun_private_zone" "zone_test" {
name          = "zone.test.com"
description   = "terraform前置资源"
proxy_pattern = "zone"
ttl           = 300
vpc_id_list   = [ctyun_vpc.vpc_test1.id, ctyun_vpc.vpc_test2.id]
}







