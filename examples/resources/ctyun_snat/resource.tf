terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

provider "ctyun" {
  region_id            = "200000002530"
}

resource "ctyun_snat" "snat_create_test"{
    region_id = "200000002530"
    nat_gateway_id = ""
    snat_ips = ""
    sourceCIDR = ""
    description = ""
}