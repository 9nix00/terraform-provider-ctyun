terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

provider "ctyun" {
  env = "prod"
  region_id   = "200000001852"
}

resource "ctyun_bandwidth" "bandwidth_test1" {
  name       = "tf-bandwidth-test3"
  cycle_type = "on_demand"
  bandwidth  = 6
}

# 创建一个包年，大小为10Mbit/s的带宽
# resource "ctyun_bandwidth" "bandwidth_test2" {
#   name        = "tf-bandwidth-test2"
#   cycle_type  = "month"
#   bandwidth   = 10
#   cycle_count = 1
# }

data "ctyun_bandwidths" "test" {
  bandwidth_id = ctyun_bandwidth.bandwidth_test1.id
}