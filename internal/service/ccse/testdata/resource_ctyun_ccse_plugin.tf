resource "ctyun_ccse_plugin" "%[1]s" {
  plugin_name = "%[2]s"
  cluster_id = "%[3]s"
  chart_name = "%[4]s"
  chart_version = "%[5]s"
  %[6]s
}
