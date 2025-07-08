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

data "ctyun_nat_dnats" "test" {
  nat_gateway_id = "natgw-asdsmh8scy"
}