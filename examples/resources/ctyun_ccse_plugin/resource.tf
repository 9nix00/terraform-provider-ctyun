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

data "ctyun_ccse_plugin_market" "ccse_monitor" {
  chart_name = "ccse-monitor"
  chart_version = "0.1.9"
  values_type = "YAML"
}

resource "ctyun_ccse_plugin" "example1" {
  cluster_id = "6bb243ec40ce4628a0e8ccf1028a10fd"
  chart_name = "ccse-monitor"
  chart_version = "0.1.9"
  plugin_name = "tf-ccse-plugin1"
  values_yaml = data.ctyun_ccse_plugin_market.ccse_monitor.values
}