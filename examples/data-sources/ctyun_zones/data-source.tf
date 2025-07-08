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

# 从环境变量读取region_id
data "ctyun_zones" "test" {

}


# 指定资源池ID
# data "ctyun_zones" "test" {
#   region_id = "200000002368"
# }

output "ctyun_test" {
  value = data.ctyun_zones.test
}

