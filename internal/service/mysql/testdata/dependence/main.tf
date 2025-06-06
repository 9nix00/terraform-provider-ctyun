// main.tf负责创建或查询单测依赖的前置资源

locals {
  vpc_name             = "tf-vpc-for-mysql-${local.random_string}"
  subnet_name1         = "tf-subnet-for-mysql-${local.random_string}1"
  subnet_name2         = "tf-subnet-for-mysql-${local.random_string}2"
  security_group_name1 = "tf-sg-for-mysql-${local.random_string}1"
  security_group_name2 = "tf-sg-for-mysql-${local.random_string}2"
  mysql_name = "tf-mysql-for-ip-${local.random_string}"
}
resource "ctyun_vpc" "vpc_test" {
  name        = local.vpc_name
  cidr        = "192.168.128.0/17"
  description = "terraform测试使用2"
  enable_ipv6 = true
}


resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = local.subnet_name1
  cidr        = "192.168.192.0/24"
  description = "terraform测试使用3"
  dns = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
  enable_ipv6 = true
}
resource "ctyun_subnet" "subnet_test2" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = local.subnet_name2
  cidr        = "192.168.193.0/24"
  description = "terraform测试使用4"
  dns = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
  enable_ipv6 = true
}

resource "ctyun_security_group" "test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = local.security_group_name1
  description = "terraform测试"
}

resource "ctyun_security_group" "test2" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = local.security_group_name2
  description = "terraform测试"
}

resource "ctyun_eip" "eip_test" {
  name                = "tf-eip-for-mysql"
  bandwidth           = 1
  cycle_type          = "on_demand"
  demand_billing_type = "upflowc"
}

resource "ctyun_mysql_instance" "mysql_test" {
  bill_mode             = "2"
  prod_version          = "5.7"
  vpc_id                = ctyun_vpc.vpc_test.id
  host_type             = "C7"
  subnet_id             = ctyun_subnet.subnet_test2.id
  security_group_id     = ctyun_security_group.test2.id
  name                  = local.mysql_name
  password              = "kqjwyk111"
  period                = 1
  purchase_count        = 1
  auto_renew_status     = 0
  prod_id               = 10001003
  node_type             = "master"
  inst_spec             = "1"
  storage_type          = "SATA"
  storage_space         = 100
  prod_performance_spec = "2C4G"
  disks                 = 1
  availability_zone_info = [
    { "availability_zone_name" : "cn-nm-het3-1a-public-ctcloud", "availability_zone_count" : 1, "node_type" : "master" }
  ]
  cpu_type = "30"
  os_type  = "11"
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
    0, 5  # 截取前10个字符
  )
}