resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-vpc"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}