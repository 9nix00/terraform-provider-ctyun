resource "ctyun_mysql_backup_setting" "%[1]s" {
  inst_id                    = "%[2]s"
  project_id                 = "%[3]s"
  storage_day                = %[4]d
  allow_earliest_time        = "%[7]s"
  trigger_days_of_week       = %[8]s
}
