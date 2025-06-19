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
  name        = "vpc-test-ccse1"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "subnet-test-ccse1"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
}

data "ctyun_ecs_flavors" "ecs_flavor_test" {
  cpu    = 4
  ram    = 8
  arch   = "x86"
  series = "C"
  type   = "CPU_C7"
}


resource "ctyun_ccse_cluster" "example" {
  base_info = {
    vpc_id     = ctyun_vpc.vpc_test.id
    subnet_id  = ctyun_subnet.subnet_test.id
    cluster_name = "fe-ccse1"
    cluster_domain = "www.ccc.s"
    network_plugin = "cubecni"
    start_port = 30000
    end_port   = 65535
    elb_prod_code = "standardI"
    pod_subnet_id_list = [ctyun_subnet.subnet_test.id]
    cycle_type  = "month"
    cycle_count = 1
    container_runtime = "containerd"
    timezone    = "Asia/Shanghai"
    cluster_version = "1.25.6"
    deploy_type   = "single"
    kube_proxy    = "iptables"
    cluster_series = "cce.standard"
  }

  master_host = {
    item_def_name =  data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].name

    sys_disk = {
      type = "SSD"
      size = 100
    }

    data_disks = [
      {
        type = "SSD"
        size = 200
      }
    ]

    az_infos = [
      {
        az_name = "cn-huadong1-jsnj1A-public-ctcloud"
        size    = 1
      }
    ]
  }

  slave_host = {
    instance_type = "ecs"
    mirror_id     = "3f80d8c0-8eb5-4afa-a506-13ba68b61872"
    mirror_type   = 1
    item_def_name = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].name

    az_infos = [
      {
        az_name = "cn-huadong1-jsnj2A-public-ctcloud"
        size    = 1
      }
    ]

    sys_disk = {
      type = "SATA"
      size = 80
    }

    data_disks = [
      {
        type = "SATA"
        size = 150
      }
    ]
  }
}

# resource "ctyun_ccse_cluster" "example2" {
#   base_info = {
#     vpc_id     = ctyun_vpc.vpc_test.id
#     subnet_id  = ctyun_subnet.subnet_test.id
#     cluster_name = "auto-sec-grqq33"
#     cluster_domain = "www.ccc.s"
#     network_plugin = "cubecni"
#     start_port = 30000
#     end_port   = 65535
#     elb_prod_code = "standardI"
#     pod_subnet_id_list = [ctyun_subnet.subnet_test.id]
#     cycle_type  = "on_demand"
#     container_runtime = "containerd"
#     timezone    = "Asia/Shanghai"
#     cluster_version = "1.23.3"
#     deploy_type   = "single"
#     kube_proxy    = "ipvs"
#     cluster_series = "cce.managed"
#     series_type = "managedbase"
#   }
#
#
#   slave_host = {
#     instance_type = "ecs"
#     mirror_id     = "3f80d8c0-8eb5-4afa-a506-13ba68b61872"
#     mirror_type   = 1
#     item_def_name = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].name
#
#     az_infos = [
#       {
#         az_name = "cn-huadong1-jsnj2A-public-ctcloud"
#         size    = 1
#       }
#     ]
#
#     sys_disk = {
#       type = "SATA"
#       size = 80
#     }
#
#     data_disks = [
#       {
#         type = "SATA"
#         size = 150
#       }
#     ]
#   }
# }