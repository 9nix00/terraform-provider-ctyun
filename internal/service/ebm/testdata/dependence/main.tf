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
  type = "common"
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-ebm"
  description = "terraform测试使用"
}

locals {
  device_type1 = "physical.s5.2xlarge4"      // az1、有本地盘、弹性、不支持云硬盘
  device_type2 = "physical.s5.2xlarge1"      // az2、无本地盘、ta
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

resource "ctyun_ebs" "ebs_test" {
  name       = "tf-ebs-for-ebm"
  mode       = "vbd"
  type       = "sata"
  size       = 60
  cycle_type = "on_demand"
  az_name   = ""
}

resource "ctyun_ebm" "ebm_test" {
  instance_name = "tf-ebm-for-ebm"
  hostname = "tf-ebm-for-ebm"
  password = "%[4]s"
  status = "%[5]s"
  ext_ip = "not_use"
  cycle_type = "on_demand"
  device_type = "%[6]s"
  image_uuid = "%[7]s"
  security_group_ids = [ctyun_security_group.security_group_test.id]
  vpc_id = "%[9]s"

  disk_list =  [{
    disk_type = "system"
    size = "100"
    type = "sata"
  }]
  network_card_list = [{
    master = true,
    subnet_id = ctyun_subnet.subnet_test.id
  }]
}

locals {
# 生成当前时间戳的哈希值
hash = sha256(timestamp())

# 从哈希结果中截取字符（转为小写并移除特殊字符）
random_string = substr(
replace(
lower(local.hash),
"/[^a-z0-9]/",
""  # 移除所有非字母数字的字符
),
0, 10  # 截取前16个字符
)
}