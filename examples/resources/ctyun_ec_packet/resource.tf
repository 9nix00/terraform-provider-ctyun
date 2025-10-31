terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 可参考index.md，在环境变量中配置ak、sk、资源池ID、可用区名称
provider "ctyun" {
  env = "prod"
}
resource "ctyun_express_connect" "example" {
  name        = "express_connect_dependence"
  description = "云间高速开发测试专用"

}

# 创建云间高速带宽包
resource "ctyun_ec_packet" "example" {
  ec_id        = ctyun_express_connect.example.id
  packet_name  = "example-ec-packet"
  bandwidth    = 10
  cycle_type   = "MONTH"
  cycle_count  = 1

  # 可选参数
  on_demand    = false
  area_a       = "china"
  area_b       = "china"

}

