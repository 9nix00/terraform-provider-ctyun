resource "ctyun_mysql_rds_parameter_template" "%[1]s" {
  inst_id    = "%[2]s"
  project_id = "%[3]s"
  parameters = %[4]s
}