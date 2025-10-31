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

data "ctyun_ec_routes" "routes_examples" {
  ec_id  = "49410d6d-fd53-48b3-9f78-cb28da38d7be"
  cgw_id = "333151a9-7d33-438b-8815-14b8a579b85d"
  rtb_id = "04f12286-2df8-41c7-a946-879e1919193b"
}