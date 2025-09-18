resource "ctyun_mysql_backup_recovery" "%[1]s" {
  inst_id      = "%[2]s"
  project_id   = "%[3]s"
  src_inst_id  = "%[4]s"
  dst_inst_id  = "%[5]s"
  to_timepoint = %[6]s
}
