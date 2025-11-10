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
data "ctyun_ec_cloud_gateways" "test" {
  ec_id       = "49410d6d-fd53-48b3-9f78-cb28da38d7be"
}

output "ctyun_ec_cloud_gateways_test" {
  value = data.ctyun_ec_cloud_gateways.test
}