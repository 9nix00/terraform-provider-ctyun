resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-peer_connect"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_vpc" "vpc_test1" {
  name        = "tf-vpc-for-peer_connect1"
  cidr        = "172.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_vpc" "vpc_test3" {
  name        = "tf-vpc-for-peer_connect2"
  cidr        = "172.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_vpc_peer_connection" "test" {
  project_id     = "0"
  name           = "peer-connection_tf"
  description    = "terraform测试使用"
  request_vpc_id = ctyun_vpc.vpc_test.id
  accept_vpc_id  = ctyun_vpc.vpc_test3.id
}


provider "ctyun" {
  alias           = "test_accpet"
  ak              = "4f996ee252e84b358a4c3f9b042f646b"                                    # 如果此值不填，则默认读取环境变量中的CTYUN_AK
  sk              = "d22076b320d7424399e4235cb626f07a"                                    # 如果此值不填，则默认读取环境变量中的CTYUN_SK
  env             = "prod"                                     # 如果此值不填，则默认读取环境变量中的CTYUN_ENV
}


resource "ctyun_vpc" "vpc_test2" {
  provider    = ctyun.test_accpet
  name        = "tf-vpc-for-peer_connect"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}