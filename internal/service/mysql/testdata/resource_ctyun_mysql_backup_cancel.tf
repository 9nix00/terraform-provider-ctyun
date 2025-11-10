resource "ctyun_mysql_backup_cancel" "%[1]s" {
  inst_id     = "%[2]s"
  project_id  = "%[3]s"
  backup_record_id = %[4]s
}
