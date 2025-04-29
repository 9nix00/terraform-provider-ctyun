terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

provider "ctyun" {
  env = "prod"
}

resource "ctyun_vpc" "vpc_test" {
  name        = "tf-vpc-test-qqq"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id = ctyun_vpc.vpc_test.id
  name        = "tf-subnet-test"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
  enable_ipv6 = true
}

data "ctyun_images" "image_test1" {
  name       = "CtyunOS 23"
  visibility = "public"
  page_no = 1
  page_size = 10
}

data "ctyun_ecs_flavors" "ecs_flavor_test1" {
  cpu    = 2
  ram    = 4
  arch   = "x86"
  series = "S"
  type   = "CPU_S7"
}


resource "ctyun_ecs" "ecs_test1" {
  instance_name       = "tf1-ecs-test"
  display_name        = "tf1-ecs-test"
  flavor_id           = data.ctyun_ecs_flavors.ecs_flavor_test1.flavors[0].id
  image_id            = data.ctyun_images.image_test1.images[0].id
  system_disk_type    = "sata"
  system_disk_size    = 40
  vpc_id = ctyun_vpc.vpc_test.id
  password            = "P@ssW0rd_1"
  cycle_type          = "on_demand"
  subnet_id = ctyun_subnet.subnet_test.id
  is_destroy_instance = false
  monitor_service = false
}


resource "ctyun_vpce_server" "test" {
  name  = "tf-vpce-server-sss"
  vpc_id = ctyun_vpc.vpc_test.id
  subnet_id = ctyun_subnet.subnet_test.id
  auto_connection = true
  type = "interface"
  instance_id = ctyun_ecs.ecs_test1.id
  instance_type = "vm"
  rules = [{
    protocol = "TCP"
    endpoint_port = 2
    server_port = 2
  },
  ]
}

