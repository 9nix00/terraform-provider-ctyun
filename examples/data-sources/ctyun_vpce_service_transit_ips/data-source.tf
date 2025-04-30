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

data "ctyun_vpce_service_transit_ips" "test" {
  endpoint_service_id = "endpser-pe60y2rtu6"
}