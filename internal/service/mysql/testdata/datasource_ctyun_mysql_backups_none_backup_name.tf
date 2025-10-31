data "ctyun_mysql_backups" "%[1]s" {
  inst_id   = "%[2]s"
  page_no   = %[3]d
  page_size = %[4]d
}
