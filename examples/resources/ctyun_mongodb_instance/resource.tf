terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}
provider "ctyun" {
  region_id = "bb9fdb42056f11eda1610242ac110002"
  az_name   = "cn-huadong1-jsnj1A-public-ctcloud"
}
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
    prevent_destroy = false
  }
}

variable "password" {
  type      = string
  sensitive = true
}

// 创建单节点
resource "ctyun_mongodb_instance" "mongodb_test" {
  cycle_type        = "on_demand"
  vpc_id            = ctyun_vpc.vpc_test.id
  host_type         = "S7"
  subnet_id         = ctyun_subnet.subnet_test.id
  security_group_id = ctyun_security_group.security_group_test.id
  name              = "mongodb_test"
  password          = var.password
  prod_id           = "Single34"
  node_info_list = [
    {
      "node_type" : "s", "instance_series" : "S", "storage_type" : "SATA", "storage_space" : 100,
      "prod_performance_spec" : "2C4G", "availability_zone_info" : [
      {
        "availability_zone_name" : "cn-huadong1-jsnj1A-public-ctcloud", "availability_zone_count" : 1,
        "node_type" : "master"
      }
    ]
    }, {
      "node_type" : "backup", "instance_series" : "S", "storage_type" : "SATA", "storage_space" : 100,
      "prod_performance_spec" : "2C4G", "disks" : 1, "availability_zone_info" : [
        {
          "availability_zone_name" : "cn-huadong1-jsnj1A-public-ctcloud", "availability_zone_count" : 1,
          "node_type" : "backup"
        }
      ]
    }
  ]
}

// 升配磁盘
resource "ctyun_mongodb_instance" "mongodb_test" {
  cycle_type        = "on_demand"
  vpc_id            = ctyun_vpc.vpc_test.id
  host_type         = "S7"
  subnet_id         = ctyun_subnet.subnet_test.id
  security_group_id = ctyun_security_group.security_group_test.id
  name              = "mongodb_test"
  password          = var.password
  prod_id           = "Single34"
  node_info_list = [
    {
      "node_type" : "master", "instance_series" : "S", "storage_type" : "SATA", "storage_space" : 120,
      "prod_performance_spec" : "2C4G", "availability_zone_info" : [
      {
        "availability_zone_name" : "cn-huadong1-jsnj1A-public-ctcloud", "availability_zone_count" : 1,
        "node_type" : "master"
      }
    ]
    }
  ]
  is_upgrade_back_up = true
}