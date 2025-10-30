
resource "ctyun_express_connect" "express_connect_dependence" {
  name        = "express_connect_dependence"
  description = "云间高速开发测试专用"
}


resource "ctyun_ec_cloud_gateway" "cloud_gateway_dependence" {
  ec_id    = ctyun_express_connect.express_connect_dependence.id
  name     = "cloud_gateway_dependence"
  description = "云间高速开发测试专用"
  region_id = "200000002401"
  region_name = "cn-hn-cs42-hncs1A-public-ctcloud"
}

