data "ctyun_vpcs" "vpc_test" {
  page_size = 50
}

locals {
  vpcs        = [for vpc in data.ctyun_vpcs.vpc_test.vpcs : vpc if vpc.name == "tf-vpc-for-iaas"]
  data_vpc_id = length(local.vpcs) > 0 ? local.vpcs[0].vpc_id : ""
}

resource "ctyun_vpc" "vpc_test" {
  count       = local.data_vpc_id == "" ? 1 : 0
  name        = "tf-vpc-for-iaas"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

locals {
  real_vpc_id = local.data_vpc_id == "" ? try(ctyun_vpc.vpc_test[0].id, "") : local.data_vpc_id
}


data "ctyun_subnets" "subnet_test" {
  vpc_id = local.real_vpc_id
}

locals {
  subnets = [
    for subnet in data.ctyun_subnets.subnet_test.subnets : subnet if subnet.name == "tf-subnet-for-iaas"
  ]
  data_subnet_id = length(local.subnets) > 0 ? local.subnets[0].subnet_id : ""
}

resource "ctyun_subnet" "subnet_test" {
  count       = local.data_vpc_id=="" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-subnet-for-iaas"
  cidr        = "192.168.1.0/24"
  description = "terraform测试使用"
  dns = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

locals {
  real_subnet_id = local.data_subnet_id == "" ? try(ctyun_subnet.subnet_test[0].id, "") : local.data_subnet_id
}


resource "ctyun_subnet" "subnet_test1" {
  vpc_id      = local.real_vpc_id
  name        = "tf-subnet-for-iaas1"
  cidr        = "192.168.2.0/24"
  description = "terraform测试使用"
  dns = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}


data "ctyun_security_groups" "security_group_test" {
  vpc_id = local.real_vpc_id
}

locals {
  security_groups = [
    for security_group in data.ctyun_security_groups.security_group_test.security_groups :security_group if security_group.name == "tf-sg-for-iaas"
  ]
  data_security_group_id = length(local.security_groups) > 0 ? local.security_groups[0].security_group_id : ""
}

resource "ctyun_security_group" "security_group_test" {
  count       = local.data_vpc_id=="" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-sg-for-iaas"
  description = "terraform测试使用"
  lifecycle {
    prevent_destroy = false
  }
}

locals {
  real_security_group_id = local.data_security_group_id == "" ? try(ctyun_security_group.security_group_test[0].id, "") : local.data_security_group_id
}


resource "ctyun_security_group" "security_group_test1" {
  vpc_id      = local.real_vpc_id
  name        = "tf-sg-for-scaling1"
  description = "terraform测试使用"
  lifecycle {
    prevent_destroy = false
  }
}


data "ctyun_images" "image_test" {
  name       = "CentOS Linux 8.4"
  visibility = "public"
  page_no = 1
  page_size = 10
}

# data "ctyun_images" "image_test1" {
#   name       = "CentOS Linux 8.2"
#   visibility = "public"
#   page_no = 1
#   page_size = 10
# }


locals {
  image_id = data.ctyun_images.image_test.images[0].id
  # image_id1 = data.ctyun_images.image_test.images[0].id
}

resource "ctyun_keypair" "scaling_test" {
  name       = "key-pair-scaling-test"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjUnAnTid4wmVtajSmElMtH03OvOyY81ybfswbUu9Gt83DVVzDnwb3rcQW1us8SeKm/gRINkgdrRAgfXAmTKR7AorYtWWc/tzb6kcDpL2E8Qk+n6cyFAxXNoX2vXBr4kC9wz1uwjGyxoSlpHLIpscfI0Ef652gMlSyfODehAJHj3JPMr8pvtPIUqsZI3JOGTUzxaA2JVC0LxQegphYYf2TxGd9GLRUv1p/0BUAPCMg1NaITXNVEj3A11hk1nrFoJMmvIwIUkLmRuQcxuNAdxeLB7GXXVjKpnKIJL4L64dyA9GWa3Gb7gCJyRaBc5UhK4hT57wmukCrldHHtdF1IJr"
}

resource "ctyun_scaling_config" "config_test" {
  name            = "sc-for-policy"
  image_id        =  local.image_id
  flavor_name     = "s7.large.2"
  use_floatings   = "diable"
  login_mode      = "key_pair"
  key_pair_id     = ctyun_keypair.scaling_test.id
  monitor_service = true
  az_names        = ["cn-huadong1-jsnj1A-public-ctcloud"]
  volumes         = [{"volume_type":"SATA", "volume_size": 40, "flag":"OS"}]
}

resource "ctyun_scaling_config" "config_test1" {
  name            = "sc-for-policy1"
  image_id        =  local.image_id
  flavor_name     = "s7.large.2"
  use_floatings   = "diable"
  login_mode      = "key_pair"
  key_pair_id     = ctyun_keypair.scaling_test.id
  monitor_service = true
  az_names        = ["cn-huadong1-jsnj1A-public-ctcloud"]
  volumes         = [{"volume_type":"SATA", "volume_size": 40, "flag":"OS"}]
}

resource "ctyun_scaling_group" "group_test" {
  security_group_id_list = [local.real_security_group_id]
  name                   = "scaling-group-for-policy"
  health_mode            = "server"
  subnet_id_list         = [local.real_subnet_id]
  move_out_strategy      = "earlier_config"
  vpc_id                 = local.real_vpc_id
  min_count              = 1
  max_count              = 3
  health_period          = 300
  use_lb                 = 2
  config_list            = [ctyun_scaling_config.config_test.id]
  az_strategy            = "uniform_distribution"
  delete_protection      = "disable"
}


data "ctyun_ecs_flavors" "ecs_flavor_test" {
  cpu    = 2
  ram    = 4
  arch   = "x86"
  series = "C"
  type   = "CPU_C7"
}
#
# #
resource "ctyun_ecs" "ecs_test" {
  instance_name       = "tf-ecs-for-scaling-ecs1"
  display_name        = "tf-ecs-for-scaling-ecs1"
  flavor_id           = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].id
  image_id            = data.ctyun_images.image_test.images[0].id
  system_disk_type    = "sata"
  system_disk_size    = 40
  vpc_id =  local.real_vpc_id
  password            = "P@ssW0rd_1"
  cycle_type          = "on_demand"
  subnet_id = local.real_subnet_id
  security_group_ids = [local.real_security_group_id]
  is_destroy_instance = true
}
#
#
resource "ctyun_ecs" "ecs_test1" {
  instance_name       = "tf-ecs-for-scaling-ecs2"
  display_name        = "tf-ecs-for-scaling-ecs2"
  flavor_id           = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].id
  image_id            = data.ctyun_images.image_test.images[0].id
  system_disk_type    = "sata"
  system_disk_size    = 40
  vpc_id =  local.real_vpc_id
  password            = "P@ssW0rd_1"
  cycle_type          = "on_demand"
  subnet_id = local.real_subnet_id
  security_group_ids = [local.real_security_group_id]
  is_destroy_instance = true
}

