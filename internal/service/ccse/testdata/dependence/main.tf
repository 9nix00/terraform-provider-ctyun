data "ctyun_vpcs" "vpc_test" {

}

locals {
  vpcs = [for vpc in data.ctyun_vpcs.vpc_test.vpcs : vpc if vpc.name == "tf-vpc-for-paas"]
  data_vpc_id = length(local.vpcs) > 0 ? local.vpcs[0].vpc_id : ""
}

resource "ctyun_vpc" "vpc_test" {
  count    = local.data_vpc_id == "" ? 1 : 0
  name        = "tf-vpc-for-paas"
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
  subnets = [for subnet in data.ctyun_subnets.subnet_test.subnets : subnet if subnet.name == "tf-subnet-for-paas"]
  data_subnet_id = length(local.subnets) > 0 ? local.subnets[0].subnet_id : ""
}

resource "ctyun_subnet" "subnet_test" {
  count    = local.data_vpc_id == "" ? 1 : 0
  vpc_id      = local.real_vpc_id
  name        = "tf-subnet-for-paas"
  cidr        = "192.168.0.0/16"
  description = "terraform测试使用"
  dns         = [
    "8.8.8.8",
    "8.8.4.4"
  ]
}

locals {
  real_subnet_id = local.data_subnet_id == "" ? try(ctyun_subnet.subnet_test[0].id, "") : local.data_subnet_id
}

data "ctyun_ecs_flavors" "ecs_flavor_test" {
  cpu    = 4
  ram    = 8
  arch   = "x86"
  series = "C"
  type   = "CPU_C7"
}

locals {
  cluster_name = "tf-ccse-cluster-${local.random_string}"
}

resource "ctyun_ccse_cluster" "test" {
  base_info = {
    vpc_id     = local.real_vpc_id
    subnet_id  = local.real_subnet_id
    cluster_name = local.cluster_name
    cluster_domain = "www.ctyun.com"
    network_plugin = "cubecni"
    start_port = 30000
    end_port   = 65535
    elb_prod_code = "standardI"
    pod_subnet_id_list = [local.real_subnet_id]
    cycle_type  = "on_demand"
    container_runtime = "containerd"
    timezone    = "Asia/Shanghai"
    cluster_version = "1.23.3"
    deploy_type   = "single"
    kube_proxy    = "ipvs"
    cluster_series = "cce.managed"
    series_type = "managedbase"
  }


  slave_host = {
    instance_type = "ecs"
    mirror_id     = "3f80d8c0-8eb5-4afa-a506-13ba68b61872"
    mirror_type   = 1
    item_def_name = data.ctyun_ecs_flavors.ecs_flavor_test.flavors[0].name

    az_infos = [
      {
        az_name = "cn-huadong1-jsnj1A-public-ctcloud"
        size    = 1
      }
    ]

    sys_disk = {
      type = "SSD"
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


locals {
  chart_name = "node-problem-detector"
}

data "ctyun_ccse_plugin_market" "test" {
  chart_name = local.chart_name
}

locals {
  chart_version1 =try(data.ctyun_ccse_plugin_market.test.versions[2].chart_version, "")
  chart_version2 =try(data.ctyun_ccse_plugin_market.test.versions[1].chart_version, "")
}

data "ctyun_ccse_plugin_market" "test1" {
  chart_name = local.chart_name
  chart_version = local.chart_version1
  values_type = "YAML"
}

data "ctyun_ccse_plugin_market" "test2" {
  chart_name = local.chart_name
  chart_version = local.chart_version2
  values_type = "JSON"
}

locals {
  # 生成当前时间戳的哈希值
  hash = sha256(timestamp())

  # 从哈希结果中截取字符（转为小写并移除特殊字符）
  random_string = substr(
    replace(
      lower(local.hash),
      "/[^a-z0-9]/",
      ""  # 移除所有非字母数字的字符
    ),
    0, 10  # 截取前16个字符
  )
}