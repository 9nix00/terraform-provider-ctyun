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

data "ctyun_nat_snats" "test"{
  nat_gateway_id = "natgw-asdsmh8scy"
}