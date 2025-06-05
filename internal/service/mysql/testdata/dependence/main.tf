// main.tf负责创建或查询单测依赖的前置资源
resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-mysql3"
  cidr        = "192.168.128.0/17"
  description = "terraform测试使用2"
  enable_ipv6 = true
}


resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-mysql5"
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
  name        = "tf-subnet-for-mysql6"
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
  name        = "tf-secureity-group-for-mysql4"
  description = "terraform测试"
}

resource "ctyun_security_group" "test2" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-secureity-group-for-mysql5"
  description = "terraform测试"
}

resource "ctyun_eip" "eip_test" {
  name                = "tf-eip-for-mysql"
  bandwidth           = 1
  cycle_type          = "on_demand"
  demand_billing_type = "upflowc"
}

resource "ctyun_mysql_instance" "mysql_test" {
  bill_mode         = "2"
  prod_version      = "5.7"
  vpc_id            = ctyun_vpc.vpc_test.id
  host_type         = "C7"
  subnet_id         = ctyun_subnet.subnet_test2.id
  security_group_id = ctyun_security_group.test2.id
  name              = "tf-mysql-for-ip5"
  password          = "kqjwyk111"
  period            = 1
  purchase_count    = 1
  auto_renew_status = 0
  prod_id           = 10001003
  mysql_node_info_list = [
    {
      "node_type" : "master", "inst_spec" : "1", "storage_type" : "SATA", "storage_space" : 100,
      "prod_performance_spec" : "2C4G", "disks" : 1, "availability_zone_info" : [
      {
        "availability_zone_name" : "cn-nm-het3-1a-public-ctcloud", "availability_zone_count" : 1, "node_type" : "master"
      }
    ]
    }
  ]
  cpu_type = "30"
  os_type  = "11"
}