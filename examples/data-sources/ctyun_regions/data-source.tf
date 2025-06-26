terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 全部资源池
data "ctyun_regions" "test" {

}

# 筛选多AZ资源池
locals {
  multi_az_regions = [for region in data.ctyun_regions.test.regions : region if length(region.zones) > 0]
}

# 指定资源池
# data "ctyun_regions" "test" {
#   name = "南昌"
# }

output "ctyun_regions_test" {
  value =local.multi_az_regions
}
