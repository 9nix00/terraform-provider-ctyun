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

# 查可用模板列表
# data "ctyun_ccse_template_market" "test" {
#
# }

# 指定名称时，可以查模板版本
# data "ctyun_ccse_template_market" "test" {
#   tpl_name = "elasticsearch"
# }

# 指定版本时，可以查模板Values
data "ctyun_ccse_template_market" "test" {
  tpl_name = "elasticsearch"
  tpl_version = "7.10.2"
  values_type = "YAML"
}

output "template" {
  value = data.ctyun_ccse_template_market.test
}