resource "ctyun_ecs" "ecs_test2" {
  instance_name       = "tf-ecs-for-scaling-ecs3"
  display_name        = "tf-ecs-for-scaling-ecs3"
  flavor_id           = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].id
  image_id            = data.ctyun_images.image_test.images[0].id
  system_disk_type    = "sata"
  system_disk_size    = 40
  vpc_id =  local.real_vpc_id
  password            = "P@ssW0rd_1"
  cycle_type          = "on_demand"
  subnet_id = local.real_subnet_id
  security_group_ids = [local.real_security_group_id]
  is_destroy_instance = true
}

#
resource "ctyun_elb_loadbalancer" "elb_test" {
  subnet_id     = local.real_subnet_id
  name          = "tf-elb-for-scaling-group"
  sla_name      = "elb.s2.small"
  resource_type = "internal"
  vpc_id        = local.real_vpc_id
  cycle_type    = "on_demand"
}

resource "ctyun_elb_health_check" "test" {
  name     = "tf-hc-for-scaling"
  protocol = "TCP"
}

resource "ctyun_elb_target_group" "target_group_test" {
  name      = "tf_target_group"
  vpc_id    = local.real_vpc_id
  algorithm = "wrr"
  health_check_id = ctyun_elb_health_check.test.id
  session_sticky_mode = "SOURCE_IP"
  source_ip_timeout = 30
  proxy_protocol = 1
  protocol = "TCP"
}

resource "ctyun_elb_listener" "elb_listener_test" {
  loadbalancer_id     = ctyun_elb_loadbalancer.elb_test.id
  name                = "tf_listener_scaling"
  protocol            = "TCP"
  protocol_port       = 12345
  default_action_type = "forward"
  target_groups = [{ target_group_id = ctyun_elb_target_group.target_group_test.id }]
  listener_cps        = 1
  establish_timeout   = 100
}

resource "ctyun_elb_loadbalancer" "elb_test1" {
  subnet_id     = local.real_subnet_id
  name          = "tf-elb-for-scaling-group1"
  sla_name      = "elb.s2.small"
  resource_type = "internal"
  vpc_id        = local.real_vpc_id
  cycle_type    = "on_demand"
}

resource "ctyun_elb_health_check" "test1" {
  name     = "tf-hc-for-scaling1"
  protocol = "TCP"
}

resource "ctyun_elb_target_group" "target_group_test1" {
  name      = "tf_target_group1"
  vpc_id    = local.real_vpc_id
  algorithm = "wrr"
  health_check_id = ctyun_elb_health_check.test1.id
  session_sticky_mode = "SOURCE_IP"
  source_ip_timeout = 30
  proxy_protocol = 1
  protocol = "TCP"
}

resource "ctyun_elb_listener" "elb_listener_test1" {
  loadbalancer_id     = ctyun_elb_loadbalancer.elb_test1.id
  name                = "tf_listener_scaling1"
  protocol            = "TCP"
  protocol_port       = 12345
  default_action_type = "forward"
  target_groups = [{ target_group_id = ctyun_elb_target_group.target_group_test1.id }]
  listener_cps        = 1
  establish_timeout   = 100
}