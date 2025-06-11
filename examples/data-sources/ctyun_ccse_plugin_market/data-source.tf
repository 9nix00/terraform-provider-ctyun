terraform {
  required_providers {
    ctyun = {
      source = "ctyun-it/ctyun"
    }
  }
}

# 查可用插件列表
# data "ctyun_ccse_plugin_market" "test" {
#
# }

# 指定名称时，可以查插件版本
# data "ctyun_ccse_plugin_market" "test" {
#   chart_name = "cubevk-profile"
# }

# 指定版本时，可以查插件Values
data "ctyun_ccse_plugin_market" "test" {
  chart_name = "cubevk-profile"
  chart_version = "1.0.4"
  values_type = "YAML"
}

output "plugin" {
  value = data.ctyun_ccse_plugin_market.test
}