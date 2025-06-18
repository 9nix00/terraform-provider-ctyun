
resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-for-paas"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-for-paas"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "tf-sg-for-paas"
  description = "terraform测试使用"
  lifecycle {
    prevent_destroy = true
  }
}

resource "ctyun_mysql_instance" "mysql_test" {
  cycle_type            = "2"
  prod_version          = "5.7"
  vpc_id                =  ctyun_vpc.vpc_test.id
  host_type             = "C7"
  subnet_id             = ctyun_subnet.subnet_test.id
  security_group_id     = ctyun_security_group.security_group_test.id
  name                  = "tf-mysql"
  password              = "**********"
  cycle_count           = 1
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
    { "availability_zone_name" : "cn-huadong1-jsnj1A-public-ctcloud", "availability_zone_count" : 1, "node_type" : "master" }
  ]
  cpu_type = "30"
  os_type  = "11"
}