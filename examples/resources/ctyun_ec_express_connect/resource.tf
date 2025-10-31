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

resource "ctyun_express_connect" "express_connect_dependence" {
  name        = "express_connect_dependence"
  description = "云间高速开发测试专用"

}