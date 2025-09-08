terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 可参考index.md，在环境变量中配置ak、sk、资源池ID、可用区名称
provider "ctyun" {
  env = "prod"
}

resource "ctyun_vpc" "vpc_test" {
  name        = "vpc-test-mq"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "subnet-test-mq"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "terraform-minchiang-mq"
  description = "terraform测试使用"
}

resource "ctyun_rabbitmq_instance" "tbidgqvfbs" {
  instance_name = "tf-rabbitmq-kkk"
  cpu_num = 4
  mem_size = 8
  node_num = 3
  zone_list = ["cn-huadong1-jsnj1A-public-ctcloud"]
  disk_type = "SSD"
  disk_size = 300
  vpc_id = ctyun_vpc.vpc_test.id
  subnet_id = ctyun_subnet.subnet_test.id
  security_group_id = ctyun_security_group.security_group_test.id
  cycle_type = "on_demand"
}
