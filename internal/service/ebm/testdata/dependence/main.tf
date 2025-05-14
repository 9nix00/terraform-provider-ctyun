resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-ebm"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-ebm"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
  enable_ipv6 = true
  type = data.ctyun_ebm_device_types.test.device_types[0].smart_nic_exist ? "common":"ebm"
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-ebm"
  description = "terraform测试使用"
}

data "ctyun_ebm_device_types" "test" {
}

data "ctyun_ebm_device_raids" "system_raid" {
  device_type = data.ctyun_ebm_device_types.test.device_types[0].device_type
  volume_type = "system"
}

data "ctyun_ebm_device_raids" "data_raid" {
  device_type = data.ctyun_ebm_device_types.test.device_types[0].device_type
  volume_type = "data"
}

data "ctyun_ebm_device_images" "test" {
  device_type = data.ctyun_ebm_device_types.test.device_types[0].device_type
  os_type = "linux"
  image_type = "standard"
}
