data "ctyun_vpcs" "vpc_test" {
  page_size = 50
}

locals {
  vpcs        = [for vpc in data.ctyun_vpcs.vpc_test.vpcs : vpc if vpc.name == "tf-vpc-for-paas"]
  data_vpc_id = length(local.vpcs) > 0 ? local.vpcs[0].vpc_id : ""
}

resource "ctyun_vpc" "vpc_test" {
  count       = local.data_vpc_id == "" ? 1 : 0
  name        = "tf-vpc-for-paas"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}
locals {
  real_vpc_id = local.data_vpc_id == "" ? try(ctyun_vpc.vpc_test[0].id, "") : local.data_vpc_id
}

#
# resource "ctyun_vpc" "vpc_test" {
#   name        = "tf-vpc-for-vpc"
#   cidr        = "192.168.0.0/16"
#   description = "terraform测试使用"
#   enable_ipv6 = true
# }
data "ctyun_subnets" "subnet_test" {
  vpc_id = local.real_vpc_id
}
locals {
  subnets = [
    for subnet in data.ctyun_subnets.subnet_test.subnets : subnet if subnet.name == "tf-subnet-for-paas"
  ]
  data_subnet_id = length(local.subnets) > 0 ? local.subnets[0].subnet_id : ""
}

resource "ctyun_subnet" "subnet_test" {
  count       = local.data_vpc_id=="" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-subnet-for-paas"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

locals {
  real_subnet_id = local.data_subnet_id == "" ? try(ctyun_subnet.subnet_test[0].id, "") : local.data_subnet_id
}
resource "ctyun_security_group" "security_group_test" {
  vpc_id      = local.real_vpc_id
  name        = "tf-sg-for-vpc"
  description = "terraform测试使用"
}

resource "ctyun_bandwidth" "bandwidth_test" {
  name       = "tf-bandwidth-for-vpc"
  bandwidth  = 5
  cycle_type  = "on_demand"
}

resource "ctyun_eip" "eip_test" {
  name        = "tf-eip-for-vpc"
  bandwidth   = 5
  cycle_type = "on_demand"
  demand_billing_type = "upflowc"
}

data "ctyun_images" "image_test" {
  name       = "CentOS Linux 8.4"
  visibility = "public"
  page_no = 1
  page_size = 10
}

data "ctyun_ecs_flavors" "ecs_flavor_test" {
  cpu    = 2
  ram    = 4
  arch   = "x86"
  series = "C"
  type   = "CPU_C7"
}

resource "ctyun_ecs" "ecs_test" {
  instance_name       = "tf-ecs-for-vpc"
  display_name        = "tf-ecs-for-vpc"
  flavor_id           = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].id
  image_id            = data.ctyun_images.image_test.images[0].id
  system_disk_type    = "sata"
  system_disk_size    = 40
  vpc_id = local.real_vpc_id
  password            = var.password
  cycle_type          = "on_demand"
  subnet_id = local.real_subnet_id
}

variable "password" {
  type      = string
  sensitive = true
}
#
#
#
#
#
# # 创建数据盘资源
resource "ctyun_ebs" "data_disk_test" {
  name       = "tf-test-data-disk"
  mode       = "vbd"
  type       = "sata"
  size       = 60
  cycle_type = "on_demand"
}

# 创建EBS与ECS的关联关系（显式挂载）
resource "ctyun_ebs_association_ecs" "data_disk_association" {
  instance_id = ctyun_ecs.ecs_test.id
  ebs_id      = ctyun_ebs.data_disk_test.id
}
# 查询网络接口资源

resource "ctyun_port" "port_test" {
  name       = "tf-test-port"
  description                = "tf-test-port"
  subnet_id                  =local.real_subnet_id
  # security_group_ids        = ["%[5]s"]
  secondary_private_ip_count = 1
}
resource "ctyun_ecs_port_association" "port_association" {
  instance_id          = ctyun_ecs.ecs_test.id
  port_id = ctyun_port.port_test.id
}
resource "ctyun_vip" "vip_test" {
  subnet_id  = local.real_subnet_id
  vpc_id     = local.real_vpc_id
  ip_address = "192.168.100.152"
  vip_type   = "v4"
}

resource "ctyun_dhcpoptionset" "dhcpoptionset_test" {
  name         = "tf-dhcpoptionset-test"
  description  = "terraform测试使用"
  domain_name  = "example.com"
  dns_list     = ["8.8.8.8", "8.8.4.4"]
}