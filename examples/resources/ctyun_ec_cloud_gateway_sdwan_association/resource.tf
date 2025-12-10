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
resource "ctyun_ec_cloud_gateway" "example" {
  ec_id       = ctyun_express_connect.example.id
  name        = "cloud_gateway_xinan1"
  description = "云间高速开发测试专用"
  region_id   = "200000002368"
  region_name = "cn-xinan1-xn1A-public-ctcloud"
}
resource "ctyun_express_connect" "example" {
  name        = "express_connect_dependence"
  description = "云间高速example专用"
}
resource "ctyun_sdwan" "demo" {
  name = "sdwan_demo"
  description = "样列"
}

resource "ctyun_ec_cloud_gateway" "cloud_gateway_dependence" {
  ec_id       = ctyun_express_connect.example.id
  name        = "cloud_gateway_dependence"
  description = "云间高速开发测试专用"
}
resource "ctyun_ec_cloud_gateway_sdwan_association" "example" {
  ec_id    = ctyun_express_connect.example.id
  sdwan_id = ctyun_sdwan.demo.id
  cgw_list = [{
    cgw_id = ctyun_ec_cloud_gateway.example.id
    rtb_id = ctyun_ec_cloud_gateway.cloud_gateway_dependence.id
  }]
}