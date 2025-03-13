terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}


provider "ctyun" {
  region_id            = "200000001852"         # 如果此值不填，则默认读取环境变量中的CTYUN_REGION_ID
  az_name              = "cn-huabei2-tj-3a-public-ctcloud"        # 如果此值不填，则默认读取环境变量中的CTYUN_AZ_NAME
  env                  = "prod"                                     # 如果此值不填，则默认读取环境变量中的CTYUN_ENV
}

# 查找1c1g x86架构的通用型的规格
data "ctyun_ebm_device_types" "test1" {

}


output "ctyun_ecs_flavor_id1" {
  value = data.ctyun_ebm_device_types.test1
}

