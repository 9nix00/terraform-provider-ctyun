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
  name        = "vpc-test-ccse"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  enable_ipv6 = true
}

resource "ctyun_subnet" "subnet_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "subnet-test-ccse"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns         = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4"
  ]
}

data "ctyun_images" "image_test" {
  name       = "CtyunOS 23"
  visibility = "public"
  page_no = 1
  page_size = 10
}

data "ctyun_ecs_flavors" "ecs_flavor_test" {
  cpu    = 4
  ram    = 8
  arch   = "x86"
  series = "C"
  type   = "CPU_C7"
}

resource "ctyun_security_group" "security_group_test" {
  vpc_id      = ctyun_vpc.vpc_test.id
  name        = "terraform-ccse-test55"
  description = "terraform测试使用"
}

resource "ctyun_ccse_cluster" "example" {
  base_info = {
    vpc_id     = ctyun_vpc.vpc_test.id
    subnet_id  = ctyun_subnet.subnet_test.id
    security_group_id = ctyun_security_group.security_group_test.id
    cluster_name = "newxiu_sl"
    cluster_domain = "www.ccc.s.a"
    network_plugin = "cubecni"
    start_port = 30000
    end_port   = 65535
    elb_prod_code = "standardI"
    pod_cidr    = "192.168.0.0/16"
    pod_subnet_id_list = [ctyun_subnet.subnet_test.id]
    cycle_type  = "on_demand"
    container_runtime = "containerd"
    timezone    = "Asia/Shanghai"
    cluster_version = "1.23.3"
    deploy_type   = "single"
    kube_proxy    = "iptables"
    cluster_series = "cce.standard"
  }

  master_host = {
    item_def_name =  data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].name

    az_infos = [
      {
        az_name = "cn-huadong1-jsnj1A-public-ctcloud"
        size    = 1
      }
    ]

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
  }

  slave_host = {
    instance_type = "ecs"
    mirror_id     = data.ctyun_images.image_test.images[0].id
    mirror_type   = 0
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