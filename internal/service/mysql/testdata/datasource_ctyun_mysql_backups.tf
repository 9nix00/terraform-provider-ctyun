data "ctyun_mysql_backups" "%[1]s" {
  inst_id     = "%[2]s"
  backup_name = %[3]s
}
