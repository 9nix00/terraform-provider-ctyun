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

resource "ctyun_ec_cda_instance" "example" {
  ec_id           = "example-ec-id"
  cgw_id          = "example-cgw-id"
  cda_id          = "example-cda-id"
  cda_name        = "example-cda-name"
  cda_cidr_v4_list = ["192.168.1.0/24"]
  rtb_id          = "example-rtb-id"
  cda_info        = "{\"key\": \"value\"}"
}