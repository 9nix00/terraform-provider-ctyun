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

data "ctyun_mysql_specs" "mysql_specs" {
  instance_series = "S"
}

resource "ctyun_mysql_instance" "mysql_test" {
  cycle_type            = "on_demand"
  prod_version          = "5.7"
  vpc_id                = local.real_vpc_id
  host_type             = "S7"
  subnet_id             = local.real_subnet_id
  security_group_id     = local.real_security_group_id
  name                  = local.mysql_name
  prod_id               = "Single57"
  instance_series       = "S"
  storage_type          = "SATA"
  storage_space         = 100
  prod_performance_spec = "2C4G"
  availability_zone_info = [
    { "availability_zone_name" : local.az_name, "availability_zone_count" : 1, "node_type" : "master" }
  ]
  cpu_type = "Intel"
  os_type  = "ctyunos"
}